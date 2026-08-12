package pemeriksaan

import (
	"fmt"
	"net/http"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	pemeriksaanService Service
	encryptionKey     string
}

func NewHandler(service Service, encryptionKey string) *Handler {
	return &Handler{
		pemeriksaanService: service,
		encryptionKey:     encryptionKey,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/pemeriksaan/{id}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPemeriksaan)))
}

func (h *Handler) DaftarPemeriksaan(w http.ResponseWriter, r *http.Request) {
	encryptedNoRawat := r.PathValue("id")
	statusLanjut := r.PathValue("status_lanjut")

	norawat, err := crypto.Decrypt(encryptedNoRawat, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	daftarPemeriksaan, err := h.pemeriksaanService.DaftarPemeriksaan(r.Context(), norawat, shared.StatusLanjut(statusLanjut))
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range daftarPemeriksaan {
		p := &daftarPemeriksaan[i]
		rawKey := fmt.Sprintf("%s~%s~%s", p.NoRawat, p.TanggalPemeriksaan, p.JamPemeriksaan)
		encrypted, err := crypto.Encrypt(rawKey, h.encryptionKey)
		if err == nil {
			p.Id = encrypted
		}
	}

	response.Success(w, "Berhasil mengambil daftar pemeriksaan", daftarPemeriksaan)
}