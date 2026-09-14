package middleware_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erm-dokter/internal/middleware"
)

func TestGetClientIP(t *testing.T) {
	t.Run("RemoteAddr with port", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.100:54321"

		ip := middleware.GetClientIP(req)
		if ip != "192.168.1.100" {
			t.Errorf("Expected 192.168.1.100, got %s", ip)
		}
	})

	t.Run("RemoteAddr IPv6 with port", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "[::1]:54321"

		ip := middleware.GetClientIP(req)
		if ip != "::1" {
			t.Errorf("Expected ::1, got %s", ip)
		}
	})

	t.Run("X-Real-IP", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.100:54321"
		req.Header.Set("X-Real-IP", "10.10.10.5")

		ip := middleware.GetClientIP(req)
		if ip != "10.10.10.5" {
			t.Errorf("Expected 10.10.10.5, got %s", ip)
		}
	})

	t.Run("X-Forwarded-For single", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.100:54321"
		req.Header.Set("X-Real-IP", "10.10.10.5")
		req.Header.Set("X-Forwarded-For", "172.16.0.2")

		ip := middleware.GetClientIP(req)
		if ip != "172.16.0.2" {
			t.Errorf("Expected 172.16.0.2, got %s", ip)
		}
	})

	t.Run("X-Forwarded-For multi-hop", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.100:54321"
		req.Header.Set("X-Forwarded-For", "172.16.0.2, 10.0.0.1, 192.168.1.1")

		ip := middleware.GetClientIP(req)
		if ip != "172.16.0.2" {
			t.Errorf("Expected 172.16.0.2, got %s", ip)
		}
	})
}

func TestLoginRateLimitMiddleware(t *testing.T) {
	t.Run("Same IP and same username hits limit", func(t *testing.T) {
		mw := middleware.LoginRateLimitMiddleware(2, 500*time.Millisecond, true)

		handler := mw(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		sendLoginReq := func() *httptest.ResponseRecorder {
			body := bytes.NewBufferString(`{"username":"DR001","password":"secret"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
			req.RemoteAddr = "192.168.30.1:12345"
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			return rr
		}

		// Req 1: Allowed
		if rr := sendLoginReq(); rr.Code != http.StatusOK {
			t.Fatalf("Req 1: expected 200, got %d", rr.Code)
		}
		// Req 2: Allowed
		if rr := sendLoginReq(); rr.Code != http.StatusOK {
			t.Fatalf("Req 2: expected 200, got %d", rr.Code)
		}
		// Req 3: Rate limited
		if rr := sendLoginReq(); rr.Code != http.StatusTooManyRequests {
			t.Fatalf("Req 3: expected 429 Too Many Requests, got %d", rr.Code)
		}
	})

	t.Run("Same IP but different usernames have independent quotas", func(t *testing.T) {
		mw := middleware.LoginRateLimitMiddleware(2, 500*time.Millisecond, true)

		handler := mw(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		sendLoginReq := func(username string) *httptest.ResponseRecorder {
			body := bytes.NewBufferString(`{"username":"` + username + `","password":"secret"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
			req.RemoteAddr = "192.168.30.1:12345" // Same gateway IP
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			return rr
		}

		// Doctor A exhausts 2 requests
		if rr := sendLoginReq("dr_ahmad"); rr.Code != http.StatusOK {
			t.Fatalf("Doctor A Req 1: expected 200, got %d", rr.Code)
		}
		if rr := sendLoginReq("dr_ahmad"); rr.Code != http.StatusOK {
			t.Fatalf("Doctor A Req 2: expected 200, got %d", rr.Code)
		}
		// Doctor A is blocked
		if rr := sendLoginReq("dr_ahmad"); rr.Code != http.StatusTooManyRequests {
			t.Fatalf("Doctor A Req 3: expected 429, got %d", rr.Code)
		}

		// Doctor B from SAME gateway IP is NOT blocked!
		if rr := sendLoginReq("dr_budi"); rr.Code != http.StatusOK {
			t.Fatalf("Doctor B Req 1: expected 200, got %d (should NOT be blocked by Doctor A)", rr.Code)
		}
		if rr := sendLoginReq("dr_budi"); rr.Code != http.StatusOK {
			t.Fatalf("Doctor B Req 2: expected 200, got %d", rr.Code)
		}
		// Doctor B now blocked after 2 requests
		if rr := sendLoginReq("dr_budi"); rr.Code != http.StatusTooManyRequests {
			t.Fatalf("Doctor B Req 3: expected 429, got %d", rr.Code)
		}
	})

	t.Run("Preserves r.Body for downstream handler", func(t *testing.T) {
		mw := middleware.LoginRateLimitMiddleware(5, 1*time.Minute, true)

		var capturedUsername string
		handler := mw(func(w http.ResponseWriter, r *http.Request) {
			var payload struct {
				Username string `json:"username"`
			}
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %v", err)
			}
			if err := json.Unmarshal(bodyBytes, &payload); err != nil {
				t.Fatalf("Failed to unmarshal body downstream: %v", err)
			}
			capturedUsername = payload.Username
			w.WriteHeader(http.StatusOK)
		})

		body := bytes.NewBufferString(`{"username":"DR007","password":"password123"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", rr.Code)
		}
		if capturedUsername != "DR007" {
			t.Fatalf("Expected downstream handler to receive DR007, got '%s'", capturedUsername)
		}
	})

	t.Run("Fallback to IP when username is missing", func(t *testing.T) {
		mw := middleware.LoginRateLimitMiddleware(2, 500*time.Millisecond, true)

		handler := mw(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		sendEmptyReq := func() *httptest.ResponseRecorder {
			body := bytes.NewBufferString(`{"password":"secret"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
			req.RemoteAddr = "10.0.0.99:9999"
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			return rr
		}

		// Req 1 & 2 allowed
		sendEmptyReq()
		sendEmptyReq()
		// Req 3 blocked
		if rr := sendEmptyReq(); rr.Code != http.StatusTooManyRequests {
			t.Fatalf("Expected 429 when username is empty, got %d", rr.Code)
		}
	})

	t.Run("Disabled rate limit bypasses throttling", func(t *testing.T) {
		mw := middleware.LoginRateLimitMiddleware(1, 500*time.Millisecond, false) // disabled

		handler := mw(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		sendReq := func() *httptest.ResponseRecorder {
			body := bytes.NewBufferString(`{"username":"dr_bypass","password":"123"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
			req.RemoteAddr = "10.0.0.1:1234"
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			return rr
		}

		// Even with rate=1, multiple requests succeed because disabled
		for i := 0; i < 5; i++ {
			if rr := sendReq(); rr.Code != http.StatusOK {
				t.Fatalf("Req %d: expected 200 when disabled, got %d", i+1, rr.Code)
			}
		}
	})
}
