package penilaianmedis

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

func (h *Handler) DetailPenilaianMedisRalanKandungan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	detail, err := h.service.DetailPenilaianMedisRalanKandungan(r.Context(), noRawat)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	detail.IdKunjungan = idKunjungan

	response.Success(w, "Berhasil mengambil detail penilaian awal medis rawat jalan kandungan", detail)
}

func (h *Handler) RiwayatPenilaianMedisRalanKandungan(w http.ResponseWriter, r *http.Request) {
	idPasien := r.PathValue("id_pasien")
	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	riwayat, err := h.service.RiwayatPenilaianMedisRalanKandunganByNoRM(r.Context(), noRM)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range riwayat {
		if enc, err := crypto.Encrypt(riwayat[i].NoRawat, h.encryptionKey); err == nil {
			riwayat[i].IdKunjungan = enc
		}
	}

	response.Success(w, "Berhasil mengambil riwayat penilaian awal medis rawat jalan kandungan pasien", riwayat)
}

func (h *Handler) SimpanPenilaianMedisRalanKandungan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req SimpanPenilaianMedisRalanKandunganRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format data penilaian medis tidak valid"))
		return
	}

	if req.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada penilaian medis tidak sesuai dengan ID kunjungan"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil, err := h.service.SimpanPenilaianMedisRalanKandungan(r.Context(), kodeDokter, noRawat, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan

	response.Created(w, "Berhasil menyimpan penilaian awal medis rawat jalan kandungan", hasil)
}

func (h *Handler) UpdatePenilaianMedisRalanKandungan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req UpdatePenilaianMedisRalanKandunganRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format data penilaian medis tidak valid"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil, err := h.service.UpdatePenilaianMedisRalanKandungan(r.Context(), kodeDokter, noRawat, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan

	response.Success(w, "Berhasil memperbarui penilaian awal medis rawat jalan kandungan", hasil)
}

func (h *Handler) HapusPenilaianMedisRalanKandungan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if err := h.service.HapusPenilaianMedisRalanKandungan(r.Context(), kodeDokter, noRawat); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil menghapus penilaian awal medis rawat jalan kandungan", nil)
}
