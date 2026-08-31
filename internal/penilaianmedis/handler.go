package penilaianmedis

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/middleware"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/response"
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

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	authMiddleware func(http.HandlerFunc) http.HandlerFunc,
	timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc,
) {
	mux.HandleFunc("GET /api/v1/penilaian-medis/referensi", authMiddleware(timeoutMiddleware(h.Referensi)))
	mux.HandleFunc("GET /api/v1/penilaian-medis/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailPenilaianMedisRalan)))
	mux.HandleFunc("GET /api/v1/penilaian-medis/ralan/pasien/{id_pasien}", authMiddleware(timeoutMiddleware(h.RiwayatPenilaianMedisRalan)))
	mux.HandleFunc("POST /api/v1/penilaian-medis/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.SimpanPenilaianMedisRalan)))
	mux.HandleFunc("PUT /api/v1/penilaian-medis/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.UpdatePenilaianMedisRalan)))
	mux.HandleFunc("DELETE /api/v1/penilaian-medis/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.HapusPenilaianMedisRalan)))
}

func (h *Handler) Referensi(w http.ResponseWriter, r *http.Request) {
	ref := h.service.Referensi(r.Context())
	response.Success(w, "Berhasil mengambil opsi referensi penilaian medis", ref)
}

func (h *Handler) DetailPenilaianMedisRalan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	detail, err := h.service.DetailPenilaianMedisRalan(r.Context(), noRawat)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil detail penilaian awal medis rawat jalan", detail)
}

func (h *Handler) RiwayatPenilaianMedisRalan(w http.ResponseWriter, r *http.Request) {
	idPasien := r.PathValue("id_pasien")
	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	riwayat, err := h.service.RiwayatPenilaianMedisRalanByNoRM(r.Context(), noRM)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil riwayat penilaian awal medis rawat jalan pasien", riwayat)
}

func (h *Handler) SimpanPenilaianMedisRalan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req SimpanPenilaianMedisRalanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	if req.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada payload tidak cocok dengan ID kunjungan"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil, err := h.service.SimpanPenilaianMedisRalan(r.Context(), kodeDokter, noRawat, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Created(w, "Berhasil menyimpan penilaian awal medis rawat jalan", hasil)
}

func (h *Handler) UpdatePenilaianMedisRalan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req UpdatePenilaianMedisRalanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil, err := h.service.UpdatePenilaianMedisRalan(r.Context(), kodeDokter, noRawat, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil memperbarui penilaian awal medis rawat jalan", hasil)
}

func (h *Handler) HapusPenilaianMedisRalan(w http.ResponseWriter, r *http.Request) {
	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if err := h.service.HapusPenilaianMedisRalan(r.Context(), kodeDokter, noRawat); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil menghapus penilaian awal medis rawat jalan", nil)
}
