package shared_test

import (
	"testing"

	"erm-dokter/internal/shared"
)

func TestStatusLanjut_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   shared.StatusLanjut
		expected bool
	}{
		{name: "Valid Ralan", status: shared.StatusLanjutRawatJalan, expected: true},
		{name: "Valid Ranap", status: shared.StatusLanjutRawatInap, expected: true},
		{name: "Invalid Empty", status: "", expected: false},
		{name: "Invalid Random", status: "IGD", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.expected {
				t.Errorf("StatusLanjut.IsValid() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestParseStatusLanjut(t *testing.T) {
	tests := []struct {
		input       string
		expected    shared.StatusLanjut
		expectedOk  bool
	}{
		{"ralan", shared.StatusLanjutRawatJalan, true},
		{"Ralan", shared.StatusLanjutRawatJalan, true},
		{"RALAN", shared.StatusLanjutRawatJalan, true},
		{"ranap", shared.StatusLanjutRawatInap, true},
		{"Ranap", shared.StatusLanjutRawatInap, true},
		{"RANAP", shared.StatusLanjutRawatInap, true},
		{"semua", "", false},
		{"igd", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got, ok := shared.ParseStatusLanjut(tt.input)
		if ok != tt.expectedOk || got != tt.expected {
			t.Errorf("ParseStatusLanjut(%q) = (%v, %v), expected (%v, %v)", tt.input, got, ok, tt.expected, tt.expectedOk)
		}
	}
}

func TestParseStatusLanjutWithSemua(t *testing.T) {
	tests := []struct {
		input      string
		expected   shared.StatusLanjut
		expectedOk bool
	}{
		{"ralan", shared.StatusLanjutRawatJalan, true},
		{"Ralan", shared.StatusLanjutRawatJalan, true},
		{"ranap", shared.StatusLanjutRawatInap, true},
		{"Ranap", shared.StatusLanjutRawatInap, true},
		{"semua", "Semua", true},
		{"Semua", "Semua", true},
		{"SEMUA", "Semua", true},
		{"invalid", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got, ok := shared.ParseStatusLanjutWithSemua(tt.input)
		if ok != tt.expectedOk || got != tt.expected {
			t.Errorf("ParseStatusLanjutWithSemua(%q) = (%v, %v), expected (%v, %v)", tt.input, got, ok, tt.expected, tt.expectedOk)
		}
	}
}

func TestSortOrder_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		sort     shared.SortOrder
		expected bool
	}{
		{name: "Valid ASC", sort: shared.SortASC, expected: true},
		{name: "Valid DESC", sort: shared.SortDESC, expected: true},
		{name: "Invalid Empty", sort: "", expected: false},
		{name: "Invalid Random", sort: "NONE", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sort.IsValid(); got != tt.expected {
				t.Errorf("SortOrder.IsValid() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestKategoriLab_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		kategori shared.KategoriLab
		expected bool
	}{
		{name: "Valid PK", kategori: shared.KategoriLabPK, expected: true},
		{name: "Valid PA", kategori: shared.KategoriLabPA, expected: true},
		{name: "Valid MB", kategori: shared.KategoriLabMB, expected: true},
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
