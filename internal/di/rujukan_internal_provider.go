package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/rujukaninternal"
)

func provideRujukanInternal(db *sql.DB, cfg *config.Config, log *logger.Logger) *rujukaninternal.Handler {
	repo := rujukaninternal.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	svc := rujukaninternal.NewService(repo, rawatJalanSvc, log, cfg.EncryptionKey)
	return rujukaninternal.NewHandler(svc, cfg.EncryptionKey)
}
