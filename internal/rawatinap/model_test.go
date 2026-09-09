package rawatinap_test

import (
	"testing"

	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/shared"
)

func TestFilterPasienRawatInap_Sanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    rawatinap.FilterPasienRawatInap
		expected rawatinap.FilterPasienRawatInap
	}{
		{
			name:  "Default values when empty",
			input: rawatinap.FilterPasienRawatInap{},
			expected: rawatinap.FilterPasienRawatInap{
				Page:            1,
				Limit:           20,
				StatusKunjungan: rawatinap.StatusKunjunganBelumKRS,
				ScopeDPJP:       rawatinap.ScopeDPJPSemua,
				OrderBy:         rawatinap.OrderByWaktuMRS,
				SortOrder:       "DESC",
			},
		},
		{
			name: "Max limit capped at 100",
			input: rawatinap.FilterPasienRawatInap{
				Page:      2,
				Limit:     250,
				SortOrder: "asc",
			},
			expected: rawatinap.FilterPasienRawatInap{
				Page:            2,
				Limit:           100,
				StatusKunjungan: rawatinap.StatusKunjunganBelumKRS,
				ScopeDPJP:       rawatinap.ScopeDPJPSemua,
				OrderBy:         rawatinap.OrderByWaktuMRS,
				SortOrder:       "ASC",
			},
		},
		{
			name: "Trim whitespace on filters",
			input: rawatinap.FilterPasienRawatInap{
				Bangsal:  "  B01  ",
				Kelas:    "  Kelas 1  ",
				Penjamin: "  BPJ  ",
				Keyword:  "  Pasien Budi  ",
				Tanggal:  "  2026-09-01,2026-09-08  ",
			},
			expected: rawatinap.FilterPasienRawatInap{
				Page:            1,
				Limit:           20,
				StatusKunjungan: rawatinap.StatusKunjunganBelumKRS,
				ScopeDPJP:       rawatinap.ScopeDPJPSemua,
				OrderBy:         rawatinap.OrderByWaktuMRS,
				SortOrder:       "DESC",
				Bangsal:         "B01",
				Kelas:           "Kelas 1",
				Penjamin:        "BPJ",
				Keyword:         "Pasien Budi",
				Tanggal:         "2026-09-01,2026-09-08",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.input
			f.Sanitize()

			if f.Page != tt.expected.Page {
				t.Errorf("Page got %d, want %d", f.Page, tt.expected.Page)
			}
			if f.Limit != tt.expected.Limit {
				t.Errorf("Limit got %d, want %d", f.Limit, tt.expected.Limit)
			}
			if f.StatusKunjungan != tt.expected.StatusKunjungan {
				t.Errorf("StatusKunjungan got %s, want %s", f.StatusKunjungan, tt.expected.StatusKunjungan)
			}
			if f.ScopeDPJP != tt.expected.ScopeDPJP {
				t.Errorf("ScopeDPJP got %s, want %s", f.ScopeDPJP, tt.expected.ScopeDPJP)
			}
			if f.OrderBy != tt.expected.OrderBy {
				t.Errorf("OrderBy got %s, want %s", f.OrderBy, tt.expected.OrderBy)
			}
			if f.SortOrder != tt.expected.SortOrder {
				t.Errorf("SortOrder got %s, want %s", f.SortOrder, tt.expected.SortOrder)
			}
			if f.Bangsal != tt.expected.Bangsal {
				t.Errorf("Bangsal got %q, want %q", f.Bangsal, tt.expected.Bangsal)
			}
			if f.Kelas != tt.expected.Kelas {
				t.Errorf("Kelas got %q, want %q", f.Kelas, tt.expected.Kelas)
			}
			if f.Penjamin != tt.expected.Penjamin {
				t.Errorf("Penjamin got %q, want %q", f.Penjamin, tt.expected.Penjamin)
			}
			if f.Keyword != tt.expected.Keyword {
				t.Errorf("Keyword got %q, want %q", f.Keyword, tt.expected.Keyword)
			}
			if f.Tanggal != tt.expected.Tanggal {
				t.Errorf("Tanggal got %q, want %q", f.Tanggal, tt.expected.Tanggal)
			}
		})
	}
}

func TestFilterPasienRawatInap_Validate(t *testing.T) {
	t.Run("Valid filter", func(t *testing.T) {
		f := rawatinap.FilterPasienRawatInap{
			SortOrder:       "ASC",
			OrderBy:         rawatinap.OrderByNamaPasien,
			StatusKunjungan: rawatinap.StatusKunjunganTanggalMRS,
			StatusPulang:    rawatinap.StatusPulangMembaik,
			ScopeDPJP:       rawatinap.ScopeDPJPDokter,
			Tanggal:         "2026-09-01,2026-09-08",
			Keyword:         "Siti",
		}
		errs := f.Validate()
		if errs != nil {
			t.Fatalf("expected nil errs, got %v", errs)
		}
	})

	t.Run("Invalid enums, status pulang, and keyword too short", func(t *testing.T) {
		f := rawatinap.FilterPasienRawatInap{
			SortOrder:       "RANDOM",
			OrderBy:         "invalid_order",
			StatusKunjungan: "invalid_status",
			StatusPulang:    "invalid_pulang",
			ScopeDPJP:       "invalid_scope",
			Keyword:         "ab",
			Tanggal:         "bukan-tanggal",
		}
		errs := f.Validate()
		if errs == nil {
			t.Fatal("expected validation errors, got nil")
		}
		if _, ok := errs["sort_order"]; !ok {
			t.Error("expected error for sort_order")
		}
		if _, ok := errs["order_by"]; !ok {
			t.Error("expected error for order_by")
		}
		if _, ok := errs["status_kunjungan"]; !ok {
			t.Error("expected error for status_kunjungan")
		}
		if _, ok := errs["status_pulang"]; !ok {
			t.Error("expected error for status_pulang")
		}
		if _, ok := errs["scope_dpjp"]; !ok {
			t.Error("expected error for scope_dpjp")
		}
		if _, ok := errs["keyword"]; !ok {
			t.Error("expected error for keyword")
		}
		if _, ok := errs["tanggal"]; !ok {
			t.Error("expected error for tanggal")
		}
	})
}

