package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/token"
)

const testSecret = "my-secret-key-12345678901234567890"

func TestJWTMiddleware(t *testing.T) {
	jwtMw := middleware.JWTMiddleware(testSecret)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		kodeDokter, err := middleware.GetKodeDokter(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.Write([]byte("OK: " + kodeDokter))
	}

	wrapped := jwtMw(dummyHandler)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})

	t.Run("Invalid Bearer Prefix", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Basic abcdef")
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})

	t.Run("Invalid Token String", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.payload")
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		validToken, err := token.GenerateToken("DR001", "dr. Ahmad", testSecret, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		if rr.Body.String() != "OK: DR001" {
			t.Errorf("Unexpected body: %s", rr.Body.String())
		}
	})
}

func TestGetKodeDokter_NoClaims(t *testing.T) {
	_, err := middleware.GetKodeDokter(context.Background())
	if err == nil {
		t.Error("Expected error when context has no claims, got nil")
	}

	ctxWithEmptyDoc := context.WithValue(context.Background(), middleware.UserClaimKey, &token.Claims{
		KodeDokter: "   ",
	})
	_, errEmpty := middleware.GetKodeDokter(ctxWithEmptyDoc)
	if errEmpty == nil {
		t.Error("Expected error when KodeDokter is empty, got nil")
	}
}

func TestCORSMiddleware(t *testing.T) {
	cors := middleware.CORSMiddleware("http://localhost:3000")
	handler := cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("cors passed"))
	}))

	t.Run("Preflight OPTIONS", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status 204 No Content for OPTIONS, got %d", rr.Code)
		}
		if rr.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("Unexpected CORS origin: %s", rr.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("GET Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}

func TestRateLimiter(t *testing.T) {
	mw := middleware.RateLimitMiddleware(2, 500*time.Millisecond)

	handler := mw(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.50:12345"

	// Req 1: Allowed
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req)
	if rr1.Code != http.StatusOK {
		t.Errorf("Req 1: expected 200, got %d", rr1.Code)
	}

	// Req 2: Allowed
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusOK {
		t.Errorf("Req 2: expected 200, got %d", rr2.Code)
	}

	// Req 3: Rate Limited (429)
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req)
	if rr3.Code != http.StatusTooManyRequests {
		t.Errorf("Req 3: expected 429 Too Many Requests, got %d", rr3.Code)
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	timeoutMw := middleware.TimeoutMiddleware(50 * time.Millisecond)

	t.Run("Within Timeout", func(t *testing.T) {
		handler := timeoutMw(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("fast response"))
		})

		req := httptest.NewRequest(http.MethodGet, "/fast", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Exceeded Timeout", func(t *testing.T) {
		handler := timeoutMw(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusGatewayTimeout {
			t.Errorf("Expected status 504 Gateway Timeout, got %d", rr.Code)
		}
	})
}
