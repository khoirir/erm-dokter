package routes

import (
	"net/http"

	"erm-dokter/internal/handler"
	"erm-dokter/internal/middleware"
)

type Handlers struct {
	Health     *handler.HealthHandler
	Auth       *handler.AuthHandler
	Penjamin   *handler.PenjaminHandler
	RawatJalan *handler.RawatJalanHandler
}

func SetupRouter(handlers *Handlers, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()
	authMW := middleware.JWTMiddleware(jwtSecret)

	mux.HandleFunc("GET /health", handlers.Health.HealthCheck)

	RegisterAuthRoutes(mux, handlers.Auth)
	RegisterPenjaminRoutes(mux, handlers.Penjamin, authMW)
	RegisterRawatJalanRoutes(mux, handlers.RawatJalan, authMW)

	return mux
}

