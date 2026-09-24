package formatter_test

import (
	"testing"

	"erm-dokter/internal/shared/formatter"
)

func TestParseRentangTanggal(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantAwal  string
		wantAkhir string
	}{
		{
			name:      "Valid rentang tanggal",
			input:     "2026-09-01,2026-09-09",
			wantAwal:  "2026-09-01",
			wantAkhir: "2026-09-09",
		},
		{
			name:      "Valid rentang tanggal dengan spasi ekstra",
			input:     " 2026-09-01 , 2026-09-09 ",
			wantAwal:  "2026-09-01",
			wantAkhir: "2026-09-09",
		},
		{
			name:      "Bukan 2 bagian (tanpa koma)",
			input:     "2026-09-01",
			wantAwal:  "",
			wantAkhir: "",
		},
		{
			name:      "Lebih dari 2 bagian",
			input:     "2026-09-01,2026-09-05,2026-09-09",
			wantAwal:  "",
			wantAkhir: "",
		},
		{
			name:      "Tanggal awal kosong",
			input:     " , 2026-09-09",
			wantAwal:  "",
			wantAkhir: "",
		},
		{
			name:      "Tanggal akhir kosong",
			input:     "2026-09-01, ",
			wantAwal:  "",
			wantAkhir: "",
		},
		{
			name:      "String input kosong",
			input:     "",
			wantAwal:  "",
			wantAkhir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			awal, akhir := formatter.ParseRentangTanggal(tt.input)
			if awal != tt.wantAwal || akhir != tt.wantAkhir {
				t.Errorf("ParseRentangTanggal(%q) = (%q, %q), want (%q, %q)", tt.input, awal, akhir, tt.wantAwal, tt.wantAkhir)
			}
		})
	}
}
