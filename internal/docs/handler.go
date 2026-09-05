package docs

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"time"
)

//go:embed openapi.yaml
var openAPISpec []byte

const scalarHTML = `<!doctype html>
<html lang="id">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>ERM Dokter - Dokumentasi API</title>
    <link rel="icon" type="image/svg+xml" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%236366f1'><path d='M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-2 10h-4v4h-2v-4H7v-2h4V7h2v4h4v2z'/></svg>" />
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/docs/openapi.yaml"
      data-configuration='{
        "theme": "purple",
        "darkMode": true,
        "layout": "modern",
        "showSidebar": true,
        "searchHotKey": "k",
        "authentication": {
          "preferredSecurityScheme": "BearerAuth"
        }
      }'>
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      (function() {
        function initLiveReload() {
          if (!window.EventSource) return;
          var es = new EventSource('/docs/live-reload');
          var wasConnected = false;

          es.onopen = function() {
            if (wasConnected) {
              console.log('[LiveReload] Server terhubung kembali, me-refresh...');
              location.reload();
            }
            wasConnected = true;
          };

          es.onmessage = function(e) {
            if (e.data === 'reload') {
              console.log('[LiveReload] openapi.yaml diperbarui, me-reload halaman...');
              location.reload();
            }
          };

          es.onerror = function() {
            // Jangan pernah panggil location.reload() dari onerror untuk mencegah infinite reload loop
            console.warn('[LiveReload] Koneksi live reload terputus, menunggu server online...');
          };
        }
        initLiveReload();
      })();
    </script>
  </body>
</html>`

type Handler struct {
	specFilePath string
}

func NewHandler() *Handler {
	return &Handler{
		specFilePath: "internal/docs/openapi.yaml",
	}
}

func NewHandlerWithFilePath(filePath string) *Handler {
	return &Handler{
		specFilePath: filePath,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /docs", h.ServeUI)
	mux.HandleFunc("GET /docs/openapi.yaml", h.ServeOpenAPISpec)
	mux.HandleFunc("GET /docs/live-reload", h.ServeLiveReload)
}

func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(scalarHTML))
}

func (h *Handler) ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	data := openAPISpec
	if h.specFilePath != "" {
		if fileContent, err := os.ReadFile(h.specFilePath); err == nil {
			data = fileContent
		}
	}

	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = w.Write(data)
	}
}

func (h *Handler) ServeLiveReload(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lastMod time.Time
	if h.specFilePath != "" {
		if info, err := os.Stat(h.specFilePath); err == nil {
			lastMod = info.ModTime()
		}
	}

	_, _ = fmt.Fprintf(w, "retry: 1500\n\n")
	_, _ = fmt.Fprintf(w, ": ping\n\n")
	if err := rc.Flush(); err != nil {
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		} else {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if h.specFilePath == "" {
				continue
			}
			info, err := os.Stat(h.specFilePath)
			if err != nil {
				continue
			}
			if lastMod.IsZero() {
				lastMod = info.ModTime()
				continue
			}
			if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				_, _ = fmt.Fprintf(w, "data: reload\n\n")
				_ = rc.Flush()
				return
			}
		}
	}
}
