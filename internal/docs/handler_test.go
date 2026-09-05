package docs_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"erm-dokter/internal/docs"
)

func TestDocsHandler_ServeUI(t *testing.T) {
	handler := docs.NewHandler()
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected Content-Type text/html, got '%s'", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "@scalar/api-reference") {
		t.Error("expected body to contain @scalar/api-reference")
	}
	if !strings.Contains(body, "/docs/openapi.yaml") {
		t.Error("expected body to reference /docs/openapi.yaml")
	}
	if !strings.Contains(body, "initLiveReload") {
		t.Error("expected body to contain live reload script")
	}
}

func TestDocsHandler_ServeOpenAPISpec(t *testing.T) {
	handler := docs.NewHandler()
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/yaml") {
		t.Errorf("expected Content-Type application/yaml, got '%s'", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "openapi: 3.1.0") {
		t.Error("expected body to contain 'openapi: 3.1.0'")
	}
	if !strings.Contains(body, "title: ERM Dokter API") {
		t.Error("expected body to contain 'title: ERM Dokter API'")
	}
}

func TestDocsHandler_ServeOpenAPISpec_CustomFileAndHead(t *testing.T) {
	tempDir := t.TempDir()
	specFile := filepath.Join(tempDir, "custom_spec.yaml")
	expectedContent := "openapi: 3.1.0\ninfo:\n  title: Custom Spec\n"
	if err := os.WriteFile(specFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("failed to write test spec file: %v", err)
	}

	handler := docs.NewHandlerWithFilePath(specFile)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Test GET custom file
	req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != expectedContent {
		t.Errorf("expected body '%s', got '%s'", expectedContent, w.Body.String())
	}

	// Test HEAD request
	reqHead := httptest.NewRequest(http.MethodHead, "/docs/openapi.yaml", nil)
	wHead := httptest.NewRecorder()
	mux.ServeHTTP(wHead, reqHead)

	if wHead.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", wHead.Code)
	}
	if wHead.Body.Len() != 0 {
		t.Errorf("expected empty body for HEAD request, got %d bytes", wHead.Body.Len())
	}
}

func TestDocsHandler_ServeLiveReload(t *testing.T) {
	tempDir := t.TempDir()
	specFile := filepath.Join(tempDir, "test_spec.yaml")
	if err := os.WriteFile(specFile, []byte("initial spec"), 0644); err != nil {
		t.Fatalf("failed to write test spec file: %v", err)
	}

	handler := docs.NewHandlerWithFilePath(specFile)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/docs/live-reload", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		mux.ServeHTTP(w, req)
	}()

	// Tunggu sebentar agar initial ping terkirim
	time.Sleep(100 * time.Millisecond)

	// Simulasikan file openapi diupdate
	time.Sleep(100 * time.Millisecond)
	newModTime := time.Now().Add(1 * time.Second)
	if err := os.WriteFile(specFile, []byte("updated spec"), 0644); err != nil {
		t.Fatalf("failed to update spec file: %v", err)
	}
	_ = os.Chtimes(specFile, newModTime, newModTime)

	// Tunggu response reload
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		cancel()
		t.Fatal("timed out waiting for live reload notification")
	}

	body := w.Body.String()
	if !strings.Contains(body, ": ping") {
		t.Errorf("expected stream to contain ': ping', got: %s", body)
	}
	if !strings.Contains(body, "data: reload") {
		t.Errorf("expected stream to contain 'data: reload', got: %s", body)
	}
}

