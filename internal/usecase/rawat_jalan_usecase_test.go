package usecase_test

import (
	"context"
	"testing"

	"erm-dokter/internal/domain"
	"erm-dokter/internal/dto"
	"erm-dokter/internal/usecase"
	"erm-dokter/pkg/logger"
)

// --- Mock Repository ---

type mockRawatJalanRepository struct {
	daftarAntreanFunc   func(ctx context.Context, filter dto.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, int, error)
	detailKunjunganFunc func(ctx context.Context, noRawat string, kodeDokter string) (*domain.KunjunganRawatJalan, error)
}

func (m *mockRawatJalanRepository) DaftarAntreanDokter(ctx context.Context, filter dto.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, int, error) {
	if m.daftarAntreanFunc != nil {
		return m.daftarAntreanFunc(ctx, filter)
	}
	return []domain.KunjunganRawatJalan{}, 0, nil
}

func (m *mockRawatJalanRepository) DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*domain.KunjunganRawatJalan, error) {
	if m.detailKunjunganFunc != nil {
		return m.detailKunjunganFunc(ctx, noRawat, kodeDokter)
	}
	return nil, nil
}

// --- Tests ---

func TestDaftarAntreanDokter_ValidationError(t *testing.T) {
	repo := &mockRawatJalanRepository{}
	log := logger.New()
	uc := usecase.NewRawatJalanUsecase(repo, log)

	filter := dto.FilterAntreanDokter{
		OrderBy:   "kolom_salah",
		SortOrder: "SALAH",
	}

	_, _, err := uc.DaftarAntreanDokter(context.Background(), filter)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	validationErr, ok := err.(dto.ValidationError)
	if !ok {
		t.Fatalf("expected dto.ValidationError, got %T", err)
	}

	if _, exists := validationErr["order_by"]; !exists {
		t.Error("expected validation error for 'order_by'")
	}
	if _, exists := validationErr["sort_order"]; !exists {
		t.Error("expected validation error for 'sort_order'")
	}
}

func TestDaftarAntreanDokter_MinKeywordLength(t *testing.T) {
	repo := &mockRawatJalanRepository{}
	log := logger.New()
	uc := usecase.NewRawatJalanUsecase(repo, log)

	filter := dto.FilterAntreanDokter{KataKunci: "ab"}

	_, _, err := uc.DaftarAntreanDokter(context.Background(), filter)
	if err == nil {
		t.Fatal("expected business error for short keyword, got nil")
	}

	businessErr, ok := err.(*domain.BusinessError)
	if !ok {
		t.Fatalf("expected *domain.BusinessError, got %T", err)
	}
	if businessErr.Message != "kata kunci pencarian minimal 3 karakter" {
		t.Errorf("unexpected error message: %s", businessErr.Message)
	}
}

func TestDaftarAntreanDokter_Success(t *testing.T) {
	repo := &mockRawatJalanRepository{
		daftarAntreanFunc: func(ctx context.Context, filter dto.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, int, error) {
			return []domain.KunjunganRawatJalan{
				{NoRawat: "2025/04/22/000001", NamaPasien: "Budi"},
				{NoRawat: "2025/04/22/000002", NamaPasien: "Ani"},
			}, 2, nil
		},
	}
	log := logger.New()
	uc := usecase.NewRawatJalanUsecase(repo, log)

	filter := dto.FilterAntreanDokter{KodeDokter: "DK001"}

	data, meta, err := uc.DaftarAntreanDokter(context.Background(), filter)
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
	repo := &mockRawatJalanRepository{
		daftarAntreanFunc: func(ctx context.Context, filter dto.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, int, error) {
			return []domain.KunjunganRawatJalan{{NoRawat: "001", NamaPasien: "Budi"}}, 15, nil
		},
	}
	log := logger.New()
	uc := usecase.NewRawatJalanUsecase(repo, log)

	filter := dto.FilterAntreanDokter{KodeDokter: "DK001", Halaman: 2, Batas: 5}

	_, meta, err := uc.DaftarAntreanDokter(context.Background(), filter)
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
	repo := &mockRawatJalanRepository{}
	log := logger.New()
	uc := usecase.NewRawatJalanUsecase(repo, log)

	_, err := uc.DetailKunjungan(context.Background(), "", "DK001")
	if err == nil {
		t.Fatal("expected error for empty noRawat, got nil")
	}
}
