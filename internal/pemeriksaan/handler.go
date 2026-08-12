package pemeriksaan

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	pemeriksaanService Service
	encryptionKey      string
}

func NewHandler(service Service, encryptionKey string) *Handler {
	return &Handler{
		pemeriksaanService: service,
		encryptionKey:      encryptionKey,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("POST /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPemeriksaan)))
	mux.HandleFunc("POST /api/v1/pemeriksaan/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPemeriksaanByPasien)))
	mux.HandleFunc("GET /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}/{id_pemeriksaan}", authMiddleware(timeoutMiddleware(h.DetailPemeriksaan)))
}

func (h *Handler) DaftarPemeriksaan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	var filter FilterDaftarPemeriksaan
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil && err.Error() != "EOF" {
		response.Error(w, http.StatusBadRequest, "Format request JSON tidak valid", nil)
		return
	}

	daftarPemeriksaan, meta, err := h.pemeriksaanService.DaftarPemeriksaan(r.Context(), noRawat, shared.StatusLanjut(statusLanjut), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range daftarPemeriksaan {
		pemeriksaan := &daftarPemeriksaan[i]
		if encrypted, err := crypto.Encrypt(pemeriksaan.CompositeKey(), h.encryptionKey); err == nil {
			pemeriksaan.Id = encrypted
		}
		pemeriksaan.IdKunjungan = idKunjungan
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar pemeriksaan", daftarPemeriksaan, meta)
}

func (h *Handler) DaftarPemeriksaanByPasien(w http.ResponseWriter, r *http.Request) {
	idPasien := r.PathValue("id_pasien")
	statusLanjut := r.PathValue("status_lanjut")

	noRekamMedis, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	var filter FilterDaftarPemeriksaan
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil && err.Error() != "EOF" {
		response.Error(w, http.StatusBadRequest, "Format request JSON tidak valid", nil)
		return
	}

	daftarPemeriksaan, meta, err := h.pemeriksaanService.DaftarPemeriksaanByRM(r.Context(), noRekamMedis, shared.StatusLanjut(statusLanjut), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range daftarPemeriksaan {
		pemeriksaan := &daftarPemeriksaan[i]
		if encrypted, err := crypto.Encrypt(pemeriksaan.CompositeKey(), h.encryptionKey); err == nil {
			pemeriksaan.Id = encrypted
		}
		if encryptedIdKunjungan, err := crypto.Encrypt(pemeriksaan.NoRawat, h.encryptionKey); err == nil {
			pemeriksaan.IdKunjungan = encryptedIdKunjungan
		}
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat pemeriksaan pasien", daftarPemeriksaan, meta)
}

func (h *Handler) DetailPemeriksaan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")
	encryptedID := r.PathValue("id_pemeriksaan")

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID kunjungan tidak valid atau kadaluarsa", nil)
		return
	}

	decrypted, err := crypto.Decrypt(encryptedID, h.encryptionKey)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID pemeriksaan tidak valid atau kadaluarsa", nil)
		return
	}

	idPemeriksaan, err := ParseIdPemeriksaan(decrypted)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if idPemeriksaan.NoRawat != noRawat {
		response.Error(w, http.StatusBadRequest, "ID pemeriksaan tidak cocok dengan ID kunjungan", nil)
		return
	}

	pemeriksaan, err := h.pemeriksaanService.DetailPemeriksaan(r.Context(), idPemeriksaan, shared.StatusLanjut(statusLanjut))
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if pemeriksaan == nil {
		response.Error(w, http.StatusNotFound, "Detail pemeriksaan tidak ditemukan", nil)
		return
	}

	pemeriksaan.Id = encryptedID
	pemeriksaan.IdKunjungan = idKunjungan

	response.Success(w, "Berhasil mengambil detail pemeriksaan", pemeriksaan)
}
