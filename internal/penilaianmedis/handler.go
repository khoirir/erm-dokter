package penilaianmedis

import (
	"net/http"

	"erm-dokter/internal/pkg/response"
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

	mux.HandleFunc("GET /api/v1/penilaian-medis/igd/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailPenilaianMedisIGD)))
	mux.HandleFunc("GET /api/v1/penilaian-medis/igd/pasien/{id_pasien}", authMiddleware(timeoutMiddleware(h.RiwayatPenilaianMedisIGD)))
	mux.HandleFunc("POST /api/v1/penilaian-medis/igd/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.SimpanPenilaianMedisIGD)))
	mux.HandleFunc("PUT /api/v1/penilaian-medis/igd/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.UpdatePenilaianMedisIGD)))
	mux.HandleFunc("DELETE /api/v1/penilaian-medis/igd/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.HapusPenilaianMedisIGD)))

	mux.HandleFunc("GET /api/v1/penilaian-medis/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailPenilaianMedisRanap)))
	mux.HandleFunc("GET /api/v1/penilaian-medis/ranap/pasien/{id_pasien}", authMiddleware(timeoutMiddleware(h.RiwayatPenilaianMedisRanap)))
	mux.HandleFunc("POST /api/v1/penilaian-medis/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.SimpanPenilaianMedisRanap)))
	mux.HandleFunc("PUT /api/v1/penilaian-medis/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.UpdatePenilaianMedisRanap)))
	mux.HandleFunc("DELETE /api/v1/penilaian-medis/ranap/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.HapusPenilaianMedisRanap)))

	mux.HandleFunc("GET /api/v1/penilaian-medis/ralan-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailPenilaianMedisRalanKandungan)))
	mux.HandleFunc("GET /api/v1/penilaian-medis/ralan-kandungan/pasien/{id_pasien}", authMiddleware(timeoutMiddleware(h.RiwayatPenilaianMedisRalanKandungan)))
	mux.HandleFunc("POST /api/v1/penilaian-medis/ralan-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.SimpanPenilaianMedisRalanKandungan)))
	mux.HandleFunc("PUT /api/v1/penilaian-medis/ralan-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.UpdatePenilaianMedisRalanKandungan)))
	mux.HandleFunc("DELETE /api/v1/penilaian-medis/ralan-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.HapusPenilaianMedisRalanKandungan)))

	mux.HandleFunc("GET /api/v1/penilaian-medis/ranap-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.DetailPenilaianMedisRanapKandungan)))
	mux.HandleFunc("GET /api/v1/penilaian-medis/ranap-kandungan/pasien/{id_pasien}", authMiddleware(timeoutMiddleware(h.RiwayatPenilaianMedisRanapKandungan)))
	mux.HandleFunc("POST /api/v1/penilaian-medis/ranap-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.SimpanPenilaianMedisRanapKandungan)))
	mux.HandleFunc("PUT /api/v1/penilaian-medis/ranap-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.UpdatePenilaianMedisRanapKandungan)))
	mux.HandleFunc("DELETE /api/v1/penilaian-medis/ranap-kandungan/{id_kunjungan}", authMiddleware(timeoutMiddleware(h.HapusPenilaianMedisRanapKandungan)))
}

func (h *Handler) Referensi(w http.ResponseWriter, r *http.Request) {
	ref := h.service.Referensi(r.Context())
	response.Success(w, "Berhasil mengambil opsi referensi penilaian medis", ref)
}
