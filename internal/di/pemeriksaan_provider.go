package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pemeriksaan"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
)

func providePemeriksaan(db *sql.DB, cfg *config.Config, log *logger.Logger) *pemeriksaan.Handler {
	repo := pemeriksaan.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	rawatInapRepo := rawatinap.NewRepository(db)
	rawatInapSvc := rawatinap.NewService(rawatInapRepo, log)
	svc := pemeriksaan.NewService(repo, rawatJalanSvc, rawatInapSvc, cfg.MaxEditRekamMedisJam, log)
	return pemeriksaan.NewHandler(svc, cfg.EncryptionKey)
}
