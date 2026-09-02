package di

import (
	"database/sql"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/config"
	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/pkg/logger"
)

func provideLaboratorium(db *sql.DB, cfg *config.Config, log *logger.Logger) *laboratorium.Handler {
	repo := laboratorium.NewRepository(db)
	berkasRepo := berkasdigital.NewRepository(db)
	berkasSvc := berkasdigital.NewService(berkasRepo, cfg.EncryptionKey, cfg.URLBerkasDigital, log)
	svc := laboratorium.NewService(repo, berkasSvc, cfg.EncryptionKey, cfg.URLBerkasDigital, cfg.KodeBerkasLabPK, cfg.KodeBerkasLabPA, cfg.KodeBerkasLabMB, log)
	return laboratorium.NewHandler(svc)
}
