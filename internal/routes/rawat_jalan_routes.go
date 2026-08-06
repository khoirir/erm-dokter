package routes

import (
	"net/http"

	"erm-dokter/internal/handler"
)

func RegisterRawatJalanRoutes(mux *http.ServeMux, h *handler.RawatJalanHandler, authMW func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/rawat-jalan/referensi-filter", authMW(h.GetReferensiFilter))
	mux.HandleFunc("POST /api/v1/rawat-jalan/antrean", authMW(h.DaftarAntreanDokter))
	mux.HandleFunc("GET /api/v1/rawat-jalan/detail/{no_rawat}", authMW(h.DetailKunjungan))
}

