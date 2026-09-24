package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/obat"
	"erm-dokter/internal/pkg/logger"
)

func provideObat(db *sql.DB, cfg *config.Config, log *logger.Logger) *obat.Handler {
	repo := obat.NewRepository(db)
	svc := obat.NewService(repo, log)
	return obat.NewHandler(svc, cfg.EncryptionKey)
}
