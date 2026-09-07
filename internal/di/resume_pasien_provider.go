package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/resumepasien"
)

func provideResumePasien(db *sql.DB, cfg *config.Config, log *logger.Logger) *resumepasien.Handler {
	repo := resumepasien.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	rawatInapRepo := rawatinap.NewRepository(db)
	rawatInapSvc := rawatinap.NewService(rawatInapRepo, log)
	svc := resumepasien.NewService(repo, rawatJalanSvc, rawatInapSvc, 48, log)
	return resumepasien.NewHandler(svc, cfg.EncryptionKey)
}
