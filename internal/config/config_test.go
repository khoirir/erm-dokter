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
}
