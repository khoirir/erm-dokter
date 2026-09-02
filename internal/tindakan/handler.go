package tindakan

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
	mux.HandleFunc("GET /api/v1/tindakan/lab/{kategori}", authMiddleware(timeoutMiddleware(h.GetDaftarTindakanLab)))
	mux.HandleFunc("GET /api/v1/tindakan/lab/{kategori}/{id_tindakan}", authMiddleware(timeoutMiddleware(h.GetDetailTindakanLab)))
}

func (h *Handler) GetDaftarTindakanLab(w http.ResponseWriter, r *http.Request) {
	kategoriRaw := r.PathValue("kategori")
	kat := shared.KategoriLab(strings.ToUpper(strings.TrimSpace(kategoriRaw)))
	if !kat.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Kategori laboratorium tidak valid (pilihan: PK, PA, MB)"))
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	keyword := r.URL.Query().Get("keyword")

	filter := FilterDaftarTindakanLab{
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
	}

	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.GetDaftarTindakanLab(r.Context(), kat, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar tindakan laboratorium", data, meta)
}

func (h *Handler) GetDetailTindakanLab(w http.ResponseWriter, r *http.Request) {
	kategoriRaw := r.PathValue("kategori")
	kat := shared.KategoriLab(strings.ToUpper(strings.TrimSpace(kategoriRaw)))
	if !kat.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("Kategori laboratorium tidak valid (pilihan: PK, PA, MB)"))
		return
	}

	idTindakan := strings.TrimSpace(r.PathValue("id_tindakan"))
	if idTindakan == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID tindakan wajib diisi"))
		return
	}

	data, err := h.service.GetDetailTindakanLab(r.Context(), kat, idTindakan)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil detail tindakan laboratorium", data)
}
