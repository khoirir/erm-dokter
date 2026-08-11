package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
	"erm-dokter/internal/repository"
	"erm-dokter/internal/usecase"
	"erm-dokter/pkg/logger"
)

func provideRawatJalan(db *sql.DB, cfg *config.Config, log *logger.Logger) *handler.RawatJalanHandler {
	repo := repository.NewRawatJalanRepository(db)
	uc := usecase.NewRawatJalanUsecase(repo, log)
	return handler.NewRawatJalanHandler(uc, cfg.EncryptionKey)
}
