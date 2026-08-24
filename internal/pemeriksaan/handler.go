package pemeriksaan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	pemeriksaanService Service
	encryptionKey      string
}

func NewHandler(service Service, encryptionKey string) *Handler {
	return &Handler{
		pemeriksaanService: service,
		encryptionKey:      encryptionKey,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.HandlerFunc) http.HandlerFunc, timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/pemeriksaan/referensi-kesadaran", authMiddleware(timeoutMiddleware(h.GetDaftarKesadaran)))
	mux.HandleFunc("GET /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPemeriksaan)))
	mux.HandleFunc("GET /api/v1/pemeriksaan/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DaftarPemeriksaanByPasien)))
	mux.HandleFunc("GET /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}/{id_pemeriksaan}", authMiddleware(timeoutMiddleware(h.DetailPemeriksaan)))
	mux.HandleFunc("POST /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.SimpanPemeriksaan)))
	mux.HandleFunc("PUT /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}/{id_pemeriksaan}", authMiddleware(timeoutMiddleware(h.UpdatePemeriksaan)))
	mux.HandleFunc("DELETE /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}/{id_pemeriksaan}", authMiddleware(timeoutMiddleware(h.HapusPemeriksaan)))
}

func (h *Handler) DaftarPemeriksaan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	q := r.URL.Query()
	halaman, _ := strconv.Atoi(q.Get("halaman"))
	batas, _ := strconv.Atoi(q.Get("batas"))

	filter := FilterDaftarPemeriksaan{
		Tanggal: q.Get("tanggal"),
		Halaman: halaman,
		Batas:   batas,
	}

	daftarPemeriksaan, meta, err := h.pemeriksaanService.DaftarPemeriksaan(r.Context(), noRawat, shared.StatusLanjut(statusLanjut), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range daftarPemeriksaan {
		pemeriksaan := &daftarPemeriksaan[i]
		if encrypted, err := crypto.Encrypt(pemeriksaan.CompositeKey(), h.encryptionKey); err == nil {
			pemeriksaan.Id = encrypted
		}
		pemeriksaan.IdKunjungan = idKunjungan
	}

	response.SuccessWithMeta(w, "Berhasil mengambil daftar pemeriksaan", daftarPemeriksaan, meta)
}

func (h *Handler) DaftarPemeriksaanByPasien(w http.ResponseWriter, r *http.Request) {
	idPasien := r.PathValue("id_pasien")
	statusLanjut := r.PathValue("status_lanjut")

	noRekamMedis, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	q := r.URL.Query()
	halaman, _ := strconv.Atoi(q.Get("halaman"))
	batas, _ := strconv.Atoi(q.Get("batas"))

	filter := FilterDaftarPemeriksaan{
		Tanggal: q.Get("tanggal"),
		Halaman: halaman,
		Batas:   batas,
	}

	daftarPemeriksaan, meta, err := h.pemeriksaanService.DaftarPemeriksaanByRM(r.Context(), noRekamMedis, shared.StatusLanjut(statusLanjut), filter)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range daftarPemeriksaan {
		pemeriksaan := &daftarPemeriksaan[i]
		if encrypted, err := crypto.Encrypt(pemeriksaan.CompositeKey(), h.encryptionKey); err == nil {
			pemeriksaan.Id = encrypted
		}
		if encryptedIdKunjungan, err := crypto.Encrypt(pemeriksaan.NoRawat, h.encryptionKey); err == nil {
			pemeriksaan.IdKunjungan = encryptedIdKunjungan
		}
	}

	response.SuccessWithMeta(w, "Berhasil mengambil riwayat pemeriksaan pasien", daftarPemeriksaan, meta)
}

func (h *Handler) DetailPemeriksaan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")
	encryptedID := r.PathValue("id_pemeriksaan")

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid atau kadaluarsa"))
		return
	}

	decrypted, err := crypto.Decrypt(encryptedID, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pemeriksaan tidak valid atau kadaluarsa"))
		return
	}

	idPemeriksaan, err := ParseIdPemeriksaan(decrypted)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError(err.Error()))
		return
	}

	if idPemeriksaan.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("ID pemeriksaan tidak cocok dengan ID kunjungan"))
		return
	}

	pemeriksaan, err := h.pemeriksaanService.DetailPemeriksaan(r.Context(), idPemeriksaan, shared.StatusLanjut(statusLanjut))
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if pemeriksaan == nil {
		apperror.HandleError(w, apperror.NewNotFoundError("Detail pemeriksaan tidak ditemukan"))
		return
	}

	pemeriksaan.Id = encryptedID
	pemeriksaan.IdKunjungan = idKunjungan

	response.Success(w, "Berhasil mengambil detail pemeriksaan", pemeriksaan)
}

