package di

import (
	"database/sql"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/config"
	"erm-dokter/internal/laboratorium"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/tindakan"
)

func provideLaboratorium(db *sql.DB, cfg *config.Config, log *logger.Logger) *laboratorium.Handler {
	repo := laboratorium.NewRepository(db)
	berkasRepo := berkasdigital.NewRepository(db)
	berkasSvc := berkasdigital.NewService(berkasRepo, cfg.URLBerkasDigital, log)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	tindakanRepo := tindakan.NewRepository(db)
	tindakanSvc := tindakan.NewService(tindakanRepo, log)

	svc := laboratorium.NewService(
		repo,
		berkasSvc,
		rawatJalanSvc,
		tindakanSvc,
		cfg.MaxEditRekamMedisJam,
		cfg.URLBerkasDigital,
		cfg.KodeBerkasLabPK,
		cfg.KodeBerkasLabPA,
		cfg.KodeBerkasLabMB,
		log,
	)
	return laboratorium.NewHandler(svc, cfg.EncryptionKey)
}
