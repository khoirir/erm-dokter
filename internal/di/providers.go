package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
	"erm-dokter/internal/routes"
	"erm-dokter/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func ProvideRouteConfig(db *sql.DB, cfg *config.Config, log *logger.Logger) *routes.RouteConfig {
	handler.SetLogger(log)
	validate := validator.New()

	return routes.NewRouteConfig(
		provideHealth(),
		provideAuth(db, cfg, validate, log),
		providePenjamin(db, log),
		provideRawatJalan(db, cfg, log),
		cfg.JWTSecret,
	)
}
