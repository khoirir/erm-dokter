package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"erm-dokter/pkg/response"
)

// timeoutResponseWriter adalah wrapper thread-safe untuk http.ResponseWriter.
// Memastikan hanya satu goroutine yang bisa menulis response.
type timeoutResponseWriter struct {
	http.ResponseWriter
	mu          sync.Mutex
	written     bool
	headersSent bool
}

func (tw *timeoutResponseWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.written {
		return
	}
	tw.headersSent = true
	tw.ResponseWriter.WriteHeader(code)
}

func (tw *timeoutResponseWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.written {
		return 0, nil
	}
	return tw.ResponseWriter.Write(b)
}

// markWritten menandai bahwa response sudah ditulis.
// Mengembalikan true jika berhasil menandai (belum pernah ditulis sebelumnya).
func (tw *timeoutResponseWriter) markWritten() bool {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.written {
		return false
	}
	tw.written = true
	return true
}

func TimeoutMiddleware(timeout time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			tw := &timeoutResponseWriter{ResponseWriter: w}

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					// Hanya tulis timeout response jika handler belum menulis apapun
					if tw.markWritten() {
						response.Error(w, http.StatusGatewayTimeout, "Request timeout", nil)
					}
				}
			}
		}
	}
}
