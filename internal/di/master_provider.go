package di

import (
	"database/sql"

	"erm-dokter/internal/master"
	"erm-dokter/internal/pkg/logger"
)

func provideMaster(db *sql.DB, log *logger.Logger) *master.Handler {
	repo := master.NewRepository(db)
	svc := master.NewService(repo, log)
	return master.NewHandler(svc)
}
