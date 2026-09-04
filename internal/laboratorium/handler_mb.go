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

func (h *Handler) SimpanPermintaanLabMB(w http.ResponseWriter, r *http.Request) {
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

	var req SimpanPermintaanLabMBRequest
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

	if err := h.decryptTindakanLabMBPayload(&req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data, err := h.service.SimpanPermintaanLabMB(r.Context(), kodeDokter, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabMB(data)
	}

	response.Created(w, "Berhasil mengirim permintaan laboratorium MB", data)
}

func (h *Handler) DaftarPermintaanLabMB(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.service.GetDaftarPermintaanLabMB(r.Context(), noRawat, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptPermintaanLabMBList(data)
	response.Success(w, "Berhasil mengambil daftar permintaan laboratorium MB", data)
}

func (h *Handler) DaftarPermintaanLabMBByRM(w http.ResponseWriter, r *http.Request) {
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

	data, meta, err := h.service.GetRiwayatPermintaanLabMBByRM(r.Context(), noRM, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptPermintaanLabMBList(data)
	response.SuccessWithMeta(w, "Berhasil mengambil riwayat permintaan laboratorium MB pasien", data, meta)
}

func (h *Handler) DetailPermintaanLabMB(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.service.GetDetailPermintaanLabMB(r.Context(), noRawat, noPermintaan, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabMB(data)
	}

	response.Success(w, "Berhasil mengambil detail permintaan laboratorium MB", data)
}

func (h *Handler) UpdatePermintaanLabMB(w http.ResponseWriter, r *http.Request) {
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

	var req SimpanPermintaanLabMBRequest
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

	if err := h.decryptTindakanLabMBPayload(&req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data, err := h.service.UpdatePermintaanLabMB(r.Context(), kodeDokter, noRawatURL, noPermintaanURL, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptDetailPermintaanLabMB(data)
	}

	response.Success(w, "Berhasil memperbarui permintaan laboratorium MB", data)
}

func (h *Handler) HapusPermintaanLabMB(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.HapusPermintaanLabMB(r.Context(), noRawat, noPermintaan, statusLanjut, kodeDokter); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil membatalkan permintaan laboratorium MB", nil)
}

func (h *Handler) decryptTindakanLabMBPayload(req *SimpanPermintaanLabMBRequest) error {
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

func (h *Handler) encryptPermintaanLabMB(item *PermintaanLabMB) {
	if item == nil {
		return
	}
	item.Id, _ = crypto.Encrypt(item.NoPermintaan, h.encryptionKey)
	item.IdKunjungan, _ = crypto.Encrypt(item.NoRawat, h.encryptionKey)
}

func (h *Handler) encryptPermintaanLabMBList(items []PermintaanLabMB) {
	for i := range items {
		h.encryptPermintaanLabMB(&items[i])
	}
}

func (h *Handler) encryptDetailPermintaanLabMB(detail *DetailPermintaanLabMB) {
	if detail == nil {
		return
	}
	h.encryptPermintaanLabMB(&detail.PermintaanLabMB)
	for i := range detail.Pemeriksaan {
		detail.Pemeriksaan[i].IdTindakan, _ = crypto.Encrypt(detail.Pemeriksaan[i].KodeTindakan, h.encryptionKey)
		for j := range detail.Pemeriksaan[i].DetailTemplate {
			if detail.Pemeriksaan[i].DetailTemplate[j].IdTemplate != "" {
				detail.Pemeriksaan[i].DetailTemplate[j].IdTemplate, _ = crypto.Encrypt(detail.Pemeriksaan[i].DetailTemplate[j].IdTemplate, h.encryptionKey)
			}
		}
	}
}
