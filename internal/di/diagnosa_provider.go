package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/diagnosa"
	"erm-dokter/internal/master"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
)

func provideDiagnosa(db *sql.DB, cfg *config.Config, log *logger.Logger) *diagnosa.Handler {
	repo := diagnosa.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	rawatInapRepo := rawatinap.NewRepository(db)
	rawatInapSvc := rawatinap.NewService(rawatInapRepo, log)
	masterRepo := master.NewRepository(db)
	masterSvc := master.NewService(masterRepo, log)
	svc := diagnosa.NewService(repo, rawatJalanSvc, rawatInapSvc, masterSvc, cfg.MaxEditRekamMedisJam, log)
	return diagnosa.NewHandler(svc, cfg.EncryptionKey)
}
