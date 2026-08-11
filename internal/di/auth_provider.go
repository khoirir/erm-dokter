package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
	"erm-dokter/internal/repository"
	"erm-dokter/internal/usecase"
	"erm-dokter/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func provideAuth(db *sql.DB, cfg *config.Config, validate *validator.Validate, log *logger.Logger) *handler.AuthHandler {
	repo := repository.NewAuthRepository(db, cfg.UserKey, cfg.PasswordKey)
	uc := usecase.NewAuthUsecase(repo, cfg.JWTSecret, validate, log)
	return handler.NewAuthHandler(uc)
}
