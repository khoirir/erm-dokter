package apperror

import (
	"errors"
	"net/http"
	"strings"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/pkg/response"
)

type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(msg string) error {
	return &BusinessError{Message: msg}
}

type ValidationError map[string]string

func (v ValidationError) Error() string {
	var errs []string
	for field, msg := range v {
		errs = append(errs, field+": "+msg)
	}
	return strings.Join(errs, "; ")
}

var log *logger.Logger

func SetLogger(l *logger.Logger) {
	log = l
}

func HandleError(w http.ResponseWriter, err error) {
	var validationErr ValidationError
	var businessErr *BusinessError

	switch {
	case errors.As(err, &validationErr):
		response.Error(w, http.StatusBadRequest, "Validasi gagal", validationErr)
	case errors.As(err, &businessErr):
		response.Error(w, http.StatusBadRequest, businessErr.Message, nil)
	default:
		if log != nil {
			log.Error("Internal server error: %v", err)
		}
		response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}
}
