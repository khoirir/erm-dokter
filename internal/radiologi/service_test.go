package radiologi_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/radiologi"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
	daftarHasilRadiologiKunjunganFn func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error)
	daftarHasilRadiologiPasienFn    func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error)
	detailHasilRadiologiFn          func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error)
}

func (m *mockRepository) DaftarHasilRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error) {
	if m.daftarHasilRadiologiKunjunganFn != nil {
		return m.daftarHasilRadiologiKunjunganFn(ctx, noRawat, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DaftarHasilRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error) {
	if m.daftarHasilRadiologiPasienFn != nil {
		return m.daftarHasilRadiologiPasienFn(ctx, noRM, statusLanjut, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) DetailHasilRadiologi(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error) {
	if m.detailHasilRadiologiFn != nil {
		return m.detailHasilRadiologiFn(ctx, idHasil, statusLanjut)
	}
	return nil, nil
}

func TestService_GetRiwayatRadiologiKunjungan_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilRadiologiKunjunganFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error) {
			if noRawat != "2026/04/22/000001" {
				t.Errorf("Expected noRawat 2026/04/22/000001, got %s", noRawat)
			}
			return []radiologi.HasilRadiologi{
				{
					NoRawat:        noRawat,
					KodeTindakan:   "RAD001",
					NamaTindakan:   "Rontgen Thorax AP/PA",
					Status:         "Ralan",
					TanggalPeriksa: "2026-04-22",
					JamPeriksa:     "10:00:00",
					Hasil:          "Cor dan Pulmo dalam batas normal",
					GambarPACS:     []string{"http://pacs.example.com/viewer?token=xyz"},
				},
			}, 1, nil
		},
	}

	svc := radiologi.NewService(mockRepo, log)
	list, meta, err := svc.GetRiwayatRadiologiKunjungan(context.Background(), "2026/04/22/000001", shared.StatusLanjutRawatJalan, radiologi.FilterRiwayatRadiologi{Page: 1, Limit: 5})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(list))
	}
	if list[0].KodeTindakan != "RAD001" {
		t.Errorf("Expected RAD001, got %s", list[0].KodeTindakan)
	}
	if len(list[0].GambarPACS) != 1 {
		t.Errorf("Expected 1 PACS image, got %d", len(list[0].GambarPACS))
	}
	if meta.TotalRecords != 1 {
		t.Errorf("Expected 1 total record, got %d", meta.TotalRecords)
	}
}

func TestService_GetRiwayatRadiologiKunjungan_Error(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilRadiologiKunjunganFn: func(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error) {
			return nil, 0, errors.New("database connection failed")
		},
	}

	svc := radiologi.NewService(mockRepo, log)
	_, _, err := svc.GetRiwayatRadiologiKunjungan(context.Background(), "2026/04/22/000001", shared.StatusLanjutRawatJalan, radiologi.FilterRiwayatRadiologi{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestService_GetRiwayatRadiologiPasien_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilRadiologiPasienFn: func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error) {
			if noRM != "123456" {
				t.Errorf("Expected noRM 123456, got %s", noRM)
			}
			return []radiologi.HasilRadiologi{
				{
					NoRawat:        "2026/04/22/000001",
					KodeTindakan:   "RAD001",
					NamaTindakan:   "Rontgen Thorax",
					TanggalPeriksa: "2026-04-22",
					JamPeriksa:     "10:00:00",
					Hasil:          "Infiltrat minimal di apex pulmo",
				},
			}, 1, nil
		},
	}

	svc := radiologi.NewService(mockRepo, log)
	list, meta, err := svc.GetRiwayatRadiologiPasien(context.Background(), "123456", "Semua", radiologi.FilterRiwayatRadiologi{Page: 1, Limit: 5})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(list))
	}
	if meta.TotalRecords != 1 {
		t.Errorf("Expected 1 total record, got %d", meta.TotalRecords)
	}
}

func TestService_GetRiwayatRadiologiPasien_Error(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		daftarHasilRadiologiPasienFn: func(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter radiologi.FilterRiwayatRadiologi) ([]radiologi.HasilRadiologi, int, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := radiologi.NewService(mockRepo, log)
	_, _, err := svc.GetRiwayatRadiologiPasien(context.Background(), "123456", "Semua", radiologi.FilterRiwayatRadiologi{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestService_GetDetailHasilRadiologi_Success(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		detailHasilRadiologiFn: func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error) {
			return &radiologi.HasilRadiologi{
				NoRawat:        idHasil.NoRawat,
				KodeTindakan:   idHasil.KodeTindakan,
				NamaTindakan:   "Rontgen Thorax",
				Status:         "Ralan",
				TanggalPeriksa: idHasil.TanggalPeriksa,
				JamPeriksa:     idHasil.JamPeriksa,
				Hasil:          "Cor dan Pulmo normal",
				GambarPACS:     []string{"http://pacs.example.com/viewer?token=abc"},
			}, nil
		},
	}

	svc := radiologi.NewService(mockRepo, log)
	idHasil := radiologi.IdHasilRadiologi{
		NoRawat:        "2026/04/22/000001",
		KodeTindakan:   "RAD001",
		TanggalPeriksa: "2026-04-22",
		JamPeriksa:     "10:00:00",
	}
	res, err := svc.GetDetailHasilRadiologi(context.Background(), idHasil, "Ralan")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res.Hasil != "Cor dan Pulmo normal" {
		t.Errorf("Expected hasil bacaan, got %s", res.Hasil)
	}
}

func TestService_GetDetailHasilRadiologi_NotFound(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		detailHasilRadiologiFn: func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error) {
			return nil, sql.ErrNoRows
		},
	}

	svc := radiologi.NewService(mockRepo, log)
	idHasil := radiologi.IdHasilRadiologi{
		NoRawat:        "2026/04/22/000001",
		KodeTindakan:   "RAD999",
		TanggalPeriksa: "2026-04-22",
		JamPeriksa:     "10:00:00",
	}
	_, err := svc.GetDetailHasilRadiologi(context.Background(), idHasil, "Semua")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	var notFoundErr *apperror.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("Expected NotFoundError, got %T: %v", err, err)
	}
}

func TestService_GetDetailHasilRadiologi_Error(t *testing.T) {
	log := logger.New()
	mockRepo := &mockRepository{
		detailHasilRadiologiFn: func(ctx context.Context, idHasil radiologi.IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*radiologi.HasilRadiologi, error) {
			return nil, errors.New("db connection failure")
		},
	}

	svc := radiologi.NewService(mockRepo, log)
	idHasil := radiologi.IdHasilRadiologi{
		NoRawat:        "2026/04/22/000001",
		KodeTindakan:   "RAD001",
		TanggalPeriksa: "2026-04-22",
		JamPeriksa:     "10:00:00",
	}
	_, err := svc.GetDetailHasilRadiologi(context.Background(), idHasil, "Semua")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
