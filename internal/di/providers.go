package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/routes"
	"erm-dokter/internal/shared/apperror"

	"github.com/go-playground/validator/v10"
)

func ProvideRouteConfig(db *sql.DB, cfg *config.Config, log *logger.Logger) *routes.RouteConfig {
	apperror.SetLogger(log)
	validate := validator.New()

	return routes.NewRouteConfig(
		provideHealth(),
		provideAuth(db, cfg, validate, log),
		providePenjamin(db, log),
		provideRawatJalan(db, cfg, log),
		cfg.JWTSecret,
	)
}
