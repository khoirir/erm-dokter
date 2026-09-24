package di

import (
	"database/sql"

	"erm-dokter/internal/auth"
	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
)

func provideAuth(db *sql.DB, cfg *config.Config, log *logger.Logger) *auth.Handler {
	repo := auth.NewRepository(db, cfg.UserKey, cfg.PasswordKey)
	svc := auth.NewService(repo, cfg.JWTSecret, log)
	return auth.NewHandler(svc)
}
