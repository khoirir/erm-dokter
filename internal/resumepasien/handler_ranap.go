package resumepasien

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

func (h *Handler) DetailResumePasienRanap(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	detail, err := h.service.DetailResumePasienRanap(r.Context(), noRawat)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	detail.IdKunjungan = idKunjungan
	response.Success(w, "Berhasil mengambil detail resume pasien rawat inap", detail)
}

func (h *Handler) RiwayatResumePasienRanapByNoRM(w http.ResponseWriter, r *http.Request) {
	idPasien := r.PathValue("id_pasien")
	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	riwayat, err := h.service.RiwayatResumePasienRanapByNoRM(r.Context(), noRM)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range riwayat {
		if enc, err := crypto.Encrypt(riwayat[i].NoRawat, h.encryptionKey); err == nil {
			riwayat[i].IdKunjungan = enc
		}
	}

	response.Success(w, "Berhasil mengambil riwayat resume pasien rawat inap", riwayat)
}

func (h *Handler) SimpanResumePasienRanap(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req SimpanResumePasienRanapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	if req.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada payload tidak cocok dengan ID kunjungan"))
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

	hasil, err := h.service.SimpanResumePasienRanap(r.Context(), noRawat, kodeDokter, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan
	response.Created(w, "Berhasil menyimpan resume pasien rawat inap", hasil)
}

func (h *Handler) UpdateResumePasienRanap(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req UpdateResumePasienRanapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
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

	hasil, err := h.service.UpdateResumePasienRanap(r.Context(), kodeDokter, noRawat, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan
	response.Success(w, "Berhasil memperbarui resume pasien rawat inap", hasil)
}

func (h *Handler) HapusResumePasienRanap(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.HapusResumePasienRanap(r.Context(), kodeDokter, noRawat); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil menghapus resume pasien rawat inap", nil)
}
