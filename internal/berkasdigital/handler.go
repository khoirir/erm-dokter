package berkasdigital

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	service       Service
	encryptionKey string
}

func NewHandler(service Service, encryptionKey string) *Handler {
	return &Handler{
		service:       service,
		encryptionKey: encryptionKey,
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

	targetURL, err := crypto.Decrypt(idBerkas, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Token berkas digital tidak valid"))
		return
	}

	stream, err := h.service.GetBerkasStream(r.Context(), targetURL)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	defer stream.Body.Close()

	w.Header().Set("Content-Type", stream.ContentType)
	if stream.ContentLength != "" {
		w.Header().Set("Content-Length", stream.ContentLength)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", stream.Filename))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, stream.Body)
}

func (h *Handler) MasterBerkasDigital(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.GetMasterBerkas(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil master berkas digital", list)
}

