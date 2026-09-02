package di

import (
	"database/sql"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
)

func provideBerkasDigital(db *sql.DB, cfg *config.Config, log *logger.Logger) *berkasdigital.Handler {
	repo := berkasdigital.NewRepository(db)
	svc := berkasdigital.NewService(repo, cfg.EncryptionKey, cfg.URLBerkasDigital, log)
	return berkasdigital.NewHandler(svc)
}
