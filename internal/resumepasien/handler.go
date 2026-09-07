package resumepasien

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/middleware"
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

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	authMiddleware func(http.HandlerFunc) http.HandlerFunc,
	timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc,
) {
	mux.HandleFunc("GET /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DetailResumePasien)))
	mux.HandleFunc("GET /api/v1/resume/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.RiwayatResumePasienByNoRM)))
	mux.HandleFunc("POST /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.SimpanResumePasien)))
	mux.HandleFunc("PUT /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.UpdateResumePasien)))
	mux.HandleFunc("DELETE /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.HapusResumePasien)))
}

func (h *Handler) parseStatusLanjut(r *http.Request) (shared.StatusLanjut, error) {
	status := shared.StatusLanjut(r.PathValue("status_lanjut"))
	if !status.IsValid() {
		return "", apperror.NewBusinessError("Status lanjut tidak valid")
	}
	if status != shared.StatusLanjutRawatJalan {
		return "", apperror.NewBusinessError("Resume pasien ini khusus untuk rawat jalan (Ralan)")
	}
	return status, nil
}

func (h *Handler) DetailResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	detail, err := h.service.DetailResumePasien(r.Context(), noRawat, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	detail.IdKunjungan = idKunjungan
	response.Success(w, "Berhasil mengambil detail resume pasien", detail)
}

func (h *Handler) RiwayatResumePasienByNoRM(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	idPasien := r.PathValue("id_pasien")
	noRM, err := crypto.Decrypt(idPasien, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID pasien tidak valid"))
		return
	}

	riwayat, err := h.service.RiwayatResumePasienByNoRM(r.Context(), noRM, statusLanjut)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	for i := range riwayat {
		if enc, err := crypto.Encrypt(riwayat[i].NoRawat, h.encryptionKey); err == nil {
			riwayat[i].IdKunjungan = enc
		}
	}

	response.Success(w, "Berhasil mengambil riwayat resume pasien", riwayat)
}

func (h *Handler) SimpanResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req SimpanResumePasienRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	if req.NoRawat != noRawat {
		apperror.HandleError(w, apperror.NewBusinessError("Nomor rawat pada payload tidak cocok dengan ID kunjungan"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil, err := h.service.SimpanResumePasien(r.Context(), noRawat, kodeDokter, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan
	response.Created(w, "Berhasil menyimpan resume pasien", hasil)
}

func (h *Handler) UpdateResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	idKunjungan := r.PathValue("id_kunjungan")
	noRawat, err := crypto.Decrypt(idKunjungan, h.encryptionKey)
	if err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("ID kunjungan tidak valid"))
		return
	}

	var req UpdateResumePasienRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	kodeDokter, err := middleware.GetKodeDokter(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil, err := h.service.UpdateResumePasien(r.Context(), kodeDokter, noRawat, statusLanjut, req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	hasil.IdKunjungan = idKunjungan
	response.Success(w, "Berhasil memperbarui resume pasien", hasil)
}

func (h *Handler) HapusResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

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

	if err := h.service.HapusResumePasien(r.Context(), kodeDokter, noRawat, statusLanjut); err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Berhasil menghapus resume pasien", nil)
}
