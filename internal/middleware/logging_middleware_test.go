package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erm-dokter/internal/pkg/token"
)

func TestLoggingMiddleware_AuthenticatedAndAnonymous(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	loggedHandler := LoggingMiddleware(dummyHandler)

	// 1. Anonymous request
	req1 := httptest.NewRequest("GET", "/health", nil)
	rr1 := httptest.NewRecorder()
	loggedHandler.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr1.Code)
	}

	// 2. Authenticated request with Bearer token
	secret := "test-secret-key-1234567890123456"
	tkn, err := token.GenerateToken("DR001", "dr. Handi", secret, 1*time.Hour)
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

	// 3. OPTIONS preflight request (should bypass logger without error)
	req3 := httptest.NewRequest("OPTIONS", "/api/v1/pemeriksaan/pasien/test/Semua", nil)
	rr3 := httptest.NewRecorder()
	loggedHandler.ServeHTTP(rr3, req3)

	if rr3.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr3.Code)
	}
}