func (h *Handler) GetDaftarKesadaran(w http.ResponseWriter, r *http.Request) {
	daftarKesadaran := h.pemeriksaanService.GetDaftarKesadaran(r.Context())
	response.Success(w, "Berhasil mengambil daftar kesadaran", daftarKesadaran)
}

func (h *Handler) SimpanPemeriksaan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")

	noRawatURL, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid atau kadaluarsa"))
		return
	}

	var req SimpanPemeriksaanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	if req.NoRawat != noRawatURL {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada payload tidak cocok dengan ID kunjungan"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	pemeriksaan, err := h.pemeriksaanService.SimpanPemeriksaan(r.Context(), kodeDokter, shared.StatusLanjut(statusLanjut), req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if pemeriksaan != nil {
		if encrypted, err := crypto.Encrypt(pemeriksaan.CompositeKey(), h.encryptionKey); err == nil {
			pemeriksaan.Id = encrypted
		}
		pemeriksaan.IdKunjungan = idKunjungan
	}

	response.Created(w, "Berhasil menyimpan data pemeriksaan", pemeriksaan)
}

func (h *Handler) UpdatePemeriksaan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")
	encryptedIdPemeriksaan := r.PathValue("id_pemeriksaan")

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid atau kadaluarsa"))
		return
	}

	decryptedIdPememeriksaan, err := crypto.Decrypt(encryptedIdPemeriksaan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pemeriksaan tidak valid atau kadaluarsa"))
		return
	}

	idPemeriksaan, err := ParseIdPemeriksaan(decryptedIdPememeriksaan)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format ID pemeriksaan tidak valid"))
		return
	}

	if idPemeriksaan.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("ID pemeriksaan tidak cocok dengan ID kunjungan"))
		return
	}

	var req UpdatePemeriksaanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	pemeriksaan, err := h.pemeriksaanService.UpdatePemeriksaan(r.Context(), kodeDokter, idPemeriksaan, shared.StatusLanjut(statusLanjut), req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if pemeriksaan != nil {
		if encrypted, err := crypto.Encrypt(pemeriksaan.CompositeKey(), h.encryptionKey); err == nil {
			pemeriksaan.Id = encrypted
		}
		pemeriksaan.IdKunjungan = idKunjungan
	}

	response.Success(w, "Berhasil memperbarui data pemeriksaan", pemeriksaan)
}

func (h *Handler) HapusPemeriksaan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	statusLanjut := r.PathValue("status_lanjut")
	encryptedIdPemeriksaan := r.PathValue("id_pemeriksaan")

	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid atau kadaluarsa"))
		return
	}

	decryptedIdPememeriksaan, err := crypto.Decrypt(encryptedIdPemeriksaan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pemeriksaan tidak valid atau kadaluarsa"))
		return
	}

	idPemeriksaan, err := ParseIdPemeriksaan(decryptedIdPememeriksaan)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format ID pemeriksaan tidak valid"))
		return
	}

	if idPemeriksaan.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("ID pemeriksaan tidak cocok dengan ID kunjungan"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if err := h.pemeriksaanService.HapusPemeriksaan(r.Context(), kodeDokter, idPemeriksaan, shared.StatusLanjut(statusLanjut)); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil menghapus data pemeriksaan", nil)
}



