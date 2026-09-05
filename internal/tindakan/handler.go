package tindakan

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
	mux.HandleFunc("GET /api/v1/tindakan/lab/{kategori}", authMiddleware(timeoutMiddleware(h.GetDaftarTindakanLab)))
	mux.HandleFunc("GET /api/v1/tindakan/lab/{kategori}/{id_tindakan}", authMiddleware(timeoutMiddleware(h.GetDetailTindakanLab)))
	mux.HandleFunc("GET /api/v1/tindakan/radiologi", authMiddleware(timeoutMiddleware(h.GetDaftarTindakanRadiologi)))
	mux.HandleFunc("GET /api/v1/tindakan/radiologi/{id_tindakan}", authMiddleware(timeoutMiddleware(h.GetDetailTindakanRadiologi)))
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

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.GetDaftarTindakanLab(r.Context(), kat, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range data {
		if encId, errEnc := crypto.Encrypt(data[i].KodeTindakan, h.encryptionKey); errEnc == nil {
			data[i].Id = encId
		}
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

	kodeTindakan, err := crypto.Decrypt(idTindakan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID tindakan lab tidak valid"))
		return
	}

	data, err := h.service.GetDetailTindakanLab(r.Context(), kat, kodeTindakan)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data.Id = idTindakan
	for i := range data.Templates {
		if encTemplateId, errEnc := crypto.Encrypt(data.Templates[i].IdTemplate, h.encryptionKey); errEnc == nil {
			data.Templates[i].IdTemplate = encTemplateId
		}
	}

	response.Success(w, "Berhasil mengambil detail tindakan laboratorium", data)
}

func (h *Handler) GetDaftarTindakanRadiologi(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	keyword := r.URL.Query().Get("keyword")

	filter := FilterDaftarTindakanRadiologi{
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.GetDaftarTindakanRadiologi(r.Context(), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range data {
		if encId, errEnc := crypto.Encrypt(data[i].KodeTindakan, h.encryptionKey); errEnc == nil {
			data[i].Id = encId
		}
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar tindakan radiologi", data, meta)
}

func (h *Handler) GetDetailTindakanRadiologi(w http.ResponseWriter, r *http.Request) {
	idTindakan := strings.TrimSpace(r.PathValue("id_tindakan"))
	if idTindakan == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID tindakan wajib diisi"))
		return
	}

	kodeTindakan, err := crypto.Decrypt(idTindakan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID tindakan radiologi tidak valid"))
		return
	}

	data, err := h.service.GetDetailTindakanRadiologi(r.Context(), kodeTindakan)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	data.Id = idTindakan

	response.Success(w, "Berhasil mengambil detail tindakan radiologi", data)
}

