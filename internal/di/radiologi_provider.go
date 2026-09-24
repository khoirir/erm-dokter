package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/radiologi"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/tindakan"
)

func provideRadiologi(db *sql.DB, cfg *config.Config, log *logger.Logger) *radiologi.Handler {
	repo := radiologi.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	rawatInapSvc := provideRawatInap(db, log)
	tindakanRepo := tindakan.NewRepository(db)
	tindakanSvc := tindakan.NewService(tindakanRepo, log)

	svc := radiologi.NewService(
		repo,
		rawatJalanSvc,
		rawatInapSvc,
		tindakanSvc,
		cfg.MaxEditRekamMedisJam,
		log,
	)
	return radiologi.NewHandler(svc, cfg.EncryptionKey)
}
