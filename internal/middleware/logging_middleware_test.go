package middleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/pkg/token"
)

func TestLoggingMiddleware_AuthenticatedAndAnonymous(t *testing.T) {
	var buf bytes.Buffer
	testLogger := logger.NewWithOptions(&buf, "json", slog.LevelDebug, "erm-dokter")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	loggedHandler := middleware.LoggingMiddlewareWithLogger(testLogger)(dummyHandler)

	req1 := httptest.NewRequest("GET", "/health", nil)
	rr1 := httptest.NewRecorder()
	loggedHandler.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr1.Code)
	}

	secret := "test-secret-key-1234567890123456"
	tkn, err := token.GenerateToken("DR001", "dr. Budi", secret, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req2 := httptest.NewRequest("GET", "/api/v1/master/depo", nil)
	req2.Header.Set("Authorization", "Bearer "+tkn)
	rr2 := httptest.NewRecorder()
	loggedHandler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr2.Code)
	}

	req3 := httptest.NewRequest("OPTIONS", "/api/v1/pemeriksaan/Semua/pasien/test", nil)
	rr3 := httptest.NewRecorder()
	loggedHandler.ServeHTTP(rr3, req3)

	if rr3.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr3.Code)
	}

	defaultHandler := middleware.LoggingMiddleware(dummyHandler)
	req4 := httptest.NewRequest("GET", "/health", nil)
	rr4 := httptest.NewRecorder()
	defaultHandler.ServeHTTP(rr4, req4)
	if rr4.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr4.Code)
	}
}

func TestLoggingMiddleware_StatusCodesAndOutputFormat(t *testing.T) {
	var buf bytes.Buffer
	testLogger := logger.NewWithOptions(&buf, "json", slog.LevelDebug, "erm-dokter")

	tests := []struct {
		name          string
		returnCode    int
		expectedLevel string
	}{
		{"200 OK -> INFO", http.StatusOK, "INFO"},
		{"404 Not Found -> WARN", http.StatusNotFound, "WARN"},
		{"500 Internal Error -> ERROR", http.StatusInternalServerError, "ERROR"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.returnCode)
			})
			wrapped := middleware.LoggingMiddlewareWithLogger(testLogger)(handler)

			req := httptest.NewRequest("POST", "/api/v1/test?filter=aktif", nil)
			req.Header.Set("X-Forwarded-For", "203.0.113.195, 198.51.100.1")
			req.Header.Set("User-Agent", "GoTestClient/1.0")

			rr := httptest.NewRecorder()
			wrapped.ServeHTTP(rr, req)

			var parsed map[string]any
			if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
				t.Fatalf("Failed to parse JSON: %v, raw: %s", err, buf.String())
			}

			if parsed["level"] != tc.expectedLevel {
				t.Errorf("expected level %s, got %v", tc.expectedLevel, parsed["level"])
			}
			if parsed["method"] != "POST" {
				t.Errorf("expected method POST, got %v", parsed["method"])
			}
			if parsed["path"] != "/api/v1/test" {
				t.Errorf("expected path /api/v1/test, got %v", parsed["path"])
			}
			if parsed["query"] != "filter=aktif" {
				t.Errorf("expected query filter=aktif, got %v", parsed["query"])
			}
			if parsed["client_ip"] != "203.0.113.195" {
				t.Errorf("expected client_ip 203.0.113.195, got %v", parsed["client_ip"])
			}
			if parsed["status"] != float64(tc.returnCode) {
				t.Errorf("expected status %d, got %v", tc.returnCode, parsed["status"])
			}
			if parsed["user_agent"] != "GoTestClient/1.0" {
				t.Errorf("expected user_agent GoTestClient/1.0, got %v", parsed["user_agent"])
			}
		})
	}
}

func TestLoggingMiddleware_WithErrorDetailAndUser(t *testing.T) {
	var buf bytes.Buffer
	testLogger := logger.NewWithOptions(&buf, "json", slog.LevelDebug, "erm-dokter")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-User-ID", "DRBUDI")
		w.Header().Set("X-Error-Detail", "Autentikasi ditolak: Kredensial login tidak cocok untuk username 'DRBUDI'")
		w.WriteHeader(http.StatusUnauthorized)
	})

	wrapped := middleware.LoggingMiddlewareWithLogger(testLogger)(handler)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}

	if rr.Header().Get("X-Error-Detail") != "" {
		t.Error("expected X-Error-Detail to be deleted from response")
	}
	if rr.Header().Get("X-User-ID") != "" {
		t.Error("expected X-User-ID to be deleted from response")
	}
	if rr.Header().Get("X-Logging-Middleware") != "" {
		t.Error("expected X-Logging-Middleware to be deleted from response")
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Failed to parse log JSON: %v, raw: %s", err, buf.String())
	}

	if parsed["level"] != "WARN" {
		t.Errorf("expected level WARN, got %v", parsed["level"])
	}
	expectedMsg := "Autentikasi ditolak: Kredensial login tidak cocok untuk username 'DRBUDI'"
	if parsed["msg"] != expectedMsg {
		t.Errorf("expected msg '%s', got '%v'", expectedMsg, parsed["msg"])
	}
	if parsed["user"] != "DRBUDI" {
		t.Errorf("expected user 'DRBUDI', got '%v'", parsed["user"])
	}
	if parsed["status"] != float64(401) {
		t.Errorf("expected status 401, got %v", parsed["status"])
	}
}

func TestLoggingMiddleware_WithSuccessMessage(t *testing.T) {
	var buf bytes.Buffer
	testLogger := logger.NewWithOptions(&buf, "json", slog.LevelDebug, "erm-dokter")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-User-ID", "DRBUDI")
		response.Success(w, "Login berhasil", map[string]string{"token": "dummy-token"})
	})

	wrapped := middleware.LoggingMiddlewareWithLogger(testLogger)(handler)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if rr.Header().Get("X-Log-Message") != "" {
		t.Error("expected X-Log-Message to be deleted from client response")
	}
	if rr.Header().Get("X-User-ID") != "" {
		t.Error("expected X-User-ID to be deleted from client response")
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Failed to parse log JSON: %v, raw: %s", err, buf.String())
	}

	if parsed["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", parsed["level"])
	}
	if parsed["msg"] != "Login berhasil" {
		t.Errorf("expected msg 'Login berhasil', got '%v'", parsed["msg"])
	}
	if parsed["user"] != "DRBUDI" {
		t.Errorf("expected user 'DRBUDI', got '%v'", parsed["user"])
	}
	if parsed["status"] != float64(200) {
		t.Errorf("expected status 200, got %v", parsed["status"])
	}
}
