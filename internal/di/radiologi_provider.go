package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/radiologi"
)

func provideRadiologi(db *sql.DB, cfg *config.Config, log *logger.Logger) *radiologi.Handler {
	repo := radiologi.NewRepository(db)
	svc := radiologi.NewService(repo, log)
	return radiologi.NewHandler(svc, cfg.EncryptionKey)
}
