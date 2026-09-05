package obat

import (
	"net/http"
	"strconv"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/shared/formatter"
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
	mux.HandleFunc("GET /api/v1/obat", authMiddleware(timeoutMiddleware(h.DaftarObat)))
	mux.HandleFunc("GET /api/v1/obat/{id_obat}", authMiddleware(timeoutMiddleware(h.DetailObat)))
	mux.HandleFunc("GET /api/v1/obat/jenis", authMiddleware(timeoutMiddleware(h.DaftarJenis)))
	mux.HandleFunc("GET /api/v1/obat/golongan", authMiddleware(timeoutMiddleware(h.DaftarGolongan)))
	mux.HandleFunc("GET /api/v1/obat/kategori", authMiddleware(timeoutMiddleware(h.DaftarKategori)))
}

func (h *Handler) DaftarObat(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := FilterDaftarObat{
		Depo:      q.Get("depo"),
		Jenis:     q.Get("jenis"),
		Golongan:  q.Get("golongan"),
		Kategori:  q.Get("kategori"),
		Keyword:   q.Get("keyword"),
		OrderBy:   q.Get("order_by"),
		SortOrder: q.Get("sort_order"),
		Page:      page,
		Limit:     limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	data, meta, err := h.service.DaftarObat(r.Context(), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range data {
		if encryptedId, err := crypto.Encrypt(data[i].KodeObat, h.encryptionKey); err == nil {
			data[i].Id = encryptedId
		}
		data[i].Harga = formatter.FormatHarga(data[i].Harga)
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar obat", data, meta)
}

func (h *Handler) DetailObat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id_obat")
	if id == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID obat wajib diisi"))
		return
	}

	kodeObat, err := crypto.Decrypt(id, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID obat tidak valid"))
		return
	}

	detail, err := h.service.DetailObat(r.Context(), kodeObat)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	detail.Id = id
	detail.Harga = formatter.FormatHarga(detail.Harga)

	response.Success(w, "Berhasil mengambil detail obat", detail)
}

func (h *Handler) DaftarJenis(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarJenis(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar jenis obat", data)
}

func (h *Handler) DaftarGolongan(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarGolongan(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar golongan obat", data)
}

func (h *Handler) DaftarKategori(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarKategori(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	response.Success(w, "Berhasil mengambil daftar kategori obat", data)
}
