package pasien_test

import (
	"testing"
	"time"

	"erm-dokter/internal/pasien"
)

func TestPasienModelFormatters(t *testing.T) {
	birthDate := time.Now().AddDate(-25, -2, -10).Format("2006-01-02")

	p := pasien.Pasien{
		NoRM:         "00123456",
		Nama:         "Budi Santoso",
		JenisKelamin: "L",
		TanggalLahir: birthDate,
	}

	if got := p.FormatNoRekamMedis(); got != "123456" {
		t.Errorf("Expected FormatNoRekamMedis '123456', got '%s'", got)
	}

	if got := p.FormatJenisKelamin(); got != "Laki-Laki" {
		t.Errorf("Expected FormatJenisKelamin 'Laki-Laki', got '%s'", got)
	}

	p2 := pasien.Pasien{
		NoRM:         "654321",
		Nama:         "Siti Rahma",
		JenisKelamin: "P",
		TanggalLahir: birthDate,
	}

	if got := p2.FormatNoRekamMedis(); got != "654321" {
		t.Errorf("Expected FormatNoRekamMedis '654321', got '%s'", got)
	}

	if got := p2.FormatJenisKelamin(); got != "Perempuan" {
		t.Errorf("Expected FormatJenisKelamin 'Perempuan', got '%s'", got)
	}

	umur := p.HitungUmurLengkap()
	if umur.Tahun != 25 {
		t.Errorf("Expected age in years to be 25, got %d", umur.Tahun)
	}

	umurStr := p.FormatUmur()
	if umurStr == "" {
		t.Errorf("Expected non-empty format umur")
	}
}
