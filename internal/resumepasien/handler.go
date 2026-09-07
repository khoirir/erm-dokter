package resumepasien

import (
	"net/http"

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
	mux.HandleFunc("GET /api/v1/resume/referensi/{status_lanjut}", authMiddleware(timeoutMiddleware(h.Referensi)))
	mux.HandleFunc("GET /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.DetailResumePasien)))
	mux.HandleFunc("GET /api/v1/resume/pasien/{id_pasien}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.RiwayatResumePasienByNoRM)))
	mux.HandleFunc("POST /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.SimpanResumePasien)))
	mux.HandleFunc("PUT /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.UpdateResumePasien)))
	mux.HandleFunc("DELETE /api/v1/resume/{id_kunjungan}/{status_lanjut}", authMiddleware(timeoutMiddleware(h.HapusResumePasien)))
}

func (h *Handler) parseStatusLanjut(r *http.Request) (shared.StatusLanjut, error) {
	status := shared.StatusLanjut(r.PathValue("status_lanjut"))
	if !status.IsValid() {
		return "", apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)")
	}
	return status, nil
}

func (h *Handler) Referensi(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		ref := h.service.ReferensiRalan(r.Context())
		response.Success(w, "Berhasil mengambil opsi referensi resume ralan", ref)
	case shared.StatusLanjutRawatInap:
		ref := h.service.ReferensiRanap(r.Context())
		response.Success(w, "Berhasil mengambil opsi referensi resume ranap", ref)
	}
}

func (h *Handler) DetailResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		h.DetailResumePasienRalan(w, r)
	case shared.StatusLanjutRawatInap:
		h.DetailResumePasienRanap(w, r)
	}
}

func (h *Handler) RiwayatResumePasienByNoRM(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		h.RiwayatResumePasienRalanByNoRM(w, r)
	case shared.StatusLanjutRawatInap:
		h.RiwayatResumePasienRanapByNoRM(w, r)
	}
}

func (h *Handler) SimpanResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		h.SimpanResumePasienRalan(w, r)
	case shared.StatusLanjutRawatInap:
		h.SimpanResumePasienRanap(w, r)
	}
}

func (h *Handler) UpdateResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		h.UpdateResumePasienRalan(w, r)
	case shared.StatusLanjutRawatInap:
		h.UpdateResumePasienRanap(w, r)
	}
}

func (h *Handler) HapusResumePasien(w http.ResponseWriter, r *http.Request) {
	statusLanjut, err := h.parseStatusLanjut(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	switch statusLanjut {
	case shared.StatusLanjutRawatJalan:
		h.HapusResumePasienRalan(w, r)
	case shared.StatusLanjutRawatInap:
		h.HapusResumePasienRanap(w, r)
	}
}
