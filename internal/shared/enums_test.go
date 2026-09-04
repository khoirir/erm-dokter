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
