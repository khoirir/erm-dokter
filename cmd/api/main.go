package main

import (
	"fmt"
	"net/http"

	"erm-dokter/internal/config"
	"erm-dokter/internal/routes"
	"erm-dokter/pkg/database"
	"erm-dokter/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New()

	db, err := database.InitMySQL(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)
	if err != nil {
		log.Error("Database MySQL connection failed: %v", err)
	} else {
		defer db.Close()
	}

	mux := routes.SetupRouter(db, cfg)

	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Info("Starting server '%s' (%s mode) on %s...", cfg.AppName, cfg.AppEnv, serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Error("Failed to start server: %v", err)
	}
}
