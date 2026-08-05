package handler

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/domain"
	"erm-dokter/internal/middleware"
	"erm-dokter/pkg/response"
)

type RawatJalanHandler struct {
	rawatJalanUsecase domain.RawatJalanUsecase
}

func NewRawatJalanHandler(usecase domain.RawatJalanUsecase) *RawatJalanHandler {
	return &RawatJalanHandler{
		rawatJalanUsecase: usecase,
	}
}

func (h *RawatJalanHandler) DaftarAntreanDokter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan", nil)
		return
	}

	var filter domain.FilterAntreanDokter
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
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar antrean dokter", daftarAntrean, meta)
}

func (h *RawatJalanHandler) DetailKunjungan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan", nil)
		return
	}

	noRawat := r.URL.Query().Get("no_rawat")
	if noRawat == "" {
		response.Error(w, http.StatusBadRequest, "Parameter no_rawat wajib diisi", nil)
		return
	}

	kunjungan, err := h.rawatJalanUsecase.DetailKunjungan(r.Context(), noRawat)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if kunjungan == nil {
		response.Error(w, http.StatusNotFound, "Detail kunjungan pasien tidak ditemukan", nil)
		return
	}

	response.Success(w, "Berhasil mengambil detail kunjungan", kunjungan)
}
