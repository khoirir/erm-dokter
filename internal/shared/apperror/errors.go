package apperror

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/pkg/response"
)

type BusinessError struct {
	Message   string
	LogDetail string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func (e *BusinessError) LogMessage() string {
	if e.LogDetail != "" {
		return e.LogDetail
	}
	return e.Message
}

func NewBusinessError(msg string, logDetail ...string) *BusinessError {
	detail := ""
	if len(logDetail) > 0 {
		detail = logDetail[0]
	}
	return &BusinessError{Message: msg, LogDetail: detail}
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
	Message   string
	LogDetail string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) LogMessage() string {
	if e.LogDetail != "" {
		return e.LogDetail
	}
	return e.Message
}

func NewUnauthorizedError(msg string, logDetail ...string) *UnauthorizedError {
	detail := ""
	if len(logDetail) > 0 {
		detail = logDetail[0]
	}
	return &UnauthorizedError{Message: msg, LogDetail: detail}
}

type NotFoundError struct {
	Message   string
	LogDetail string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) LogMessage() string {
	if e.LogDetail != "" {
		return e.LogDetail
	}
	return e.Message
}

func NewNotFoundError(msg string, logDetail ...string) *NotFoundError {
	detail := ""
	if len(logDetail) > 0 {
		detail = logDetail[0]
	}
	return &NotFoundError{Message: msg, LogDetail: detail}
}

type ForbiddenError struct {
	Message   string
	LogDetail string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func (e *ForbiddenError) LogMessage() string {
	if e.LogDetail != "" {
		return e.LogDetail
	}
	return e.Message
}

func NewForbiddenError(msg string, logDetail ...string) *ForbiddenError {
	detail := ""
	if len(logDetail) > 0 {
		detail = logDetail[0]
	}
	return &ForbiddenError{Message: msg, LogDetail: detail}
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

	var logMsg string
	var statusCode int
	var respMsg string
	var errData any

	switch {
	case errors.As(err, &validationErr):
		logMsg = fmt.Sprintf("Validasi input gagal: %v", validationErr)
		statusCode = http.StatusBadRequest
		respMsg = "Validasi gagal"
		errData = validationErr
	case errors.As(err, &businessErr):
		logMsg = fmt.Sprintf("Aturan bisnis ditolak: %s", businessErr.LogMessage())
		statusCode = http.StatusBadRequest
		respMsg = businessErr.Message
	case errors.As(err, &unauthorizedErr):
		logMsg = fmt.Sprintf("Autentikasi ditolak: %s", unauthorizedErr.LogMessage())
		statusCode = http.StatusUnauthorized
		respMsg = unauthorizedErr.Message
	case errors.As(err, &forbiddenErr):
		logMsg = fmt.Sprintf("Akses terlarang: %s", forbiddenErr.LogMessage())
		statusCode = http.StatusForbidden
		respMsg = forbiddenErr.Message
	case errors.As(err, &notFoundErr):
		logMsg = fmt.Sprintf("Data tidak ditemukan: %s", notFoundErr.LogMessage())
		statusCode = http.StatusNotFound
		respMsg = notFoundErr.Message
	default:
		logMsg = fmt.Sprintf("Internal server error: %v", err)
		statusCode = http.StatusInternalServerError
		respMsg = "Terjadi kesalahan pada server"
	}

	isUnderLoggingMiddleware := w.Header().Get("X-Logging-Middleware") == "true"
	if isUnderLoggingMiddleware {
		w.Header().Set("X-Error-Detail", logMsg)
	} else {
		targetLog := log
		if targetLog != nil {
			if reqID := w.Header().Get("X-Request-ID"); reqID != "" {
				targetLog = targetLog.With("request_id", reqID)
			}
			if user := w.Header().Get("X-User-ID"); user != "" {
				targetLog = targetLog.With("user", user)
			}
			if method := w.Header().Get("X-Request-Method"); method != "" {
				targetLog = targetLog.With("method", method)
			}
			if path := w.Header().Get("X-Request-Path"); path != "" {
				targetLog = targetLog.With("path", path)
			}

			if statusCode >= 500 {
				targetLog.Error(logMsg)
			} else {
				targetLog.Warn(logMsg)
			}
		}

		w.Header().Del("X-User-ID")
		w.Header().Del("X-Request-Method")
		w.Header().Del("X-Request-Path")
	}

	response.Error(w, statusCode, respMsg, errData)
}
