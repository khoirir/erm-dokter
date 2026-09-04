package laboratorium

import (
	"net/http"
	"strconv"
	"strings"

	"erm-dokter/internal/berkasdigital"
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
	mux.HandleFunc("GET /api/v1/laboratorium/{kategori}/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarHasilLab)))
	mux.HandleFunc("GET /api/v1/laboratorium/{kategori}/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarHasilLabByRM)))
	mux.HandleFunc("GET /api/v1/laboratorium/{kategori}/{id_kunjungan}/{status_lanjut}/{id_hasil}", authMiddleware(timeoutMiddleware(h.DetailHasilLab)))

	mux.HandleFunc("POST /api/v1/laboratorium/pk/permintaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.SimpanPermintaanLabPK)))
	mux.HandleFunc("GET /api/v1/laboratorium/pk/permintaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPermintaanLabPK)))
	mux.HandleFunc("GET /api/v1/laboratorium/pk/permintaan/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPermintaanLabPKByRM)))
	mux.HandleFunc("GET /api/v1/laboratorium/pk/permintaan/{id_kunjungan}/{status_lanjut}/{id_permintaan}", authMiddleware(timeoutMiddleware(h.DetailPermintaanLabPK)))
	mux.HandleFunc("DELETE /api/v1/laboratorium/pk/permintaan/{id_kunjungan}/{status_lanjut}/{id_permintaan}", authMiddleware(timeoutMiddleware(h.HapusPermintaanLabPK)))

	mux.HandleFunc("POST /api/v1/laboratorium/pa/permintaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.SimpanPermintaanLabPA)))
	mux.HandleFunc("GET /api/v1/laboratorium/pa/permintaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPermintaanLabPA)))
	mux.HandleFunc("GET /api/v1/laboratorium/pa/permintaan/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPermintaanLabPAByRM)))
	mux.HandleFunc("GET /api/v1/laboratorium/pa/permintaan/{id_kunjungan}/{status_lanjut}/{id_permintaan}", authMiddleware(timeoutMiddleware(h.DetailPermintaanLabPA)))
	mux.HandleFunc("DELETE /api/v1/laboratorium/pa/permintaan/{id_kunjungan}/{status_lanjut}/{id_permintaan}", authMiddleware(timeoutMiddleware(h.HapusPermintaanLabPA)))
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

	data, meta, err := h.service.GetRiwayatLabKunjungan(r.Context(), kat, noRawat, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if data != nil {
		h.encryptHasilLabList(data.HasilPemeriksaan)
		h.encryptBerkasDigitalList(data.BerkasDigital)
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

	data, meta, err := h.service.GetRiwayatLabPasien(r.Context(), kat, noRM, statusLanjut, filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptHasilLabList(data)
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

	kunjunganNoRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	idHasil := strings.TrimSpace(r.PathValue("id_hasil"))
	if idHasil == "" {
		apperror.HandleError(w, apperror.NewBusinessError("ID hasil laboratorium wajib diisi"))
		return
	}

	plainId, err := crypto.Decrypt(idHasil, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID hasil laboratorium tidak valid"))
		return
	}

	parts := strings.Split(plainId, "~")
	if len(parts) != 4 {
		apperror.HandleError(w, apperror.NewBusinessError("Format ID hasil laboratorium tidak valid"))
		return
	}

	noRawat := parts[0]
	kodeTindakan := parts[1]
	tanggalPeriksa := parts[2]
	jamPeriksa := parts[3]

	if noRawat != kunjunganNoRawat {
		apperror.HandleError(w, apperror.NewBusinessError("ID hasil laboratorium tidak sesuai dengan kunjungan pasien"))
		return
	}

	data, err := h.service.GetDetailHasilLab(r.Context(), kat, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	h.encryptHasilLab(data)
	response.Success(w, "Berhasil mengambil detail hasil laboratorium", data)
}

func (h *Handler) encryptHasilLab(item *HasilLaboratorium) {
	if item == nil {
		return
	}
	item.Id, _ = crypto.Encrypt(item.NoRawat+"~"+item.KodeTindakan+"~"+item.TanggalPeriksa+"~"+item.JamPeriksa, h.encryptionKey)
	item.IdKunjungan, _ = crypto.Encrypt(item.NoRawat, h.encryptionKey)
	for i := range item.DetailPK {
		if item.DetailPK[i].IdTemplate != "" {
			item.DetailPK[i].IdTemplate, _ = crypto.Encrypt(item.DetailPK[i].IdTemplate, h.encryptionKey)
		}
	}
}

func (h *Handler) encryptHasilLabList(items []HasilLaboratorium) {
	for i := range items {
		h.encryptHasilLab(&items[i])
	}
}

func (h *Handler) encryptBerkasDigitalList(list []berkasdigital.BerkasDigital) {
	for i := range list {
		if list[i].IdBerkas != "" {
			encId, err := crypto.Encrypt(list[i].IdBerkas, h.encryptionKey)
			if err == nil {
				list[i].IdBerkas = encId
				list[i].UrlBerkas = "/api/v1/berkas-digital/" + encId
			}
		}
	}
}
