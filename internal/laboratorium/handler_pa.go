package laboratorium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

func (h *Handler) SimpanPermintaanLabPA(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)"))
		return
	}

	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	noRawatURL, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req SimpanPermintaanLabPARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	if req.NoRawat != noRawatURL {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada payload tidak cocok dengan ID kunjungan"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	if err := h.decryptTindakanLabPAPayload(&req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data, err := h.service.SimpanPermintaanLabPA(r.Context(), kodeDokter, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabPA(data)
	}

	response.Created(w, "Berhasil mengirim permintaan laboratorium PA", data)
}

func (h *Handler) DaftarPermintaanLabPA(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap, Semua)"))
		return
	}

	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	data, err := h.service.GetDaftarPermintaanLabPA(r.Context(), noRawat, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptPermintaanLabPAList(data)
	response.Success(w, "Berhasil mengambil daftar permintaan laboratorium PA", data)
}

func (h *Handler) DaftarPermintaanLabPAByRM(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap, Semua)"))
		return
	}

	idPasien := strings.TrimSpace(r.PathValue("id_pasien"))
	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	tanggal := r.URL.Query().Get("tanggal")

	filter := FilterRiwayatLab{
		Tanggal: tanggal,
		Page:    page,
		Limit:   limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.GetRiwayatPermintaanLabPAByRM(r.Context(), noRM, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptPermintaanLabPAList(data)
	response.SuccessWithMeta(w, "Berhasil mengambil riwayat permintaan laboratorium PA pasien", data, meta)
}

func (h *Handler) DetailPermintaanLabPA(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap, Semua)"))
		return
	}

	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	idPermintaan := strings.TrimSpace(r.PathValue("id_permintaan"))
	noPermintaan, err := crypto.Decrypt(idPermintaan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID permintaan laboratorium tidak valid"))
		return
	}

	data, err := h.service.GetDetailPermintaanLabPA(r.Context(), noRawat, noPermintaan, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabPA(data)
	}

	response.Success(w, "Berhasil mengambil detail permintaan laboratorium PA", data)
}

func (h *Handler) UpdatePermintaanLabPA(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)"))
		return
	}

	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	noRawatURL, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	idPermintaan := strings.TrimSpace(r.PathValue("id_permintaan"))
	noPermintaanURL, err := crypto.Decrypt(idPermintaan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID permintaan laboratorium tidak valid"))
		return
	}

	var req SimpanPermintaanLabPARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	if req.NoRawat != noRawatURL {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada payload tidak cocok dengan ID kunjungan"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	if err := h.decryptTindakanLabPAPayload(&req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data, err := h.service.UpdatePermintaanLabPA(r.Context(), kodeDokter, noRawatURL, noPermintaanURL, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabPA(data)
	}

	response.Success(w, "Berhasil memperbarui permintaan laboratorium PA", data)
}

func (h *Handler) HapusPermintaanLabPA(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap, Semua)"))
		return
	}

	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	idPermintaan := strings.TrimSpace(r.PathValue("id_permintaan"))
	noPermintaan, err := crypto.Decrypt(idPermintaan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID permintaan laboratorium tidak valid"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if err := h.service.HapusPermintaanLabPA(r.Context(), noRawat, noPermintaan, statusLanjut, kodeDokter); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil membatalkan permintaan laboratorium PA", nil)
}

func (h *Handler) decryptTindakanLabPAPayload(req *SimpanPermintaanLabPARequest) error {
	for i := range req.Pemeriksaan {
		kodeTindakan, err := crypto.Decrypt(req.Pemeriksaan[i].IdTindakan, h.encryptionKey)
		if err != nil {
			return apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan ke-%d: ID tindakan tidak valid", i+1))
		}
		req.Pemeriksaan[i].KodeTindakan = kodeTindakan
	}
	return nil
}

func (h *Handler) encryptPermintaanLabPA(item *PermintaanLabPA) {
	if item == nil {
		return
	}
	item.Id, _ = crypto.Encrypt(item.NoPermintaan, h.encryptionKey)
	item.IdKunjungan, _ = crypto.Encrypt(item.NoRawat, h.encryptionKey)
}

func (h *Handler) encryptPermintaanLabPAList(items []PermintaanLabPA) {
	for i := range items {
		h.encryptPermintaanLabPA(&items[i])
	}
}

func (h *Handler) encryptDetailPermintaanLabPA(detail *DetailPermintaanLabPA) {
	if detail == nil {
		return
	}
	h.encryptPermintaanLabPA(&detail.PermintaanLabPA)
	for i := range detail.Pemeriksaan {
		detail.Pemeriksaan[i].IdTindakan, _ = crypto.Encrypt(detail.Pemeriksaan[i].KodeTindakan, h.encryptionKey)
	}
}
