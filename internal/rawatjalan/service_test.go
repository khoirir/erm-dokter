package rawatjalan_test

import (
	"context"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
	daftarAntreanFunc   func(ctx context.Context, kodeDokter string, filter rawatjalan.FilterAntreanDokter) ([]rawatjalan.KunjunganRawatJalan, int, error)
	detailKunjunganFunc func(ctx context.Context, noRawat string, kodeDokter string) (*rawatjalan.KunjunganRawatJalan, error)
}

func (m *mockRepository) DaftarAntreanDokter(ctx context.Context, kodeDokter string, filter rawatjalan.FilterAntreanDokter) ([]rawatjalan.KunjunganRawatJalan, int, error) {
	if m.daftarAntreanFunc != nil {
		return m.daftarAntreanFunc(ctx, kodeDokter, filter)
	}
	return []rawatjalan.KunjunganRawatJalan{}, 0, nil
}

func (m *mockRepository) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*rawatjalan.KunjunganRawatJalan, error) {
	if m.detailKunjunganFunc != nil {
		return m.detailKunjunganFunc(ctx, noRawat, kodeDokter)
	}
	return nil, nil
}

func TestDaftarAntreanDokter_ValidationError(t *testing.T) {
	repo := &mockRepository{}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)

	filter := rawatjalan.FilterAntreanDokter{
		OrderBy:   "kolom_salah",
		SortOrder: "SALAH",
	}

	_, _, err := uc.DaftarAntreanDokter(context.Background(), "DK001", filter)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	validationErr, ok := err.(apperror.ValidationError)
	if !ok {
		t.Fatalf("expected apperror.ValidationError, got %T", err)
	}

	if _, exists := validationErr["order_by"]; !exists {
		t.Error("expected validation error for 'order_by'")
	}
	if _, exists := validationErr["sort_order"]; !exists {
		t.Error("expected validation error for 'sort_order'")
	}
}

func TestDaftarAntreanDokter_MinKeywordLength(t *testing.T) {
	repo := &mockRepository{}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)

	filter := rawatjalan.FilterAntreanDokter{KataKunci: "ab"}

	_, _, err := uc.DaftarAntreanDokter(context.Background(), "DK001", filter)
	if err == nil {
		t.Fatal("expected business error for short keyword, got nil")
	}

	businessErr, ok := err.(*apperror.BusinessError)
	if !ok {
		t.Fatalf("expected *apperror.BusinessError, got %T", err)
	}
	if businessErr.Message != "kata kunci pencarian minimal 3 karakter" {
		t.Errorf("unexpected error message: %s", businessErr.Message)
	}
}

func TestDaftarAntreanDokter_Success(t *testing.T) {
	repo := &mockRepository{
		daftarAntreanFunc: func(ctx context.Context, kodeDokter string, filter rawatjalan.FilterAntreanDokter) ([]rawatjalan.KunjunganRawatJalan, int, error) {
			return []rawatjalan.KunjunganRawatJalan{
				{NoRawat: "2025/04/22/000001", NamaPasien: "Budi"},
				{NoRawat: "2025/04/22/000002", NamaPasien: "Ani"},
			}, 2, nil
		},
	}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)

	filter := rawatjalan.FilterAntreanDokter{}

	data, meta, err := uc.DaftarAntreanDokter(context.Background(), "DK001", filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 2 {
		t.Errorf("expected 2 results, got %d", len(data))
	}
	if meta.TotalData != 2 {
		t.Errorf("expected TotalData=2, got %d", meta.TotalData)
	}
}

func TestDaftarAntreanDokter_Pagination(t *testing.T) {
	repo := &mockRepository{
		daftarAntreanFunc: func(ctx context.Context, kodeDokter string, filter rawatjalan.FilterAntreanDokter) ([]rawatjalan.KunjunganRawatJalan, int, error) {
			return []rawatjalan.KunjunganRawatJalan{{NoRawat: "001", NamaPasien: "Budi"}}, 15, nil
		},
	}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)

	filter := rawatjalan.FilterAntreanDokter{Halaman: 2, Batas: 5}

	_, meta, err := uc.DaftarAntreanDokter(context.Background(), "DK001", filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.TotalHalaman != 3 {
		t.Errorf("expected TotalHalaman=3, got %d", meta.TotalHalaman)
	}
	if meta.HalamanAktif != 2 {
		t.Errorf("expected HalamanAktif=2, got %d", meta.HalamanAktif)
	}
}

func TestDetailKunjungan_EmptyNoRawat(t *testing.T) {
	repo := &mockRepository{}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)

	_, err := uc.DetailKunjungan(context.Background(), "", "DK001")
	if err == nil {
		t.Fatal("expected error for empty noRawat, got nil")
	}
}
