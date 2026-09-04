package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/penilaianmedis"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
)

func providePenilaianMedis(db *sql.DB, cfg *config.Config, log *logger.Logger) *penilaianmedis.Handler {
	repo := penilaianmedis.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	rawatInapRepo := rawatinap.NewRepository(db)
	rawatInapSvc := rawatinap.NewService(rawatInapRepo, log)
	svc := penilaianmedis.NewService(repo, rawatJalanSvc, rawatInapSvc, 48, log)
	return penilaianmedis.NewHandler(svc, cfg.EncryptionKey)
}
