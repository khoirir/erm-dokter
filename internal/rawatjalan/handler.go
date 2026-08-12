package rawatjalan

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	rawatJalanService Service
	encryptionKey     string
}

func NewHandler(service Service, encryptionKey string) *Handler {
	return &Handler{
		rawatJalanService: service,
		encryptionKey:     encryptionKey,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/rawat-jalan/referensi-filter", authMiddleware(timeoutMiddleware(h.GetReferensiFilter)))
	mux.HandleFunc("POST /api/v1/rawat-jalan/antrean", authMiddleware(timeoutMiddleware(h.DaftarAntreanDokter)))
	mux.HandleFunc("GET /api/v1/rawat-jalan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailKunjungan)))
}

func (h *Handler) DaftarAntreanDokter(w http.ResponseWriter, r *http.Request) {
	var filter FilterAntreanDokter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil && err.Error() != "EOF" {
		response.Error(w, http.StatusBadRequest, "Format request JSON tidak valid", nil)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	daftarAntrean, meta, err := h.rawatJalanService.DaftarAntreanDokter(r.Context(), kodeDokter, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range daftarAntrean {
		kunjungan := &daftarAntrean[i]
		if encryptedId, err := crypto.Encrypt(kunjungan.NoRawat, h.encryptionKey); err == nil {
			kunjungan.Id = encryptedId
		}
		if encryptedIdPasien, err := crypto.Encrypt(kunjungan.NoRekamMedis, h.encryptionKey); err == nil {
			kunjungan.IdPasien = encryptedIdPasien
		}
		kunjungan.NoRekamMedis = kunjungan.FormatNoRekamMedis()
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar antrean dokter", daftarAntrean, meta)
}

func (h *Handler) DetailKunjungan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	if idKunjungan == "" {
		response.Error(w, http.StatusBadRequest, "ID kunjungan tidak ditemukan", nil)
		return
	}

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID kunjungan tidak valid atau kadaluarsa", nil)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	kunjungan, err := h.rawatJalanService.DetailKunjungan(r.Context(), noRawat, kodeDokter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if kunjungan == nil {
		response.Error(w, http.StatusNotFound, "Detail kunjungan pasien tidak ditemukan", nil)
		return
	}
	kunjungan.Id = idKunjungan
	if encryptedIdPasien, err := crypto.Encrypt(kunjungan.NoRekamMedis, h.encryptionKey); err == nil {
		kunjungan.IdPasien = encryptedIdPasien
	}
	kunjungan.NoRekamMedis = kunjungan.FormatNoRekamMedis()

	response.Success(w, "Berhasil mengambil detail kunjungan pasien", kunjungan)
}

func (h *Handler) GetReferensiFilter(w http.ResponseWriter, r *http.Request) {
	referensi := h.rawatJalanService.GetReferensiFilter(r.Context())
	response.Success(w, "Berhasil mengambil referensi filter", referensi)
}
