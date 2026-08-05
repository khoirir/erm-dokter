package routes

import (
	"database/sql"
	"net/http"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
)

func SetupRouter(db *sql.DB, cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	healthHandler := handler.NewHealthHandler()
	mux.HandleFunc("GET /health", healthHandler.HealthCheck)

	RegisterAuthRoutes(mux, db, cfg)
	RegisterRawatJalanRoutes(mux, db, cfg)

	return mux
}
