package config_test

import (
	"os"
	"testing"

	"erm-dokter/internal/config"
)

func TestConfigLoad(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "9999")
	os.Setenv("MAX_EDIT_REKAM_MEDIS_JAM", "72")
	os.Setenv("KODE_BERKAS_LAB_PK", "005, 006")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("MAX_EDIT_REKAM_MEDIS_JAM")
		os.Unsetenv("KODE_BERKAS_LAB_PK")
	}()

	cfg := config.Load()

	if cfg.AppEnv != "test" {
		t.Errorf("Expected AppEnv 'test', got '%s'", cfg.AppEnv)
	}
	if cfg.AppPort != "9999" {
		t.Errorf("Expected AppPort '9999', got '%s'", cfg.AppPort)
	}
	if cfg.MaxEditRekamMedisJam != 72 {
		t.Errorf("Expected MaxEditRekamMedisJam 72, got %d", cfg.MaxEditRekamMedisJam)
	}
	if len(cfg.KodeBerkasLabPK) != 2 || cfg.KodeBerkasLabPK[0] != "005" || cfg.KodeBerkasLabPK[1] != "006" {
		t.Errorf("Unexpected KodeBerkasLabPK: %v", cfg.KodeBerkasLabPK)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("Expected default LogFormat 'json', got '%s'", cfg.LogFormat)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected default LogLevel 'info', got '%s'", cfg.LogLevel)
	}
	if cfg.LogFilePath != "logs/app.log" {
		t.Errorf("Expected default LogFilePath 'logs/app.log', got '%s'", cfg.LogFilePath)
	}
	if !cfg.LoginRateLimitEnabled {
		t.Errorf("Expected default LoginRateLimitEnabled true, got %v", cfg.LoginRateLimitEnabled)
	}
	if cfg.LoginRateLimitRate != 10 {
		t.Errorf("Expected default LoginRateLimitRate 10, got %d", cfg.LoginRateLimitRate)
	}
	if cfg.LoginRateLimitWindowMinutes != 1 {
		t.Errorf("Expected default LoginRateLimitWindowMinutes 1, got %d", cfg.LoginRateLimitWindowMinutes)
	}
}

func TestConfigLoad_CustomRateLimit(t *testing.T) {
	os.Setenv("LOGIN_RATE_LIMIT_ENABLED", "false")
	os.Setenv("LOGIN_RATE_LIMIT_RATE", "5")
	os.Setenv("LOGIN_RATE_LIMIT_WINDOW_MINUTES", "2")

	defer func() {
		os.Unsetenv("LOGIN_RATE_LIMIT_ENABLED")
		os.Unsetenv("LOGIN_RATE_LIMIT_RATE")
		os.Unsetenv("LOGIN_RATE_LIMIT_WINDOW_MINUTES")
	}()

	cfg := config.Load()

	if cfg.LoginRateLimitEnabled {
		t.Errorf("Expected LoginRateLimitEnabled false, got %v", cfg.LoginRateLimitEnabled)
	}
	if cfg.LoginRateLimitRate != 5 {
		t.Errorf("Expected LoginRateLimitRate 5, got %d", cfg.LoginRateLimitRate)
	}
	if cfg.LoginRateLimitWindowMinutes != 2 {
		t.Errorf("Expected LoginRateLimitWindowMinutes 2, got %d", cfg.LoginRateLimitWindowMinutes)
	}
}
