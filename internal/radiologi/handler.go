package radiologi

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

type Handler struct {
	service       Service
	encryptionKey string
}

func NewHandler(service Service, encryptionKey string) *Handler {
	return &Handler{
		service:       service,
		encryptionKey: encryptionKey,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/radiologi/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarHasilRadiologi)))
	mux.HandleFunc("GET /api/v1/radiologi/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarHasilRadiologiByRM)))
	mux.HandleFunc("GET /api/v1/radiologi/{id_kunjungan}/{status_lanjut}/{id_radiologi}", authMiddleware(timeoutMiddleware(h.DetailHasilRadiologi)))

	mux.HandleFunc("POST /api/v1/radiologi/permintaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.SimpanPermintaanRadiologi)))
	mux.HandleFunc("GET /api/v1/radiologi/permintaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPermintaanRadiologi)))
	mux.HandleFunc("GET /api/v1/radiologi/pasien/{id_pasien}/permintaan/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPermintaanRadiologiByRM)))
	mux.HandleFunc("GET /api/v1/radiologi/permintaan/{id_kunjungan}/{status_lanjut}/{id_permintaan}", authMiddleware(timeoutMiddleware(h.DetailPermintaanRadiologi)))
	mux.HandleFunc("PUT /api/v1/radiologi/permintaan/{id_kunjungan}/{status_lanjut}/{id_permintaan}", authMiddleware(timeoutMiddleware(h.UpdatePermintaanRadiologi)))
	mux.HandleFunc("DELETE /api/v1/radiologi/permintaan/{id_kunjungan}/{status_lanjut}/{id_permintaan}", authMiddleware(timeoutMiddleware(h.HapusPermintaanRadiologi)))
}


func (h *Handler) DaftarHasilRadiologi(w http.ResponseWriter, r *http.Request) {
	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	if idKunjungan == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan wajib diisi"))
		return
	}

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Semua, Ralan, Ranap)"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	tanggal := r.URL.Query().Get("tanggal")

	filter := FilterRiwayatRadiologi{
		Tanggal: tanggal,
		Page:    page,
		Limit:   limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	list, meta, err := h.service.GetRiwayatRadiologiKunjungan(r.Context(), noRawat, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range list {
		encId, err := crypto.Encrypt(list[i].CompositeKey(), h.encryptionKey)
		if err == nil {
			list[i].Id = encId
		}
		list[i].IdKunjungan = idKunjungan
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat hasil radiologi kunjungan", list, meta)
}

func (h *Handler) DaftarHasilRadiologiByRM(w http.ResponseWriter, r *http.Request) {
	idPasien := strings.TrimSpace(r.PathValue("id_pasien"))
	if idPasien == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien wajib diisi"))
		return
	}

	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Semua, Ralan, Ranap)"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	tanggal := r.URL.Query().Get("tanggal")

	filter := FilterRiwayatRadiologi{
		Tanggal: tanggal,
		Page:    page,
		Limit:   limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	list, meta, err := h.service.GetRiwayatRadiologiPasien(r.Context(), noRM, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range list {
		encId, err := crypto.Encrypt(list[i].CompositeKey(), h.encryptionKey)
		if err == nil {
			list[i].Id = encId
		}
		encNoRawat, err := crypto.Encrypt(list[i].NoRawat, h.encryptionKey)
		if err == nil {
			list[i].IdKunjungan = encNoRawat
		}
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat hasil radiologi pasien", list, meta)
}

func (h *Handler) DetailHasilRadiologi(w http.ResponseWriter, r *http.Request) {
	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	if idKunjungan == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan wajib diisi"))
		return
	}

	kunjunganNoRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)"))
		return
	}

	idRadiologi := strings.TrimSpace(r.PathValue("id_radiologi"))
	if idRadiologi == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID hasil radiologi wajib diisi"))
		return
	}

	plainId, err := crypto.Decrypt(idRadiologi, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID hasil radiologi tidak valid"))
		return
	}

	idHasil, err := ParseIdHasilRadiologi(plainId)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format ID hasil radiologi tidak valid"))
		return
	}

	if idHasil.NoRawat != kunjunganNoRawat {
		apperror.HandleError(w, apperror.NewBusinessError("ID hasil radiologi tidak cocok dengan ID kunjungan"))
		return
	}

	data, err := h.service.GetDetailHasilRadiologi(r.Context(), idHasil, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data.Id = idRadiologi
	data.IdKunjungan = idKunjungan

	response.Success(w, "Berhasil mengambil detail hasil radiologi", data)
}

func (h *Handler) encryptDetailPermintaanRadiologi(detail *DetailPermintaanRadiologi) {
	if detail == nil {
		return
	}
	if encOrder, err := crypto.Encrypt(detail.NoPermintaan, h.encryptionKey); err == nil {
		detail.NoPermintaan = encOrder
	}
	for i := range detail.Pemeriksaan {
		if encTindakan, err := crypto.Encrypt(detail.Pemeriksaan[i].KodeTindakan, h.encryptionKey); err == nil {
			detail.Pemeriksaan[i].IdTindakan = encTindakan
		}
	}
}

func (h *Handler) SimpanPermintaanRadiologi(w http.ResponseWriter, r *http.Request) {
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

	var req SimpanPermintaanRadiologiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	if req.NoRawat != noRawatURL {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan pada payload tidak cocok dengan ID kunjungan pada URL"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil || kodeDokter == "" {
		apperror.HandleError(w, apperror.NewUnauthorizedError("Sesi dokter tidak valid"))
		return
	}

	for i := range req.Pemeriksaan {
		plainTindakan, err := crypto.Decrypt(req.Pemeriksaan[i].IdTindakan, h.encryptionKey)
		if err != nil {
			apperror.HandleError(w, apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan ke-%d: ID tindakan tidak valid", i+1)))
			return
		}
		req.Pemeriksaan[i].KodeTindakan = plainTindakan
	}

	data, err := h.service.SimpanPermintaanRadiologi(r.Context(), kodeDokter, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptDetailPermintaanRadiologi(data)
	response.Created(w, "Berhasil menyimpan permintaan radiologi", data)
}

func (h *Handler) DaftarPermintaanRadiologi(w http.ResponseWriter, r *http.Request) {
	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	if idKunjungan == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan wajib diisi"))
		return
	}

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Semua, Ralan, Ranap)"))
		return
	}

	list, err := h.service.GetDaftarPermintaanRadiologi(r.Context(), noRawat, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range list {
		h.encryptDetailPermintaanRadiologi(&list[i])
	}

	response.Success(w, "Berhasil mengambil daftar permintaan radiologi", list)
}

func (h *Handler) DaftarPermintaanRadiologiByRM(w http.ResponseWriter, r *http.Request) {
	idPasien := strings.TrimSpace(r.PathValue("id_pasien"))
	if idPasien == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien wajib diisi"))
		return
	}

	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	statusLanjutRaw := strings.TrimSpace(r.PathValue("status_lanjut"))
	statusLanjut := shared.StatusLanjut(statusLanjutRaw)
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Semua, Ralan, Ranap)"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	tanggal := r.URL.Query().Get("tanggal")

	filter := FilterRiwayatPermintaanRadiologi{
		Tanggal: tanggal,
		Page:    page,
		Limit:   limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	list, meta, err := h.service.GetRiwayatPermintaanRadiologiByRM(r.Context(), noRM, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range list {
		h.encryptDetailPermintaanRadiologi(&list[i])
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat permintaan radiologi pasien", list, meta)
}

func (h *Handler) DetailPermintaanRadiologi(w http.ResponseWriter, r *http.Request) {
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
		apperror.HandleError(w, apperror.NewBusinessError("ID permintaan radiologi tidak valid"))
		return
	}

	data, err := h.service.GetDetailPermintaanRadiologi(r.Context(), noRawatURL, noPermintaanURL, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptDetailPermintaanRadiologi(data)
	response.Success(w, "Berhasil mengambil detail permintaan radiologi", data)
}

func (h *Handler) UpdatePermintaanRadiologi(w http.ResponseWriter, r *http.Request) {
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
		apperror.HandleError(w, apperror.NewBusinessError("ID permintaan radiologi tidak valid"))
		return
	}

	var req SimpanPermintaanRadiologiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	if req.NoRawat != noRawatURL {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan pada payload tidak cocok dengan ID kunjungan pada URL"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil || kodeDokter == "" {
		apperror.HandleError(w, apperror.NewUnauthorizedError("Sesi dokter tidak valid"))
		return
	}

	for i := range req.Pemeriksaan {
		plainTindakan, err := crypto.Decrypt(req.Pemeriksaan[i].IdTindakan, h.encryptionKey)
		if err != nil {
			apperror.HandleError(w, apperror.NewBusinessError(fmt.Sprintf("Pemeriksaan ke-%d: ID tindakan tidak valid", i+1)))
			return
		}
		req.Pemeriksaan[i].KodeTindakan = plainTindakan
	}

	data, err := h.service.UpdatePermintaanRadiologi(r.Context(), kodeDokter, noRawatURL, noPermintaanURL, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptDetailPermintaanRadiologi(data)
	response.Success(w, "Berhasil memperbarui permintaan radiologi", data)
}

func (h *Handler) HapusPermintaanRadiologi(w http.ResponseWriter, r *http.Request) {
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
		apperror.HandleError(w, apperror.NewBusinessError("ID permintaan radiologi tidak valid"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil || kodeDokter == "" {
		apperror.HandleError(w, apperror.NewUnauthorizedError("Sesi dokter tidak valid"))
		return
	}

	err = h.service.HapusPermintaanRadiologi(r.Context(), noRawatURL, noPermintaanURL, statusLanjut, kodeDokter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil membatalkan permintaan radiologi", nil)
}
