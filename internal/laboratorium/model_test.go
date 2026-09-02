package laboratorium_test

import (
	"testing"

	"erm-dokter/internal/laboratorium"
)

func TestFilterRiwayatLab_Sanitize(t *testing.T) {
	filter := laboratorium.FilterRiwayatLab{
		Page:  -1,
		Limit: 0,
	}
	filter.Sanitize()

	if filter.Page != 1 {
		t.Errorf("Expected Page 1, got %d", filter.Page)
	}
	if filter.Limit != 5 {
		t.Errorf("Expected Limit 5, got %d", filter.Limit)
	}

	filterHigh := laboratorium.FilterRiwayatLab{
		Page:  2,
		Limit: 200,
	}
	filterHigh.Sanitize()
	if filterHigh.Limit != 100 {
		t.Errorf("Expected Limit 100, got %d", filterHigh.Limit)
	}
}

func TestFilterRiwayatLab_Offset(t *testing.T) {
	filter := laboratorium.FilterRiwayatLab{
		Page:  3,
		Limit: 10,
	}
	if filter.Offset() != 20 {
		t.Errorf("Expected Offset 20, got %d", filter.Offset())
	}
}

func TestFilterRiwayatLab_Validate(t *testing.T) {
	tests := []struct {
		name    string
		filter  laboratorium.FilterRiwayatLab
		wantErr bool
	}{
		{
			name: "Valid without date",
			filter: laboratorium.FilterRiwayatLab{
				Page:  1,
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "Valid with same date range",
			filter: laboratorium.FilterRiwayatLab{
				Tanggal: "2026-04-22,2026-04-22",
				Page:    1,
				Limit:   10,
			},
			wantErr: false,
		},
		{
			name: "Valid with date range",
			filter: laboratorium.FilterRiwayatLab{
				Tanggal: "2026-04-01,2026-04-30",
				Page:    1,
				Limit:   10,
			},
			wantErr: false,
		},
		{
			name: "Invalid single date without comma",
			filter: laboratorium.FilterRiwayatLab{
				Tanggal: "2026-04-22",
				Page:    1,
				Limit:   10,
			},
			wantErr: true,
		},
		{
			name: "Invalid date format",
			filter: laboratorium.FilterRiwayatLab{
				Tanggal: "invalid-date",
				Page:    1,
				Limit:   10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.filter.Validate()
			if (errs != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", errs, tt.wantErr)
			}
		})
	}
}
