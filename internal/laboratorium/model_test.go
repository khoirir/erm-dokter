package laboratorium_test

import (
	"testing"
	"time"

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

func TestSimpanPermintaanLabPKRequest_ValidationAndSanitization(t *testing.T) {
	nowStr := time.Now().Format("2006-01-02")
	jamNowStr := time.Now().Format("15:04:05")

	t.Run("Valid request with full fields", func(t *testing.T) {
		req := laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Febris H-3",
			InformasiTambahan: "Cek rutin",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{
					IdTindakan: "enc-tindakan-1",
					IdTemplate: []string{"enc-tmpl-1", "enc-tmpl-2"},
				},
			},
		}

		errs := req.Validate()
		if errs != nil {
			t.Fatalf("Expected valid, got errors: %v", errs)
		}
		if req.InformasiTambahan != "Cek rutin" {
			t.Errorf("Expected informasi tambahan tetap, got %s", req.InformasiTambahan)
		}
	})

	t.Run("Valid request with HH:MM format and empty informasi tambahan defaults to hyphen", func(t *testing.T) {
		req := laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "2026/09/03/000001",
			TanggalPermintaan: nowStr,
			JamPermintaan:     "10:30",
			DiagnosaKlinis:    "Susp TBC",
			InformasiTambahan: "",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{
					IdTindakan: "enc-tindakan-2",
				},
			},
		}

		errs := req.Validate()
		if errs != nil {
			t.Fatalf("Expected valid, got errors: %v", errs)
		}
		if req.JamPermintaan != "10:30:00" {
			t.Errorf("Expected jam 10:30:00, got %s", req.JamPermintaan)
		}
		if req.InformasiTambahan != "-" {
			t.Errorf("Expected default '-', got %s", req.InformasiTambahan)
		}
	})

	t.Run("Invalid empty required fields", func(t *testing.T) {
		req := laboratorium.SimpanPermintaanLabPKRequest{
			NoRawat:           "",
			TanggalPermintaan: "",
			JamPermintaan:     "",
			DiagnosaKlinis:    "",
			Pemeriksaan:       []laboratorium.ItemPemeriksaanLabPKRequest{},
		}

		errs := req.Validate()
		if errs == nil {
			t.Fatal("Expected validation errors, got nil")
		}
		if _, ok := errs["no_rawat"]; !ok {
			t.Error("Expected error on no_rawat")
		}
		if _, ok := errs["tanggal_permintaan"]; !ok {
			t.Error("Expected error on tanggal_permintaan")
		}
		if _, ok := errs["jam_permintaan"]; !ok {
			t.Error("Expected error on jam_permintaan")
		}
		if _, ok := errs["diagnosa_klinis"]; !ok {
			t.Error("Expected error on diagnosa_klinis")
		}
		if _, ok := errs["pemeriksaan"]; !ok {
			t.Error("Expected error on pemeriksaan")
		}
	})

	t.Run("Invalid date and time formats", func(t *testing.T) {
		req := laboratorium.SimpanPermintaanLabPKRequest{
			TanggalPermintaan: "03-09-2026",
			JamPermintaan:     "jam-10",
			DiagnosaKlinis:    "Anemia",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{
					IdTindakan: "",
				},
			},
		}

		errs := req.Validate()
		if errs == nil {
			t.Fatal("Expected validation errors, got nil")
		}
		if _, ok := errs["tanggal_permintaan"]; !ok {
			t.Error("Expected error on format tanggal_permintaan")
		}
		if _, ok := errs["jam_permintaan"]; !ok {
			t.Error("Expected error on format jam_permintaan")
		}
		if _, ok := errs["pemeriksaan[0].id_tindakan"]; !ok {
			t.Error("Expected error on empty id_tindakan")
		}
	})

	t.Run("Invalid future time", func(t *testing.T) {
		besok := time.Now().Add(24 * time.Hour).Format("2006-01-02")
		req := laboratorium.SimpanPermintaanLabPKRequest{
			TanggalPermintaan: besok,
			JamPermintaan:     "10:00:00",
			DiagnosaKlinis:    "Anemia",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{IdTindakan: "enc-tindakan-1"},
			},
		}

		errs := req.Validate()
		if errs == nil {
			t.Fatal("Expected validation errors for future time, got nil")
		}
		if _, ok := errs["tanggal_permintaan"]; !ok {
			t.Error("Expected error on future tanggal_permintaan")
		}
	})

	t.Run("Invalid empty id_template in slice", func(t *testing.T) {
		req := laboratorium.SimpanPermintaanLabPKRequest{
			TanggalPermintaan: nowStr,
			JamPermintaan:     jamNowStr,
			DiagnosaKlinis:    "Anemia",
			Pemeriksaan: []laboratorium.ItemPemeriksaanLabPKRequest{
				{
					IdTindakan: "enc-tindakan-1",
					IdTemplate: []string{""},
				},
			},
		}

		errs := req.Validate()
		if errs == nil {
			t.Fatal("Expected validation errors for empty id_template, got nil")
		}
		if _, ok := errs["pemeriksaan[0].id_template[0]"]; !ok {
			t.Error("Expected error on pemeriksaan[0].id_template[0]")
		}
	})
}

