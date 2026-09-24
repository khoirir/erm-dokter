package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/pkg/response"
)

func TestRequestIDMiddleware_AutoGenerate(t *testing.T) {
	var capturedID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = logger.GetRequestID(r.Context())
		response.Success(w, "OK", nil)
	})

	wrapped := middleware.RequestIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	headerID := rr.Header().Get("X-Request-ID")
	if headerID == "" {
		t.Fatal("Expected non-empty X-Request-ID header")
	}
	if len(headerID) != 16 {
		t.Errorf("Expected 16 hex characters, got length %d ('%s')", len(headerID), headerID)
	}
	if strings.HasPrefix(headerID, "req_") {
		t.Errorf("Expected request ID without 'req_' prefix, got '%s'", headerID)
	}

	if capturedID != headerID {
		t.Errorf("Expected context ID '%s' to match header ID '%s'", capturedID, headerID)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("Failed to parse JSON body: %v", err)
	}
	if _, exists := body["request_id"]; exists {
		t.Errorf("Expected JSON body not to contain request_id, got '%v'", body["request_id"])
	}
}

func TestRequestIDMiddleware_PreserveIncomingHeader(t *testing.T) {
	customID := "client-trace-id-998877"
	var capturedID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = logger.GetRequestID(r.Context())
		response.Error(w, http.StatusBadRequest, "Invalid input", nil)
	})

	wrapped := middleware.RequestIDMiddleware(handler)

	req := httptest.NewRequest("POST", "/api/v1/test", nil)
	req.Header.Set("X-Request-ID", customID)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	headerID := rr.Header().Get("X-Request-ID")
	if headerID != customID {
		t.Errorf("Expected preserved header ID '%s', got '%s'", customID, headerID)
	}
	if capturedID != customID {
		t.Errorf("Expected captured context ID '%s', got '%s'", customID, capturedID)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("Failed to parse JSON body: %v", err)
	}
	if _, exists := body["request_id"]; exists {
		t.Errorf("Expected JSON body not to contain request_id, got '%v'", body["request_id"])
	}
}
