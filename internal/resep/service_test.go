package resep_test

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/resep"
	"erm-dokter/internal/shared"
)

type mockRepository struct {
	daftarResepFunc       func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error)
	daftarResepByRMFunc   func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error)
	daftarAturanPakaiFunc func(ctx context.Context, keyword string) ([]resep.AturanPakai, error)
	daftarMetodeRacikFunc func(ctx context.Context) ([]resep.MetodeRacik, error)
}

func (m *mockRepository) DaftarResep(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error) {
	if m.daftarResepFunc != nil {
		return m.daftarResepFunc(ctx, noRawat, statusLanjut, filter)
	}
	return []resep.Resep{}, 0, nil
}

func (m *mockRepository) DaftarResepByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error) {
	if m.daftarResepByRMFunc != nil {
		return m.daftarResepByRMFunc(ctx, noRM, statusLanjut, filter)
	}
	return []resep.Resep{}, 0, nil
}

func (m *mockRepository) DaftarAturanPakai(ctx context.Context, keyword string) ([]resep.AturanPakai, error) {
	if m.daftarAturanPakaiFunc != nil {
		return m.daftarAturanPakaiFunc(ctx, keyword)
	}
	return []resep.AturanPakai{}, nil
}

func (m *mockRepository) DaftarMetodeRacik(ctx context.Context) ([]resep.MetodeRacik, error) {
	if m.daftarMetodeRacikFunc != nil {
		return m.daftarMetodeRacikFunc(ctx)
	}
	return []resep.MetodeRacik{}, nil
}

func TestFilterDaftarResep_Validate(t *testing.T) {
	t.Run("Valid tanggal range", func(t *testing.T) {
		f := resep.FilterDaftarResep{
			Tanggal: "2026-04-01,2026-04-30",
			Page:    1,
			Limit:   10,
		}
		if errs := f.Validate(); errs != nil {
			t.Fatalf("expected nil validation errors, got %v", errs)
		}
	})

	t.Run("Invalid tanggal format", func(t *testing.T) {
		f := resep.FilterDaftarResep{
			Tanggal: "invalid-date",
		}
		if errs := f.Validate(); errs == nil {
			t.Fatal("expected validation error for invalid date format, got nil")
		}
	})

	t.Run("Invalid tanggal reversed range", func(t *testing.T) {
		f := resep.FilterDaftarResep{
			Tanggal: "2026-04-30,2026-04-01",
		}
		if errs := f.Validate(); errs == nil {
			t.Fatal("expected validation error for reversed date range, got nil")
		}
	})

	t.Run("Sanitize limit and page", func(t *testing.T) {
		f := resep.FilterDaftarResep{
			Page:  0,
			Limit: 200,
		}
		f.Sanitize()
		if f.Page != 1 {
			t.Errorf("expected page 1, got %d", f.Page)
		}
		if f.Limit != 100 {
			t.Errorf("expected limit 100, got %d", f.Limit)
		}
	})
}

