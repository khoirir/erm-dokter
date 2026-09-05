package radiologi

import (
	"net/http"
	"strconv"
	"strings"

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

	data, meta, err := h.service.GetRiwayatRadiologiKunjungan(r.Context(), noRawat, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range data {
		item := &data[i]
		if encrypted, err := crypto.Encrypt(item.CompositeKey(), h.encryptionKey); err == nil {
			item.Id = encrypted
		}
		if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, h.encryptionKey); err == nil {
			item.IdKunjungan = encryptedIdKunjungan
		}
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat hasil radiologi", data, meta)
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

	data, meta, err := h.service.GetRiwayatRadiologiPasien(r.Context(), noRM, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range data {
		item := &data[i]
		if encrypted, err := crypto.Encrypt(item.CompositeKey(), h.encryptionKey); err == nil {
			item.Id = encrypted
		}
		if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, h.encryptionKey); err == nil {
			item.IdKunjungan = encryptedIdKunjungan
		}
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat hasil radiologi pasien", data, meta)
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
	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Semua, Ralan, Ranap)"))
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
