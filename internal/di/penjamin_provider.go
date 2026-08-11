package di

import (
	"database/sql"

	"erm-dokter/internal/handler"
	"erm-dokter/internal/repository"
	"erm-dokter/internal/usecase"
	"erm-dokter/pkg/logger"
)

func providePenjamin(db *sql.DB, log *logger.Logger) *handler.PenjaminHandler {
	repo := repository.NewPenjaminRepository(db)
	uc := usecase.NewPenjaminUsecase(repo, log)
	return handler.NewPenjaminHandler(uc)
}
