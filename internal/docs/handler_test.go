package docs_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
