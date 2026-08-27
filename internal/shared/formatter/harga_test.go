package formatter_test

import (
	"testing"

	"erm-dokter/internal/shared/formatter"
)

func TestFormatHarga(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{0, "Rp 0,00"},
		{0.0, "Rp 0,00"},
		{500, "Rp 500,00"},
		{1500, "Rp 1.500,00"},
		{15000, "Rp 15.000,00"},
		{1250000, "Rp 1.250.000,00"},
		{1500.50, "Rp 1.500,50"},
		{1500.75, "Rp 1.500,75"},
		{"1500.00", "Rp 1.500,00"},
		{"15000", "Rp 15.000,00"},
		{-5000, "-Rp 5.000,00"},
	}

	for _, tc := range tests {
		got := formatter.FormatHarga(tc.input)
		if got != tc.expected {
			t.Errorf("FormatHarga(%v) = %s; expected %s", tc.input, got, tc.expected)
		}
	}
}
