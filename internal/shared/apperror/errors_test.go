package apperror_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared/apperror"
)

func TestAppErrors(t *testing.T) {
	apperror.SetLogger(logger.New())

	tests := []struct {
		name         string
		err          error
		expectedCode int
		checkMsg     string
	}{
		{
			name:         "BusinessError",
			err:          apperror.NewBusinessError("Stok obat tidak mencukupi"),
			expectedCode: http.StatusBadRequest,
			checkMsg:     "Stok obat tidak mencukupi",
		},
		{
			name: "ValidationError",
			err: apperror.ValidationError{
				"nama": "wajib diisi",
			},
			expectedCode: http.StatusBadRequest,
			checkMsg:     "Validasi gagal",
		},
		{
			name:         "UnauthorizedError",
			err:          apperror.NewUnauthorizedError("Token kedaluwarsa"),
			expectedCode: http.StatusUnauthorized,
			checkMsg:     "Token kedaluwarsa",
		},
		{
			name:         "ForbiddenError",
			err:          apperror.NewForbiddenError("Bukan dokter pemeriksa"),
			expectedCode: http.StatusForbidden,
			checkMsg:     "Bukan dokter pemeriksa",
		},
		{
			name:         "NotFoundError",
			err:          apperror.NewNotFoundError("Data pasien tidak ditemukan"),
			expectedCode: http.StatusNotFound,
			checkMsg:     "Data pasien tidak ditemukan",
		},
		{
			name:         "InternalServerError (Generic)",
			err:          errors.New("unexpected database connection dropped"),
			expectedCode: http.StatusInternalServerError,
			checkMsg:     "Terjadi kesalahan pada server",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() == "" {
				t.Error("Expected non-empty error message")
			}

			rr := httptest.NewRecorder()
			apperror.HandleError(rr, tc.err)

			if rr.Code != tc.expectedCode {
				t.Errorf("Expected status %d, got %d", tc.expectedCode, rr.Code)
			}

			var resp struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal JSON response: %v", err)
			}

			if resp.Success {
				t.Errorf("Expected success to be false, got true")
			}
			if resp.Message != tc.checkMsg {
				t.Errorf("Expected message '%s', got '%s'", tc.checkMsg, resp.Message)
			}
		})
	}
}
