package shared_test

import (
	"testing"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

func TestParseWaktu(t *testing.T) {
	tests := []struct {
		name        string
		tanggal     string
		jam         string
		expectError bool
		checkHour   int
		checkMin    int
		checkSec    int
	}{
		{
			name:        "Tanggal dan jam lengkap (HH:mm:ss)",
			tanggal:     "2026-08-21",
			jam:         "14:30:45",
			expectError: false,
			checkHour:   14,
			checkMin:    30,
			checkSec:    45,
		},
		{
			name:        "Tanggal dan jam tanpa detik (HH:mm)",
			tanggal:     "2026-08-21",
			jam:         "09:15",
			expectError: false,
			checkHour:   9,
			checkMin:    15,
			checkSec:    0,
		},
		{
			name:        "Tanggal saja tanpa jam (jam kosong)",
			tanggal:     "2026-08-21",
			jam:         "",
			expectError: false,
			checkHour:   0,
			checkMin:    0,
			checkSec:    0,
		},
		{
			name:        "Tanggal kosong",
			tanggal:     "",
			jam:         "12:00:00",
			expectError: true,
		},
		{
			name:        "Format tanggal salah",
			tanggal:     "21-08-2026",
			jam:         "12:00:00",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := shared.ParseWaktu(tt.tanggal, tt.jam)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError {
				if res.Hour() != tt.checkHour || res.Minute() != tt.checkMin || res.Second() != tt.checkSec {
					t.Errorf("expected time %02d:%02d:%02d, got %02d:%02d:%02d",
						tt.checkHour, tt.checkMin, tt.checkSec, res.Hour(), res.Minute(), res.Second())
				}
			}
		})
	}
}

func TestValidasiRentangTanggal(t *testing.T) {
	errs := make(apperror.ValidationError)
	shared.ValidasiRentangTanggal("2026-01-01,2026-01-31", errs)
	if len(errs) != 0 {
		t.Fatalf("expected 0 errors, got %+v", errs)
	}

	errsInvalid := make(apperror.ValidationError)
	shared.ValidasiRentangTanggal("2026-01-31,2026-01-01", errsInvalid)
	if _, exists := errsInvalid["tanggal"]; !exists {
		t.Error("expected validation error when tglAwal > tglAkhir")
	}
}
