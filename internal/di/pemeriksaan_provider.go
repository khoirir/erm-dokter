package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
)

func providePemeriksaan(db *sql.DB, cfg *config.Config, log *logger.Logger) *pemeriksaan.Handler {
	repo := pemeriksaan.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	svc := pemeriksaan.NewService(repo, rawatJalanRepo, log)
	return pemeriksaan.NewHandler(svc, cfg.EncryptionKey)
}
