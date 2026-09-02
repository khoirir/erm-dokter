package laboratorium

import (
	"net/http"
	"strconv"
	"strings"

	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/laboratorium/{kategori}/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarHasilLab)))
	mux.HandleFunc("GET /api/v1/laboratorium/{kategori}/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarHasilLabByRM)))
	mux.HandleFunc("GET /api/v1/laboratorium/{kategori}/{id_kunjungan}/{status_lanjut}/{id_hasil}", authMiddleware(timeoutMiddleware(h.DetailHasilLab)))
}

func (h *Handler) DaftarHasilLab(w http.ResponseWriter, r *http.Request) {
	kategoriRaw := r.PathValue("kategori")
	kat := shared.KategoriLab(strings.ToUpper(strings.TrimSpace(kategoriRaw)))
	if !kat.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Kategori laboratorium tidak valid (pilihan: PK, PA, MB)"))
		return
	}

	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	if idKunjungan == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan wajib diisi"))
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

	filter := FilterRiwayatLab{
		Tanggal: tanggal,
		Page:    page,
		Limit:   limit,
	}

	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.GetRiwayatLabKunjungan(r.Context(), kat, idKunjungan, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat hasil laboratorium", data, meta)
}

func (h *Handler) DaftarHasilLabByRM(w http.ResponseWriter, r *http.Request) {
	kategoriRaw := r.PathValue("kategori")
	kat := shared.KategoriLab(strings.ToUpper(strings.TrimSpace(kategoriRaw)))
	if !kat.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Kategori laboratorium tidak valid (pilihan: PK, PA, MB)"))
		return
	}

	idPasien := strings.TrimSpace(r.PathValue("id_pasien"))
	if idPasien == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien wajib diisi"))
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

	filter := FilterRiwayatLab{
		Tanggal: tanggal,
		Page:    page,
		Limit:   limit,
	}

	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.GetRiwayatLabPasien(r.Context(), kat, idPasien, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat hasil laboratorium pasien", data, meta)
}

func (h *Handler) DetailHasilLab(w http.ResponseWriter, r *http.Request) {
	kategoriRaw := r.PathValue("kategori")
	kat := shared.KategoriLab(strings.ToUpper(strings.TrimSpace(kategoriRaw)))
	if !kat.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Kategori laboratorium tidak valid (pilihan: PK, PA, MB)"))
		return
	}

	idKunjungan := strings.TrimSpace(r.PathValue("id_kunjungan"))
	if idKunjungan == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan wajib diisi"))
		return
	}

	idHasil := strings.TrimSpace(r.PathValue("id_hasil"))
	if idHasil == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID hasil lab wajib diisi"))
		return
	}

	data, err := h.service.GetDetailHasilLab(r.Context(), kat, idKunjungan, idHasil)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil detail hasil laboratorium", data)
}
