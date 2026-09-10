package shared_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"erm-dokter/internal/shared"
)

type dummyPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestDecodeJSON_Success(t *testing.T) {
	raw := `{"username":"dokter1","password":"secretpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(raw)))
	rr := httptest.NewRecorder()

	var payload dummyPayload
	ok := shared.DecodeJSON(rr, req, &payload)
	if !ok {
		t.Fatal("expected DecodeJSON to return true for valid payload")
	}

	if payload.Username != "dokter1" || payload.Password != "secretpassword" {
		t.Errorf("unexpected parsed payload: %+v", payload)
	}
}

func TestDecodeJSON_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("   ")))
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Logging-Middleware", "true")

	var payload dummyPayload
	ok := shared.DecodeJSON(rr, req, &payload, "Data login tidak valid")
	if ok {
		t.Fatal("expected DecodeJSON to return false for empty body")
	}

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp.Message != "Data login tidak valid" {
		t.Errorf("expected client message 'Data login tidak valid', got '%s'", resp.Message)
	}

	errDetail := rr.Header().Get("X-Error-Detail")
	if !strings.Contains(errDetail, "Payload request kosong") {
		t.Errorf("expected error detail to mention empty payload, got: %s", errDetail)
	}
}

func TestDecodeJSON_MalformedJSON_WithPasswordMasking(t *testing.T) {
	raw := `{"username":"dokter1","password":"secretpassword",}`
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(raw)))
	rr := httptest.NewRecorder()
	rr.Header().Set("X-Logging-Middleware", "true")

	var payload dummyPayload
	ok := shared.DecodeJSON(rr, req, &payload, "Data login tidak valid")
	if ok {
		t.Fatal("expected DecodeJSON to return false for malformed json")
	}

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp.Message != "Data login tidak valid" {
		t.Errorf("expected client message 'Data login tidak valid', got '%s'", resp.Message)
	}

	errDetail := rr.Header().Get("X-Error-Detail")
	if strings.Contains(errDetail, "secretpassword") {
		t.Errorf("password leaked in error detail! Got: %s", errDetail)
	}
	if !strings.Contains(errDetail, `"password":"***"`) {
		t.Errorf("expected masked password in error detail, got: %s", errDetail)
	}
}

func TestMaskSensitiveData(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `{"password":"rahasia123"}`,
			expected: `{"password":"***"}`,
		},
		{
			input:    `{"token":"bearer-abc-123"}`,
			expected: `{"token":"***"}`,
		},
		{
			input:    `{"pin":"123456"}`,
			expected: `{"pin":"***"}`,
		},
		{
			input:    `{"api_key":"secret-api-key"}`,
			expected: `{"api_key":"***"}`,
		},
	}

	for _, tc := range tests {
		result := shared.MaskSensitiveData(tc.input)
		if result != tc.expected {
			t.Errorf("MaskSensitiveData(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}
