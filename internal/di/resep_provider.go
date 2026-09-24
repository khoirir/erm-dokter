package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/obat"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/resep"
)

func provideResep(db *sql.DB, cfg *config.Config, log *logger.Logger) *resep.Handler {
	repo := resep.NewRepository(db)
	rawatJalanRepo := rawatjalan.NewRepository(db)
	rawatJalanSvc := rawatjalan.NewService(rawatJalanRepo, log)
	rawatInapSvc := provideRawatInap(db, log)
	obatRepo := obat.NewRepository(db)
	obatSvc := obat.NewService(obatRepo, log)
	svc := resep.NewService(repo, rawatJalanSvc, rawatInapSvc, obatSvc, cfg.MaxEditRekamMedisJam, log)
	return resep.NewHandler(svc, cfg.EncryptionKey)
}
