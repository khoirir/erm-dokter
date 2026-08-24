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
		providePenjamin(db, log),
		provideRawatJalan(db, cfg, log),
		providePemeriksaan(db, cfg, log),
		cfg.JWTSecret,
	)
}
