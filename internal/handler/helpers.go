package handler

import (
	"errors"
	"net/http"

	"erm-dokter/internal/domain"
	"erm-dokter/internal/dto"
	"erm-dokter/pkg/response"
)

func handleError(w http.ResponseWriter, err error) {
	var validationErr dto.ValidationError
	var businessErr *domain.BusinessError

	switch {
	case errors.As(err, &validationErr):
		response.Error(w, http.StatusBadRequest, "Validasi gagal", validationErr)
	case errors.As(err, &businessErr):
		response.Error(w, http.StatusBadRequest, businessErr.Message, nil)
	default:
		response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}
}

