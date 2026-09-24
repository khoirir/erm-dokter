package master

import (
	"net/http"
	"strconv"

	"erm-dokter/internal/pkg/response"
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
	mux.HandleFunc("GET /api/v1/master/penjamin", authMiddleware(timeoutMiddleware(h.DaftarPenjamin)))
	mux.HandleFunc("GET /api/v1/master/depo", authMiddleware(timeoutMiddleware(h.DaftarDepo)))
	mux.HandleFunc("GET /api/v1/master/poliklinik", authMiddleware(timeoutMiddleware(h.DaftarPoliklinik)))
	mux.HandleFunc("GET /api/v1/master/bangsal", authMiddleware(timeoutMiddleware(h.DaftarBangsal)))
	mux.HandleFunc("GET /api/v1/master/kelas", authMiddleware(timeoutMiddleware(h.DaftarKelas)))
	mux.HandleFunc("GET /api/v1/master/icd10", authMiddleware(timeoutMiddleware(h.DaftarICD10)))
	mux.HandleFunc("GET /api/v1/master/icd9", authMiddleware(timeoutMiddleware(h.DaftarICD9)))
	mux.HandleFunc("POST /api/v1/master/sync-icd", authMiddleware(timeoutMiddleware(h.SyncICD)))
}

func (h *Handler) DaftarPenjamin(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarPenjamin(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar penjamin", data)
}

func (h *Handler) DaftarDepo(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarDepo(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar depo", data)
}

func (h *Handler) DaftarPoliklinik(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarPoliklinik(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar poliklinik", data)
}

func (h *Handler) DaftarBangsal(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarBangsal(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar bangsal", data)
}

func (h *Handler) DaftarKelas(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarKelas(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar kelas kamar", data)
}

func (h *Handler) DaftarICD10(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	keyword := r.URL.Query().Get("keyword")

	filter := FilterMasterICD{
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.DaftarICD10(r.Context(), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar master ICD-10", data, meta)
}

func (h *Handler) DaftarICD9(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	keyword := r.URL.Query().Get("keyword")

	filter := FilterMasterICD{
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.DaftarICD9(r.Context(), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar master ICD-9", data, meta)
}

func (h *Handler) SyncICD(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.SyncICD(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil menyinkronkan data master ICD-10 dan ICD-9 ke memori", result)
}
