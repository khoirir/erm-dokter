package routes

import (
	"database/sql"
	"net/http"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
	"erm-dokter/internal/repository"
	"erm-dokter/internal/usecase"
)

func RegisterAuthRoutes(mux *http.ServeMux, db *sql.DB, cfg *config.Config) {
	repo := repository.NewAuthRepository(db, cfg.UserKey, cfg.PasswordKey)
	uc := usecase.NewAuthUsecase(repo, cfg.JWTSecret)
	h := handler.NewAuthHandler(uc)

	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
}
