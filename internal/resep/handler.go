package resep

import (
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /api/v1/resep/aturan-pakai", authMiddleware(timeoutMiddleware(h.DaftarAturanPakai)))
	mux.HandleFunc("GET /api/v1/resep/metode-racik", authMiddleware(timeoutMiddleware(h.DaftarMetodeRacik)))
	mux.HandleFunc("GET /api/v1/resep/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarResep)))
	mux.HandleFunc("GET /api/v1/resep/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarResepByRM)))
}

func (h *Handler) DaftarAturanPakai(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	data, err := h.service.DaftarAturanPakai(r.Context(), keyword)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil daftar aturan pakai", data)
}

func (h *Handler) DaftarMetodeRacik(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.DaftarMetodeRacik(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil daftar metode racik", data)
}

func (h *Handler) DaftarResep(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")

	status := shared.StatusLanjut(statusLanjut)
	if status != "Semua" && !status.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("status lanjut tidak valid"))
		return
	}

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := FilterDaftarResep{
		Tanggal: q.Get("tanggal"),
		Page:    page,
		Limit:   limit,
	}

	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	daftarResep, meta, err := h.service.DaftarResep(r.Context(), noRawat, status, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptResepResponse(daftarResep, idKunjungan)

	response.SuccessWithMeta(w, "Berhasil mengambil daftar resep", daftarResep, meta)
}

func (h *Handler) DaftarResepByRM(w http.ResponseWriter, r *http.Request) {
	idPasien := r.PathValue("id_pasien")
	statusLanjut := r.PathValue("status_lanjut")

	status := shared.StatusLanjut(statusLanjut)
	if status != "Semua" && !status.IsValid() {
		apperror.HandleError(w, apperror.NewBusinessError("status lanjut tidak valid"))
		return
	}

	noRekamMedis, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := FilterDaftarResep{
		Tanggal: q.Get("tanggal"),
		Page:    page,
		Limit:   limit,
	}

	if errs := filter.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	daftarResep, meta, err := h.service.DaftarResepByRM(r.Context(), noRekamMedis, status, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptResepResponse(daftarResep, "")

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat resep pasien", daftarResep, meta)
}

func (h *Handler) encryptResepResponse(daftarResep []Resep, defaultIdKunjungan string) {
	for i := range daftarResep {
		res := &daftarResep[i]

		if encryptedId, err := crypto.Encrypt(res.NoResep, h.encryptionKey); err == nil {
			res.Id = encryptedId
		}

		if defaultIdKunjungan != "" {
			res.IdKunjungan = defaultIdKunjungan
		} else if encryptedKunjungan, err := crypto.Encrypt(res.NoRawat, h.encryptionKey); err == nil {
			res.IdKunjungan = encryptedKunjungan
		}

		for j := range res.ResepDokter {
			rd := &res.ResepDokter[j]
			if encryptedObat, err := crypto.Encrypt(rd.KodeObat, h.encryptionKey); err == nil {
				rd.IdObat = encryptedObat
			}
		}

		for j := range res.ResepDokterRacikan {
			rdr := &res.ResepDokterRacikan[j]
			for k := range rdr.DetailRacikan {
				detail := &rdr.DetailRacikan[k]
				if encryptedObat, err := crypto.Encrypt(detail.KodeObat, h.encryptionKey); err == nil {
					detail.IdObat = encryptedObat
				}
			}
		}
	}
}
