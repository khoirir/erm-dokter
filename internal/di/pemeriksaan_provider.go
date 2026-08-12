package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/pemeriksaan"
)

func providePemeriksaan(db *sql.DB, cfg *config.Config, log *logger.Logger) *pemeriksaan.Handler {
	repo := pemeriksaan.NewRepository(db)
	svc := pemeriksaan.NewService(repo, log)
	return pemeriksaan.NewHandler(svc, cfg.EncryptionKey)
}
