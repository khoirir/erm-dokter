package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/resep"
)

func provideResep(db *sql.DB, cfg *config.Config, log *logger.Logger) *resep.Handler {
	repo := resep.NewRepository(db)
	svc := resep.NewService(repo, log)
	return resep.NewHandler(svc, cfg.EncryptionKey)
}
