package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/tindakan"
)

func provideTindakan(db *sql.DB, cfg *config.Config, log *logger.Logger) *tindakan.Handler {
	repo := tindakan.NewRepository(db)
	svc := tindakan.NewService(repo, cfg.EncryptionKey, log)
	return tindakan.NewHandler(svc)
}
