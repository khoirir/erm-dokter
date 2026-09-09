package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/routes"
	"erm-dokter/internal/shared/apperror"
)

func ProvideRouteConfig(db *sql.DB, cfg *config.Config, log *logger.Logger) *routes.RouteConfig {
	apperror.SetLogger(log)

	return routes.NewRouteConfig(
		provideHealth(),
		provideDocs(),
		provideAuth(db, cfg, log),
		provideMaster(db, log),
		provideObat(db, cfg, log),
		provideRawatJalan(db, cfg, log),
		provideRawatInapHandler(db, cfg, log),
		providePemeriksaan(db, cfg, log),
		provideResep(db, cfg, log),
		provideRujukanInternal(db, cfg, log),
		providePenilaianMedis(db, cfg, log),
		provideTindakan(db, cfg, log),
		provideLaboratorium(db, cfg, log),
		provideRadiologi(db, cfg, log),
		provideBerkasDigital(db, cfg, log),
		provideResumePasien(db, cfg, log),
		providePasien(db, cfg, log),
		cfg.JWTSecret,
		cfg.ServiceAPIKey,
	)
}
