package health

import (
	"net/http"

	"erm-dokter/internal/pkg/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.HealthCheck)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "healthy",
		"service": "erm-dokter",
	}
	response.Success(w, "Service is up and running", data)
}
