package response

import (
	"encoding/json"
	"net/http"
)

// Response standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// JSON sends a JSON response with status code
func JSON(w http.ResponseWriter, code int, success bool, message string, data interface{}, meta interface{}, err interface{}) {
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

// Success sends 200 OK success response without metadata
func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, true, message, data, nil, nil)
}

// SuccessWithMeta sends 200 OK success response with pagination metadata
func SuccessWithMeta(w http.ResponseWriter, message string, data interface{}, meta interface{}) {
	JSON(w, http.StatusOK, true, message, data, meta, nil)
}

// Error sends error response
func Error(w http.ResponseWriter, code int, message string, err interface{}) {
	JSON(w, code, false, message, nil, nil, err)
}