func TestStatusPulang_Values(t *testing.T) {
	expectedValues := []string{
		"Sehat", "Rujuk", "APS", "+", "Meninggal", "Sembuh", "Membaik",
		"Pulang Paksa", "-", "Pindah Kamar", "Atas Persetujuan Dokter",
		"Atas Permintaan Sendiri", "Isoman", "Lain-lain",
	}

	for _, v := range expectedValues {
		sp := rawatinap.StatusPulang(v)
		if !sp.IsValid() {
			t.Errorf("StatusPulang %q should be valid", v)
		}
	}

	if rawatinap.StatusPulang("Sembarang").IsValid() {
		t.Error("invalid status pulang should not be valid")
	}
}

func TestStatusKunjunganRanap_Values(t *testing.T) {
	validStatuses := []rawatinap.StatusKunjunganRanap{
		rawatinap.StatusKunjunganBelumKRS,
		rawatinap.StatusKunjunganTanggalMRS,
		rawatinap.StatusKunjunganTanggalKRS,
	}

	for _, s := range validStatuses {
		if !s.IsValid() {
			t.Errorf("StatusKunjunganRanap %q should be valid", s)
		}
	}

	invalidStatuses := []rawatinap.StatusKunjunganRanap{
		"belum_keluar",
		"belum_pulang",
		"tanggal_masuk",
		"tanggal_keluar",
		"tanggal_pulang",
		"semua",
		"invalid",
	}

	for _, s := range invalidStatuses {
		if s.IsValid() {
			t.Errorf("StatusKunjunganRanap %q should be invalid", s)
		}
	}
}

func TestScopeDPJP_Values(t *testing.T) {
	validScopes := []rawatinap.ScopeDPJP{
		rawatinap.ScopeDPJPDokter,
		rawatinap.ScopeDPJPSemua,
	}

	for _, s := range validScopes {
		if !s.IsValid() {
			t.Errorf("ScopeDPJP %q should be valid", s)
		}
	}

	invalidScopes := []rawatinap.ScopeDPJP{
		"dpjp_saya",
		"tanpa_dpjp",
		"invalid",
	}

	for _, s := range invalidScopes {
		if s.IsValid() {
			t.Errorf("ScopeDPJP %q should be invalid", s)
		}
	}
}

func TestOrderByRanap_Values(t *testing.T) {
	validOrderBys := []rawatinap.OrderByRanap{
		rawatinap.OrderByWaktuMRS,
		rawatinap.OrderByWaktuKRS,
		rawatinap.OrderByNamaPasien,
		rawatinap.OrderByKamar,
		rawatinap.OrderByKelas,
		rawatinap.OrderByBangsal,
		rawatinap.OrderByPenjamin,
		rawatinap.OrderByStatusPulang,
	}

	for _, ob := range validOrderBys {
		if !ob.IsValid() {
			t.Errorf("OrderByRanap %q should be valid", ob)
		}
	}

	invalidOrderBys := []rawatinap.OrderByRanap{
		"waktu_masuk",
		"waktu_keluar",
		"waktu_pulang",
		"status_keluar",
		"no_rawat",
		"invalid",
	}

	for _, ob := range invalidOrderBys {
		if ob.IsValid() {
			t.Errorf("OrderByRanap %q should be invalid", ob)
		}
	}
}

func TestFilterPasienRawatInap_Offset(t *testing.T) {
	f := rawatinap.FilterPasienRawatInap{
		Page:  3,
		Limit: 15,
	}
	if got := f.Offset(); got != 30 {
		t.Errorf("Offset() = %d, want 30", got)
	}
}

func TestKunjunganRawatInap_Formatters(t *testing.T) {
	k := rawatinap.KunjunganRawatInap{
		NoRekamMedis: "00123456",
		JenisKelamin: "L",
		TanggalLahir: "2000-01-01",
		StatusLanjut: shared.StatusLanjutRawatInap,
	}

	if k.FormatNoRekamMedis() != "123456" {
		t.Errorf("FormatNoRekamMedis() = %s, want 123456", k.FormatNoRekamMedis())
	}
	if k.FormatJenisKelamin() != "Laki-Laki" {
		t.Errorf("FormatJenisKelamin() = %s, want Laki-Laki", k.FormatJenisKelamin())
	}
	if k.FormatUmur() == "" {
		t.Error("FormatUmur() should not be empty")
	}
}
