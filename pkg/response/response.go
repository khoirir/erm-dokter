package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, code int, success bool, message string, data any, meta any, err any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Success: success,
		Message: message,
		Data:    data,
		Meta:    meta,
		Error:   err,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusOK, true, message, data, nil, nil)
}

func SuccessWithMeta(w http.ResponseWriter, message string, data any, meta any) {
	JSON(w, http.StatusOK, true, message, data, meta, nil)
}

func Error(w http.ResponseWriter, code int, message string, err any) {
	JSON(w, code, false, message, nil, nil, err)
}
