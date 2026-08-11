package di

import "erm-dokter/internal/handler"

func provideHealth() *handler.HealthHandler {
	return handler.NewHealthHandler()
}
