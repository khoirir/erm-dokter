package rawatinap

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	rawatInapService Service
	encryptionKey    string
}

func NewHandler(service Service, encryptionKey string) *Handler {
	return &Handler{
		rawatInapService: service,
		encryptionKey:    encryptionKey,
	}
}

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	authMiddleware func(http.HandlerFunc) http.HandlerFunc,
	serviceAuthMiddleware func(http.HandlerFunc) http.HandlerFunc,
	timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc,
) {
	mux.HandleFunc("GET /api/v1/rawat-inap/status-pulang", authMiddleware(timeoutMiddleware(h.DaftarStatusPulang)))
	mux.HandleFunc("GET /api/v1/rawat-inap/pasien", serviceAuthMiddleware(timeoutMiddleware(h.DaftarPasienRawatInap)))
	mux.HandleFunc("GET /api/v1/rawat-inap/{id}", authMiddleware(timeoutMiddleware(h.DetailPasienRawatInap)))
}

func (h *Handler) DaftarStatusPulang(w http.ResponseWriter, r *http.Request) {
	data := h.rawatInapService.DaftarStatusPulang(r.Context())
	response.Success(w, "Berhasil mengambil referensi status pulang", data)
}

func (h *Handler) DaftarPasienRawatInap(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := FilterPasienRawatInap{
		Bangsal:         q.Get("bangsal"),
		Kelas:           q.Get("kelas"),
		Penjamin:        q.Get("penjamin"),
		ScopeDPJP:       ScopeDPJP(q.Get("scope_dpjp")),
		StatusKunjungan: StatusKunjunganRanap(q.Get("status_kunjungan")),
		StatusPulang:    StatusPulang(q.Get("status_pulang")),
		Tanggal:         q.Get("tanggal"),
		Keyword:         q.Get("keyword"),
		OrderBy:         OrderByRanap(q.Get("order_by")),
		SortOrder:       q.Get("sort_order"),
		Page:            page,
		Limit:           limit,
	}

	filter.Sanitize()
	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	kodeDokter, isService, err := middleware.GetKodeDokterOrEmpty(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if isService && filter.ScopeDPJP == "" {
		filter.ScopeDPJP = ScopeDPJPSemua
	}

	daftar, meta, err := h.rawatInapService.DaftarPasienRawatInap(r.Context(), kodeDokter, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range daftar {
		item := &daftar[i]
		idComposite := fmt.Sprintf("%s~%s~%s", item.NoRawat, item.TanggalMasuk, item.JamMasuk)
		if encryptedId, err := crypto.Encrypt(idComposite, h.encryptionKey); err == nil {
			item.Id = encryptedId
		}
		if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, h.encryptionKey); err == nil {
			item.IdKunjungan = encryptedIdKunjungan
		}
		if encryptedIdPasien, err := crypto.Encrypt(item.NoRekamMedis, h.encryptionKey); err == nil {
			item.IdPasien = encryptedIdPasien
		}
		item.NoRekamMedis = item.FormatNoRekamMedis()
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar pasien rawat inap", daftar, meta)
}

func (h *Handler) DetailPasienRawatInap(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID rawat inap wajib diisi"))
		return
	}

	decrypted, err := crypto.Decrypt(id, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID rawat inap tidak valid"))
		return
	}

	parts := strings.Split(decrypted, "~")
	if len(parts) != 3 {
		apperror.HandleError(w, apperror.NewBusinessError("ID rawat inap tidak valid"))
		return
	}

	noRawat := strings.TrimSpace(parts[0])
	tglMasuk := strings.TrimSpace(parts[1])
	jamMasuk := strings.TrimSpace(parts[2])

	if noRawat == "" || tglMasuk == "" || jamMasuk == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID rawat inap tidak valid"))
		return
	}

	item, err := h.rawatInapService.DetailPasienRawatInap(r.Context(), noRawat, tglMasuk, jamMasuk)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if item == nil {
		apperror.HandleError(w, apperror.NewNotFoundError("Detail pasien rawat inap tidak ditemukan"))
		return
	}

	idComposite := fmt.Sprintf("%s~%s~%s", item.NoRawat, item.TanggalMasuk, item.JamMasuk)
	if encryptedId, err := crypto.Encrypt(idComposite, h.encryptionKey); err == nil {
		item.Id = encryptedId
	}
	if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, h.encryptionKey); err == nil {
		item.IdKunjungan = encryptedIdKunjungan
	}
	if encryptedIdPasien, err := crypto.Encrypt(item.NoRekamMedis, h.encryptionKey); err == nil {
		item.IdPasien = encryptedIdPasien
	}
	item.NoRekamMedis = item.FormatNoRekamMedis()

	response.Success(w, "Berhasil mengambil detail pasien rawat inap", item)
}
