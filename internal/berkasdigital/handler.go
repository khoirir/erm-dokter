package berkasdigital

import (
	"net/http"
	"strings"

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
	mux.HandleFunc("GET /api/v1/berkas-digital/{id_berkas}", authMiddleware(h.StreamBerkasDigital))
	mux.HandleFunc("GET /api/v1/master/berkas-digital", authMiddleware(timeoutMiddleware(h.MasterBerkasDigital)))
}

func (h *Handler) StreamBerkasDigital(w http.ResponseWriter, r *http.Request) {
	idBerkas := strings.TrimSpace(r.PathValue("id_berkas"))
	if idBerkas == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID berkas digital wajib diisi"))
		return
	}

	err := h.service.StreamBerkasDigital(r.Context(), idBerkas, w)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
}

func (h *Handler) MasterBerkasDigital(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.GetMasterBerkas(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil master berkas digital", list)
}
