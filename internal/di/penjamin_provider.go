package di

import (
	"database/sql"

	"erm-dokter/internal/penjamin"
	"erm-dokter/internal/pkg/logger"
)

func providePenjamin(db *sql.DB, log *logger.Logger) *penjamin.Handler {
	repo := penjamin.NewRepository(db)
	svc := penjamin.NewService(repo, log)
	return penjamin.NewHandler(svc)
}
