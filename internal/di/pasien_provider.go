package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pasien"
	"erm-dokter/internal/pkg/logger"
)

func providePasien(db *sql.DB, cfg *config.Config, log *logger.Logger) *pasien.Handler {
	repo := pasien.NewRepository(db)
	svc := pasien.NewService(repo, log)
	return pasien.NewHandler(svc, cfg.EncryptionKey)
}
