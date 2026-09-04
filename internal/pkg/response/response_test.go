package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/pkg/response"
)

func TestResponseHelpers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		rr := httptest.NewRecorder()
		response.Success(rr, "Data berhasil diambil", map[string]string{"foo": "bar"})

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp response.Response
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if !resp.Success || resp.Message != "Data berhasil diambil" {
			t.Errorf("Unexpected response: %+v", resp)
		}
	})

	t.Run("Created", func(t *testing.T) {
		rr := httptest.NewRecorder()
		response.Created(rr, "Data berhasil dibuat", map[string]string{"id": "123"})

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rr.Code)
		}
	})

	t.Run("SuccessWithMeta", func(t *testing.T) {
		rr := httptest.NewRecorder()
		meta := map[string]int{"page": 1, "total": 10}
		response.SuccessWithMeta(rr, "Daftar data", []string{"a", "b"}, meta)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var resp response.Response
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Meta == nil {
			t.Error("Expected meta in response")
		}
	})

	t.Run("Error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		response.Error(rr, http.StatusBadRequest, "Permintaan tidak valid", map[string]string{"id": "wajib diisi"})

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}

		var resp response.Response
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Success {
			t.Error("Expected success=false for Error helper")
		}
		if resp.Error == nil {
			t.Error("Expected error field in response")
		}
	})
}
