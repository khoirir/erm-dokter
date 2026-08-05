package handler

import (
	"net/http"

	"erm-dokter/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "healthy",
		"service": "erm-dokter",
	}
	response.Success(w, "Service is up and running", data)
}
