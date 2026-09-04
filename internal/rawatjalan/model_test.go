package rawatjalan_test

import (
	"testing"
	"time"

	"erm-dokter/internal/rawatjalan"
)

func TestFilterAntreanDokter_Validation(t *testing.T) {
	t.Run("Valid filter", func(t *testing.T) {
		f := rawatjalan.FilterAntreanDokter{
			Tanggal: "2026-09-03,2026-09-03",
			Page:    1,
			Limit:   20,
		}
		f.Sanitize()
		if errs := f.Validate(); errs != nil {
			t.Errorf("Expected valid filter, got: %+v", errs)
		}
	})

	t.Run("Sanitize defaults", func(t *testing.T) {
		var f rawatjalan.FilterAntreanDokter
		f.Sanitize()
		if f.Page != 1 {
			t.Errorf("Expected Page=1, got %d", f.Page)
		}
		if f.Limit != 20 {
			t.Errorf("Expected Limit=20, got %d", f.Limit)
		}
		if f.OrderBy != "waktu_registrasi" {
			t.Errorf("Expected OrderBy=waktu_registrasi, got %s", f.OrderBy)
		}
		if f.SortOrder != "ASC" {
			t.Errorf("Expected SortOrder=ASC, got %s", f.SortOrder)
		}
		if f.Tanggal == "" {
			t.Error("Expected non-empty default Tanggal")
		}
	})

	t.Run("Invalid OrderBy and SortOrder", func(t *testing.T) {
		f := rawatjalan.FilterAntreanDokter{
			OrderBy:   "malicious",
			SortOrder: "HACK",
		}
		errs := f.Validate()
		if errs == nil {
			t.Fatal("Expected validation error")
		}
		if _, exists := errs["order_by"]; !exists {
			t.Error("Expected order_by error")
		}
		if _, exists := errs["sort_order"]; !exists {
			t.Error("Expected sort_order error")
		}
	})
}

func TestKunjunganRawatJalan_Formatters(t *testing.T) {
	birthDate := time.Now().AddDate(-40, -1, -5).Format("2006-01-02")
	k := rawatjalan.KunjunganRawatJalan{
		NoRawat:      "2026/09/03/000001",
		NoRekamMedis: "00998877",
		TanggalLahir: birthDate,
		JenisKelamin: "L",
	}

	if k.FormatNoRekamMedis() != "998877" {
		t.Errorf("Expected 998877, got %s", k.FormatNoRekamMedis())
	}

	if k.FormatJenisKelamin() != "Laki-Laki" {
		t.Errorf("Expected Laki-Laki, got %s", k.FormatJenisKelamin())
	}

	if k.FormatUmur() == "" {
		t.Error("Expected non-empty format umur")
	}
}
