package resumepasien

import (
	"net/http"
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
	mux.HandleFunc("GET /api/v1/resume/ralan/referensi", authMiddleware(timeoutMiddleware(h.ReferensiRalan)))
	mux.HandleFunc("GET /api/v1/resume/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailResumePasienRalan)))
	mux.HandleFunc("GET /api/v1/resume/ralan/pasien/{id_pasien}", authMiddleware(timeoutMiddleware(h.RiwayatResumePasienRalanByNoRM)))
	mux.HandleFunc("POST /api/v1/resume/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.SimpanResumePasienRalan)))
	mux.HandleFunc("PUT /api/v1/resume/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.UpdateResumePasienRalan)))
	mux.HandleFunc("DELETE /api/v1/resume/ralan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.HapusResumePasienRalan)))

	mux.HandleFunc("GET /api/v1/resume/ranap/referensi", authMiddleware(timeoutMiddleware(h.ReferensiRanap)))
	mux.HandleFunc("GET /api/v1/resume/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailResumePasienRanap)))
	mux.HandleFunc("GET /api/v1/resume/ranap/pasien/{id_pasien}", authMiddleware(timeoutMiddleware(h.RiwayatResumePasienRanapByNoRM)))
	mux.HandleFunc("POST /api/v1/resume/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.SimpanResumePasienRanap)))
	mux.HandleFunc("PUT /api/v1/resume/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.UpdateResumePasienRanap)))
	mux.HandleFunc("DELETE /api/v1/resume/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.HapusResumePasienRanap)))
}