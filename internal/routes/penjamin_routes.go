package routes

import (
	"net/http"

	"erm-dokter/internal/handler"
)

func RegisterPenjaminRoutes(mux *http.ServeMux, h *handler.PenjaminHandler, authMW func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /api/v1/penjamin", authMW(h.DaftarPenjamin))
}