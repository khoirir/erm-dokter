package routes

import (
	"net/http"

	"erm-dokter/internal/handler"
)

func RegisterAuthRoutes(mux *http.ServeMux, h *handler.AuthHandler) {
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
}

