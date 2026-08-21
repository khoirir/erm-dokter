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

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(msg string) error {
	return &UnauthorizedError{Message: msg}
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{Message: msg}
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewForbiddenError(msg string) error {
	return &ForbiddenError{Message: msg}
}

var log *logger.Logger

func SetLogger(l *logger.Logger) {
	log = l
}

func HandleError(w http.ResponseWriter, err error) {
	var validationErr ValidationError
	var businessErr *BusinessError
	var unauthorizedErr *UnauthorizedError
	var notFoundErr *NotFoundError
	var forbiddenErr *ForbiddenError

	switch {
	case errors.As(err, &validationErr):
		response.Error(w, http.StatusBadRequest, "Validasi gagal", validationErr)
	case errors.As(err, &businessErr):
		response.Error(w, http.StatusBadRequest, businessErr.Message, nil)
	case errors.As(err, &unauthorizedErr):
		response.Error(w, http.StatusUnauthorized, unauthorizedErr.Message, nil)
	case errors.As(err, &forbiddenErr):
		response.Error(w, http.StatusForbidden, forbiddenErr.Message, nil)
	case errors.As(err, &notFoundErr):
		response.Error(w, http.StatusNotFound, notFoundErr.Message, nil)
	default:
		if log != nil {
			log.Error("Internal server error: %v", err)
		}
		response.Error(w, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
	}
}
