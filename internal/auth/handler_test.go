package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"erm-dokter/internal/auth"
	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/token"
	"erm-dokter/internal/shared/apperror"
)

type mockAuthService struct {
	loginFn func(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error)
}

func (m *mockAuthService) Login(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, req)
	}
	return nil, nil
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		loginFn: func(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error) {
			if req.Username != "DR01" || req.Password != "rahasia" {
				return nil, apperror.NewUnauthorizedError("Username atau password salah")
			}
			return &auth.LoginResponse{
				Token:      "jwt-token-abc",
				KodeDokter: "DR01",
				NamaDokter: "dr. Ahmad",
			}, nil
		},
	}

	handler := auth.NewHandler(mockSvc)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, noOpMw, noOpMw, noOpMw)

	body, _ := json.Marshal(auth.LoginRequest{
		Username: "DR01",
		Password: "rahasia",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Success bool               `json:"success"`
		Message string             `json:"message"`
		Data    auth.LoginResponse `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Success || resp.Data.Token != "jwt-token-abc" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestAuthHandler_Login_ValidationAndWrongPassword(t *testing.T) {
	mockSvc := &mockAuthService{
		loginFn: func(ctx context.Context, req auth.LoginRequest) (*auth.LoginResponse, error) {
			return nil, apperror.NewUnauthorizedError("Username atau password salah")
		},
	}

	handler := auth.NewHandler(mockSvc)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, noOpMw, noOpMw, noOpMw)

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte("{invalid-json")))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid json, got %d", rr.Code)
		}

		var resp struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Message != "Data login tidak valid" {
			t.Errorf("Expected message 'Data login tidak valid', got '%s'", resp.Message)
		}
	})

	t.Run("Invalid JSON Payload with Password Masking", func(t *testing.T) {
		rawBody := `{"username":"DR01","password":"supersecretpassword",}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(rawBody)))
		rr := httptest.NewRecorder()
		rr.Header().Set("X-Logging-Middleware", "true")
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid json, got %d", rr.Code)
		}

		errDetail := rr.Header().Get("X-Error-Detail")
		if strings.Contains(errDetail, "supersecretpassword") {
			t.Errorf("Password was NOT masked in error detail! Got: %s", errDetail)
		}
		if !strings.Contains(errDetail, `"password":"***"`) {
			t.Errorf("Expected masked password in error detail, got: %s", errDetail)
		}
	})

	t.Run("Validation Error", func(t *testing.T) {
		body, _ := json.Marshal(auth.LoginRequest{Username: "", Password: ""})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for empty login, got %d", rr.Code)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		body, _ := json.Marshal(auth.LoginRequest{Username: "DR01", Password: "wrong"})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for wrong password, got %d", rr.Code)
		}
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	mockSvc := &mockAuthService{}
	handler := auth.NewHandler(mockSvc)
	mux := http.NewServeMux()

	authMw := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), middleware.UserClaimKey, &token.Claims{KodeDokter: "DR01"})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }

	handler.RegisterRoutes(mux, noOpMw, authMw, noOpMw)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}
