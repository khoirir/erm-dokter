package di

import (
	"database/sql"

	"erm-dokter/internal/config"
	"erm-dokter/internal/handler"
	"erm-dokter/internal/repository"
	"erm-dokter/internal/routes"
	"erm-dokter/internal/usecase"
)

func ProvideHandlers(db *sql.DB, cfg *config.Config) *routes.Handlers {
	return &routes.Handlers{
		Health:     handler.NewHealthHandler(),
		Auth:       provideAuth(db, cfg),
		Penjamin:   providePenjamin(db),
		RawatJalan: provideRawatJalan(db, cfg),
	}
}

func provideAuth(db *sql.DB, cfg *config.Config) *handler.AuthHandler {
	repo := repository.NewAuthRepository(db, cfg.UserKey, cfg.PasswordKey)
	uc := usecase.NewAuthUsecase(repo, cfg.JWTSecret)
	return handler.NewAuthHandler(uc)
}

func providePenjamin(db *sql.DB) *handler.PenjaminHandler {
	repo := repository.NewPenjaminRepository(db)
	uc := usecase.NewPenjaminUsecase(repo)
	return handler.NewPenjaminHandler(uc)
}

func provideRawatJalan(db *sql.DB, cfg *config.Config) *handler.RawatJalanHandler {
	repo := repository.NewRawatJalanRepository(db)
	uc := usecase.NewRawatJalanUsecase(repo)
	return handler.NewRawatJalanHandler(uc, cfg.EncryptionKey)
}
