package tindakan_test

import (
	"testing"

	"erm-dokter/internal/tindakan"
)

func TestKategoriLab_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		kategori tindakan.KategoriLab
		expected bool
	}{
		{name: "Valid PK", kategori: tindakan.KategoriLabPK, expected: true},
		{name: "Valid PA", kategori: tindakan.KategoriLabPA, expected: true},
		{name: "Valid MB", kategori: tindakan.KategoriLabMB, expected: true},
		{name: "Invalid Empty", kategori: "", expected: false},
		{name: "Invalid Random", kategori: "XYZ", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kategori.IsValid(); got != tt.expected {
				t.Errorf("KategoriLab.IsValid() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestFilterDaftarTindakanLab_Sanitize(t *testing.T) {
	f := tindakan.FilterDaftarTindakanLab{
		Keyword: "   Darah Lengkap   ",
		Page:    -1,
		Limit:   0,
	}
	f.Sanitize()

	if f.Page != 1 {
		t.Errorf("Expected page to default to 1, got %d", f.Page)
	}
	if f.Limit != 20 {
		t.Errorf("Expected limit to default to 20, got %d", f.Limit)
	}
	if f.Keyword != "Darah Lengkap" {
		t.Errorf("Expected trimmed keyword, got %q", f.Keyword)
	}

	fOver := tindakan.FilterDaftarTindakanLab{
		Page:  2,
		Limit: 200,
	}
	fOver.Sanitize()
	if fOver.Limit != 100 {
		t.Errorf("Expected max limit 100, got %d", fOver.Limit)
	}
}

func TestFilterDaftarTindakanLab_Offset(t *testing.T) {
	f := tindakan.FilterDaftarTindakanLab{Page: 1, Limit: 20}
	if f.Offset() != 0 {
		t.Errorf("Expected offset 0 for page 1, got %d", f.Offset())
	}

	f2 := tindakan.FilterDaftarTindakanLab{Page: 3, Limit: 25}
	if f2.Offset() != 50 {
		t.Errorf("Expected offset 50 for page 3 with limit 25, got %d", f2.Offset())
	}
}

func TestFilterDaftarTindakanLab_Validate(t *testing.T) {
	tests := []struct {
		name        string
		filter      tindakan.FilterDaftarTindakanLab
		expectError bool
		errorField  string
	}{
		{
			name: "Valid Filter without keyword",
			filter: tindakan.FilterDaftarTindakanLab{
				Page:  1,
				Limit: 20,
			},
			expectError: false,
		},
		{
			name: "Valid Filter with keyword >= 3 chars",
			filter: tindakan.FilterDaftarTindakanLab{
				Keyword: "Darah",
				Page:    1,
				Limit:   20,
			},
			expectError: false,
		},
		{
			name: "Invalid Keyword < 3 chars",
			filter: tindakan.FilterDaftarTindakanLab{
				Keyword: "Da",
				Page:    1,
				Limit:   20,
			},
			expectError: true,
			errorField:  "keyword",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.filter.Validate()
			if tt.expectError {
				if errs == nil {
					t.Fatalf("Expected validation error, got nil")
				}
				if _, ok := errs[tt.errorField]; !ok {
					t.Errorf("Expected error in field %s, got errors: %v", tt.errorField, errs)
				}
				return
			}
			if errs != nil {
				t.Fatalf("Expected no validation errors, got: %v", errs)
			}
		})
	}
}
