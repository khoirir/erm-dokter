package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
)

func provideRawatInap(db *sql.DB, log *logger.Logger) rawatinap.Service {
	repo := rawatinap.NewRepository(db)
	return rawatinap.NewService(repo, log)
}

func provideRawatInapHandler(db *sql.DB, cfg *config.Config, log *logger.Logger) *rawatinap.Handler {
	svc := provideRawatInap(db, log)
	return rawatinap.NewHandler(svc, cfg.EncryptionKey)
}
