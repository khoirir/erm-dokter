package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/health"
)

func TestHealthCheck(t *testing.T) {
	handler := health.NewHandler()
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Data    map[string]string `json:"data"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true, got false")
	}

	if resp.Data["status"] != "healthy" || resp.Data["service"] != "erm-dokter" {
		t.Errorf("Unexpected health data: %+v", resp.Data)
	}
}
