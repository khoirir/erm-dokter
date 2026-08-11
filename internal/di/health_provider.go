package di

import "erm-dokter/internal/health"

func provideHealth() *health.Handler {
	return health.NewHandler()
}
