package di

import (
	"database/sql"

	"erm-dokter/internal/auth"
	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func provideAuth(db *sql.DB, cfg *config.Config, validate *validator.Validate, log *logger.Logger) *auth.Handler {
	repo := auth.NewRepository(db, cfg.UserKey, cfg.PasswordKey)
	svc := auth.NewService(repo, cfg.JWTSecret, validate, log)
	return auth.NewHandler(svc)
}
