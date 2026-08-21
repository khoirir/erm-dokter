package shared_test

import (
	"testing"
	"time"

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

func TestValidasiBatasWaktuRekamMedis(t *testing.T) {
	now := time.Now()

	// 1. Data baru (1 jam yang lalu) - harus lolos
	recent := now.Add(-1 * time.Hour)
	err := shared.ValidasiBatasWaktuRekamMedis(recent.Format("2006-01-02"), recent.Format("15:04:05"), 48, "dihapus")
	if err != nil {
		t.Fatalf("expected nil error for recent record, got: %v", err)
	}

	// 2. Data tanpa jam (hanya tanggal hari ini) - harus lolos
	err = shared.ValidasiBatasWaktuRekamMedis(now.Format("2006-01-02"), "", 48, "diubah")
	if err != nil {
		t.Fatalf("expected nil error for date-only today, got: %v", err)
	}

	// 3. Data lama (50 jam yang lalu) - harus ditolak
	old := now.Add(-50 * time.Hour)
	err = shared.ValidasiBatasWaktuRekamMedis(old.Format("2006-01-02"), old.Format("15:04:05"), 48, "dihapus")
	if err == nil {
		t.Fatal("expected ForbiddenError for record older than 48 hours, got nil")
	}
	if _, ok := err.(*apperror.ForbiddenError); !ok {
		t.Fatalf("expected *apperror.ForbiddenError, got %T", err)
	}

	// 4. Data lama tanpa jam (misal 5 hari lalu) - harus ditolak
	oldDateOnly := now.AddDate(0, 0, -5)
	err = shared.ValidasiBatasWaktuRekamMedis(oldDateOnly.Format("2006-01-02"), "", 48, "diubah")
	if err == nil {
		t.Fatal("expected ForbiddenError for date-only 5 days ago, got nil")
	}

	// 5. maxJam <= 0 (unlimited) - harus selalu lolos
	err = shared.ValidasiBatasWaktuRekamMedis("2020-01-01", "00:00:00", 0, "dihapus")
	if err != nil {
		t.Fatalf("expected nil error when maxJam is 0, got: %v", err)
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
