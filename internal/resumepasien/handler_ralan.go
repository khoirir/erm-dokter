package resumepasien

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

func (h *Handler) ReferensiRalan(w http.ResponseWriter, r *http.Request) {
	ref := h.service.ReferensiRalan(r.Context())
	response.Success(w, "Berhasil mengambil opsi referensi resume ralan", ref)
}

func (h *Handler) DetailResumePasienRalan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	detail, err := h.service.DetailResumePasienRalan(r.Context(), noRawat)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	detail.IdKunjungan = idKunjungan
	response.Success(w, "Berhasil mengambil detail resume pasien rawat jalan", detail)
}

func (h *Handler) RiwayatResumePasienRalanByNoRM(w http.ResponseWriter, r *http.Request) {
	idPasien := r.PathValue("id_pasien")
	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	riwayat, err := h.service.RiwayatResumePasienRalanByNoRM(r.Context(), noRM)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range riwayat {
		if enc, err := crypto.Encrypt(riwayat[i].NoRawat, h.encryptionKey); err == nil {
			riwayat[i].IdKunjungan = enc
		}
	}

	response.Success(w, "Berhasil mengambil riwayat resume pasien rawat jalan", riwayat)
}

func (h *Handler) SimpanResumePasienRalan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req SimpanResumePasienRalanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format data resume tidak valid"))
		return
	}

	if req.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada resume tidak sesuai dengan ID kunjungan"))
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

	hasil, err := h.service.SimpanResumePasienRalan(r.Context(), noRawat, kodeDokter, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan
	response.Created(w, "Berhasil menyimpan resume pasien rawat jalan", hasil)
}

func (h *Handler) UpdateResumePasienRalan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req UpdateResumePasienRalanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format data resume tidak valid"))
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

	hasil, err := h.service.UpdateResumePasienRalan(r.Context(), kodeDokter, noRawat, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan
	response.Success(w, "Berhasil memperbarui resume pasien rawat jalan", hasil)
}

func (h *Handler) HapusResumePasienRalan(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.HapusResumePasienRalan(r.Context(), kodeDokter, noRawat); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil menghapus resume pasien rawat jalan", nil)
}
