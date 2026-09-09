package master

import (
	"net/http"

	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/master/penjamin", authMiddleware(timeoutMiddleware(h.DaftarPenjamin)))
	mux.HandleFunc("GET /api/v1/master/depo", authMiddleware(timeoutMiddleware(h.DaftarDepo)))
	mux.HandleFunc("GET /api/v1/master/poliklinik", authMiddleware(timeoutMiddleware(h.DaftarPoliklinik)))
	mux.HandleFunc("GET /api/v1/master/bangsal", authMiddleware(timeoutMiddleware(h.DaftarBangsal)))
	mux.HandleFunc("GET /api/v1/master/kelas", authMiddleware(timeoutMiddleware(h.DaftarKelas)))
}

func (h *Handler) DaftarPenjamin(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarPenjamin(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar penjamin", data)
}

func (h *Handler) DaftarDepo(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarDepo(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar depo", data)
}

func (h *Handler) DaftarPoliklinik(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarPoliklinik(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar poliklinik", data)
}

func (h *Handler) DaftarBangsal(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarBangsal(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar bangsal", data)
}

func (h *Handler) DaftarKelas(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarKelas(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar kelas kamar", data)
}

