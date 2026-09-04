package di

import (
	"database/sql"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
)

func provideRawatInap(db *sql.DB, log *logger.Logger) rawatinap.Service {
	repo := rawatinap.NewRepository(db)
	return rawatinap.NewService(repo, log)
}
