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

func (h *Handler) SimpanPermintaanLabPK(w http.ResponseWriter, r *http.Request) {
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

	var req SimpanPermintaanLabPKRequest
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

	if err := h.decryptTindakanLabPayload(&req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data, err := h.service.SimpanPermintaanLabPK(r.Context(), kodeDokter, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabPK(data)
	}

	response.Created(w, "Berhasil mengirim permintaan laboratorium PK", data)
}

func (h *Handler) DaftarPermintaanLabPK(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.service.GetDaftarPermintaanLabPK(r.Context(), noRawat, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptPermintaanLabPKList(data)
	response.Success(w, "Berhasil mengambil daftar permintaan laboratorium PK", data)
}

func (h *Handler) DaftarPermintaanLabPKByRM(w http.ResponseWriter, r *http.Request) {
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

	data, meta, err := h.service.GetRiwayatPermintaanLabPKByRM(r.Context(), noRM, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptPermintaanLabPKList(data)
	response.SuccessWithMeta(w, "Berhasil mengambil riwayat permintaan laboratorium PK pasien", data, meta)
}

func (h *Handler) DetailPermintaanLabPK(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)"))
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

	data, err := h.service.GetDetailPermintaanLabPK(r.Context(), noRawat, noPermintaan, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabPK(data)
	}

	response.Success(w, "Berhasil mengambil detail permintaan laboratorium PK", data)
}

func (h *Handler) UpdatePermintaanLabPK(w http.ResponseWriter, r *http.Request) {
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

	var req SimpanPermintaanLabPKRequest
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

	if err := h.decryptTindakanLabPayload(&req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data, err := h.service.UpdatePermintaanLabPK(r.Context(), kodeDokter, noRawatURL, noPermintaanURL, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabPK(data)
	}

	response.Success(w, "Berhasil memperbarui permintaan laboratorium PK", data)
}

func (h *Handler) HapusPermintaanLabPK(w http.ResponseWriter, r *http.Request) {
	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)"))
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

	if err := h.service.HapusPermintaanLabPK(r.Context(), noRawat, noPermintaan, statusLanjut, kodeDokter); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil membatalkan permintaan laboratorium PK", nil)
}

func (h *Handler) decryptTindakanLabPayload(req *SimpanPermintaanLabPKRequest) error {
	for i := range req.Pemeriksaan {
		kodeTindakan, err := crypto.Decrypt(req.Pemeriksaan[i].IdTindakan, h.encryptionKey)
		if err != nil {
			return apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan ke-%d: ID tindakan tidak valid", i+1))
		}
		req.Pemeriksaan[i].KodeTindakan = kodeTindakan

		for j, encIdTemplate := range req.Pemeriksaan[i].IdTemplate {
			idTemplateStr, err := crypto.Decrypt(encIdTemplate, h.encryptionKey)
			if err != nil {
				return apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan ke-%d parameter ke-%d: ID template pengujian tidak valid", i+1, j+1))
			}
			idTemplateInt, err := strconv.Atoi(idTemplateStr)
			if err != nil {
				return apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan ke-%d parameter ke-%d: Format ID template tidak valid", i+1, j+1))
			}
			req.Pemeriksaan[i].KodeTemplate = append(req.Pemeriksaan[i].KodeTemplate, idTemplateInt)
		}
	}
	return nil
}

func (h *Handler) encryptPermintaanLabPK(item *PermintaanLabPK) {
	if item == nil {
		return
	}
	item.Id, _ = crypto.Encrypt(item.NoPermintaan, h.encryptionKey)
	item.IdKunjungan, _ = crypto.Encrypt(item.NoRawat, h.encryptionKey)
}

func (h *Handler) encryptPermintaanLabPKList(items []PermintaanLabPK) {
	for i := range items {
		h.encryptPermintaanLabPK(&items[i])
	}
}

func (h *Handler) encryptDetailPermintaanLabPK(detail *DetailPermintaanLabPK) {
	if detail == nil {
		return
	}
	h.encryptPermintaanLabPK(&detail.PermintaanLabPK)
	for i := range detail.Pemeriksaan {
		detail.Pemeriksaan[i].IdTindakan, _ = crypto.Encrypt(detail.Pemeriksaan[i].KodeTindakan, h.encryptionKey)
		for j := range detail.Pemeriksaan[i].DetailTemplate {
			if detail.Pemeriksaan[i].DetailTemplate[j].IdTemplate != "" {
				detail.Pemeriksaan[i].DetailTemplate[j].IdTemplate, _ = crypto.Encrypt(detail.Pemeriksaan[i].DetailTemplate[j].IdTemplate, h.encryptionKey)
			}
		}
	}
}
