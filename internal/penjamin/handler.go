package penjamin

import (
	"net/http"

	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	penjaminService Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		penjaminService: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/penjamin", authMiddleware(timeoutMiddleware(h.DaftarPenjamin)))
}

func (h *Handler) DaftarPenjamin(w http.ResponseWriter, r *http.Request) {
	daftarPenjamin, err := h.penjaminService.DaftarPenjamin(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil daftar penjamin", daftarPenjamin)
}