func TestDaftarResep_Success(t *testing.T) {
	mockData := []resep.Resep{
		{
			NoResep:          "202604230001",
			NoRawat:          "2026/04/23/000001",
			TanggalPeresepan: "2026-04-23",
			JamPeresepan:     "10:00:00",
			Status:           "ralan",
			KodeDokter:       "DK001",
			NamaDokter:       "dr. Budi",
			ResepDokter: []resep.ResepDokter{
				{
					KodeObat:    "B001",
					NamaObat:    "Paracetamol 500mg",
					Jumlah:      10,
					Satuan:      "TAB",
					AturanPakai: "3x1",
				},
			},
			ResepDokterRacikan: []resep.ResepDokterRacikan{},
		},
	}

	mockRepo := &mockRepository{
		daftarResepFunc: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error) {
			return mockData, 1, nil
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	t.Run("Status Ralan", func(t *testing.T) {
		data, meta, err := svc.DaftarResep(ctx, "2026/04/23/000001", shared.StatusLanjutRawatJalan, resep.FilterDaftarResep{Page: 1, Limit: 5})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(data) != 1 {
			t.Fatalf("expected 1 resep, got %d", len(data))
		}
		if meta.TotalRecords != 1 {
			t.Fatalf("expected total records 1, got %d", meta.TotalRecords)
		}
	})

	t.Run("Status Semua", func(t *testing.T) {
		data, meta, err := svc.DaftarResep(ctx, "2026/04/23/000001", "Semua", resep.FilterDaftarResep{Page: 1, Limit: 5})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(data) != 1 {
			t.Fatalf("expected 1 resep, got %d", len(data))
		}
		if meta.TotalPages != 1 {
			t.Fatalf("expected total pages 1, got %d", meta.TotalPages)
		}
	})
}

func TestDaftarResep_RepoError(t *testing.T) {
	mockRepo := &mockRepository{
		daftarResepFunc: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	_, _, err := svc.DaftarResep(ctx, "2026/04/23/000001", shared.StatusLanjutRawatJalan, resep.FilterDaftarResep{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDaftarResepByRM_Success(t *testing.T) {
	mockData := []resep.Resep{
		{
			NoResep:          "202604230002",
			NoRawat:          "2026/04/23/000002",
			TanggalPeresepan: "2026-04-23",
			JamPeresepan:     "11:00:00",
			Status:           "ranap",
			KodeDokter:       "DK002",
			NamaDokter:       "dr. Siti",
			ResepDokter:      []resep.ResepDokter{},
			ResepDokterRacikan: []resep.ResepDokterRacikan{
				{
					NoRacik:         "1",
					NamaRacik:       "Puyer Flu",
					KodeMetodeRacik: "PULV",
					NamaMetodeRacik: "Puyer",
					JumlahRacikan:   10,
					AturanPakai:     "3x1 bungkus",
					Keterangan:      "Sesudah makan",
					DetailRacikan: []resep.ResepDokterRacikanDetail{
						{
							KodeObat:  "B002",
							NamaObat:  "CTM",
							Kandungan: "4mg",
							Jumlah:    2.5,
							Satuan:    "TAB",
						},
					},
				},
			},
		},
	}

	mockRepo := &mockRepository{
		daftarResepByRMFunc: func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error) {
			return mockData, 1, nil
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	t.Run("Status Ranap", func(t *testing.T) {
		data, meta, err := svc.DaftarResepByRM(ctx, "123456", shared.StatusLanjutRawatInap, resep.FilterDaftarResep{Page: 1, Limit: 5})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(data) != 1 {
			t.Fatalf("expected 1 resep, got %d", len(data))
		}
		if len(data[0].ResepDokterRacikan) != 1 {
			t.Fatalf("expected 1 racikan, got %d", len(data[0].ResepDokterRacikan))
		}
		if len(data[0].ResepDokterRacikan[0].DetailRacikan) != 1 {
			t.Fatalf("expected 1 detail racikan, got %d", len(data[0].ResepDokterRacikan[0].DetailRacikan))
		}
		if meta.TotalRecords != 1 {
			t.Fatalf("expected total records 1, got %d", meta.TotalRecords)
		}
	})

	t.Run("Status Semua", func(t *testing.T) {
		data, meta, err := svc.DaftarResepByRM(ctx, "123456", "Semua", resep.FilterDaftarResep{Page: 1, Limit: 5})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(data) != 1 {
			t.Fatalf("expected 1 resep, got %d", len(data))
		}
		if meta.TotalRecords != 1 {
			t.Fatalf("expected total records 1, got %d", meta.TotalRecords)
		}
	})
}

func TestDaftarResepByRM_RepoError(t *testing.T) {
	mockRepo := &mockRepository{
		daftarResepByRMFunc: func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter resep.FilterDaftarResep) ([]resep.Resep, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	_, _, err := svc.DaftarResepByRM(ctx, "123456", shared.StatusLanjutRawatJalan, resep.FilterDaftarResep{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDaftarAturanPakai_Success(t *testing.T) {
	mockData := []resep.AturanPakai{
		{AturanPakai: "3 x 1 tablet sesudah makan"},
		{AturanPakai: "2 x 1 tablet sesudah makan"},
	}

	mockRepo := &mockRepository{
		daftarAturanPakaiFunc: func(ctx context.Context, keyword string) ([]resep.AturanPakai, error) {
			if keyword == "3 x 1" {
				return []resep.AturanPakai{mockData[0]}, nil
			}
			return mockData, nil
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	t.Run("Tanpa keyword", func(t *testing.T) {
		data, err := svc.DaftarAturanPakai(ctx, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(data) != 2 {
			t.Fatalf("expected 2 aturan pakai, got %d", len(data))
		}
	})

	t.Run("Dengan keyword", func(t *testing.T) {
		data, err := svc.DaftarAturanPakai(ctx, "3 x 1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(data) != 1 {
			t.Fatalf("expected 1 aturan pakai, got %d", len(data))
		}
		if data[0].AturanPakai != "3 x 1 tablet sesudah makan" {
			t.Errorf("expected '3 x 1 tablet sesudah makan', got '%s'", data[0].AturanPakai)
		}
	})
}

func TestDaftarAturanPakai_RepoError(t *testing.T) {
	mockRepo := &mockRepository{
		daftarAturanPakaiFunc: func(ctx context.Context, keyword string) ([]resep.AturanPakai, error) {
			return nil, errors.New("database connection failed")
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	_, err := svc.DaftarAturanPakai(ctx, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDaftarMetodeRacik_Success(t *testing.T) {
	mockData := []resep.MetodeRacik{
		{Kode: "PULV", Nama: "Puyer / Serbuk"},
		{Kode: "CAPS", Nama: "Kapsul"},
	}

	mockRepo := &mockRepository{
		daftarMetodeRacikFunc: func(ctx context.Context) ([]resep.MetodeRacik, error) {
			return mockData, nil
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	data, err := svc.DaftarMetodeRacik(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 metode racik, got %d", len(data))
	}
	if data[0].Kode != "PULV" || data[0].Nama != "Puyer / Serbuk" {
		t.Errorf("expected PULV, got %+v", data[0])
	}
}

func TestDaftarMetodeRacik_RepoError(t *testing.T) {
	mockRepo := &mockRepository{
		daftarMetodeRacikFunc: func(ctx context.Context) ([]resep.MetodeRacik, error) {
			return nil, errors.New("database error")
		},
	}

	log := logger.New()
	svc := resep.NewService(mockRepo, log)
	ctx := context.Background()

	_, err := svc.DaftarMetodeRacik(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
