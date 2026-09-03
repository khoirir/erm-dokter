package formatter_test

import (
	"testing"
	"time"

	"erm-dokter/internal/shared/formatter"
)

func TestFormatJenisKelamin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"L", "Laki-Laki"},
		{"l", "Laki-Laki"},
		{"P", "Perempuan"},
		{"p", "Perempuan"},
		{"-", "-"},
		{"Unknown", "Unknown"},
	}

	for _, tc := range tests {
		if got := formatter.FormatJenisKelamin(tc.input); got != tc.expected {
			t.Errorf("FormatJenisKelamin(%s) = %s, expected %s", tc.input, got, tc.expected)
		}
	}
}

func TestFormatNoRekamMedis(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"00123456", "123456"},
		{"00000001", "000001"},
		{"123456", "123456"},
		{"  00123456  ", "123456"},
		{"-", "-"},
	}

	for _, tc := range tests {
		if got := formatter.FormatNoRekamMedis(tc.input); got != tc.expected {
			t.Errorf("FormatNoRekamMedis(%s) = %s, expected %s", tc.input, got, tc.expected)
		}
	}
}

func TestHitungUmur(t *testing.T) {
	now := time.Now()
	birthDate := now.AddDate(-30, -5, -15).Format("2006-01-02")

	u := formatter.HitungUmur(birthDate)
	if u.Tahun != 30 {
		t.Errorf("Expected tahun 30, got %d", u.Tahun)
	}

	str := formatter.FormatUmur(birthDate)
	if str == "" {
		t.Errorf("Expected non-empty format umur")
	}

	invalid := formatter.HitungUmur("invalid-date")
	if invalid.Tahun != 0 || invalid.Bulan != 0 || invalid.Hari != 0 {
		t.Errorf("Expected empty umur for invalid date format")
	}
}
