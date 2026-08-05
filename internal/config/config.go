package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     string
	AppName     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
	JWTSecret   string
	UserKey     string
	PasswordKey string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO]\tFile .env tidak ditemukan, membaca dari environment variable sistem")
	}
	return &Config{
		AppEnv:      os.Getenv("APP_ENV"),
		AppPort:     getEnvOrDefault("APP_PORT", "8080"),
		AppName:     os.Getenv("APP_NAME"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPass:      os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		UserKey:     os.Getenv("USER_KEY"),
		PasswordKey: os.Getenv("PASSWORD_KEY"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
