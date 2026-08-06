package routes

import (
	"database/sql"
	"net/http"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
	"erm-dokter/internal/middleware"
	"erm-dokter/internal/repository"
	"erm-dokter/internal/usecase"
)

func RegisterPenjaminRoutes(mux *http.ServeMux, db *sql.DB, cfg *config.Config) {
	repo := repository.NewPenjaminRepository(db)
	uc := usecase.NewPenjaminUsecase(repo)
	h := handler.NewPenjaminHandler(uc)

	authMiddleware := middleware.JWTMiddleware(cfg.JWTSecret)

	mux.HandleFunc("GET /api/v1/penjamin", authMiddleware(h.DaftarPenjamin))
}