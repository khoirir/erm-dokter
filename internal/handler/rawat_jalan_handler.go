package handler

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/domain"
	"erm-dokter/internal/dto"
	"erm-dokter/internal/middleware"
	"erm-dokter/pkg/crypto"
	"erm-dokter/pkg/response"
)

type RawatJalanHandler struct {
	rawatJalanUsecase domain.RawatJalanUsecase
	encryptionKey     string
}

func NewRawatJalanHandler(usecase domain.RawatJalanUsecase, encryptionKey string) *RawatJalanHandler {
	return &RawatJalanHandler{
		rawatJalanUsecase: usecase,
		encryptionKey:     encryptionKey,
	}
}

func (h *RawatJalanHandler) DaftarAntreanDokter(w http.ResponseWriter, r *http.Request) {
	var filter dto.FilterAntreanDokter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil && err.Error() != "EOF" {
		response.Error(w, http.StatusBadRequest, "Format request JSON tidak valid", nil)
		return
	}

	claim := middleware.GetUserClaim(r.Context())
	if claim != nil && filter.KodeDokter == "" {
		filter.KodeDokter = claim.KodeDokter
	}

	daftarAntrean, meta, err := h.rawatJalanUsecase.DaftarAntreanDokter(r.Context(), filter)
	if err != nil {
		handleError(w, err)
		return
	}

	// Enkripsi no_rawat dan buat detail URL (transport concern)
	for i := range daftarAntrean {
		encrypted, err := crypto.Encrypt(daftarAntrean[i].NoRawat, h.encryptionKey)
		if err == nil {
			daftarAntrean[i].DetailURL = "/api/v1/rawat-jalan/detail/" + encrypted
		}
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar antrean dokter", daftarAntrean, meta)
}

func (h *RawatJalanHandler) DetailKunjungan(w http.ResponseWriter, r *http.Request) {
	encryptedNoRawat := r.PathValue("no_rawat")
	if encryptedNoRawat == "" {
		response.Error(w, http.StatusBadRequest, "Parameter no_rawat wajib diisi", nil)
		return
	}

	// Dekripsi no_rawat dari URL (transport concern)
	noRawat, err := crypto.Decrypt(encryptedNoRawat, h.encryptionKey)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Parameter no_rawat tidak valid", nil)
		return
	}

	claim := middleware.GetUserClaim(r.Context())
	kodeDokter := ""
	if claim != nil {
		kodeDokter = claim.KodeDokter
	}

	kunjungan, err := h.rawatJalanUsecase.DetailKunjungan(r.Context(), noRawat, kodeDokter)
	if err != nil {
		handleError(w, err)
		return
	}

	if kunjungan == nil {
		response.Error(w, http.StatusNotFound, "Detail kunjungan pasien tidak ditemukan", nil)
		return
	}

	response.Success(w, "Berhasil mengambil detail kunjungan", kunjungan)
}

func (h *RawatJalanHandler) GetReferensiFilter(w http.ResponseWriter, r *http.Request) {
	referensi := h.rawatJalanUsecase.GetReferensiFilter(r.Context())
	response.Success(w, "Berhasil mengambil referensi filter", referensi)
}

