package obat

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared/apperror"
)

type mockRepository struct {
	obatData     []Obat
	obatTotal    int64
	detailData   *Obat
	jenisData    []JenisObat
	golonganData []GolonganObat
	kategoriData []KategoriObat
	err          error
}

func (m *mockRepository) DaftarObat(ctx context.Context, filter FilterDaftarObat) ([]Obat, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.obatData, m.obatTotal, nil
}

func (m *mockRepository) DetailObat(ctx context.Context, kodeObat string) (*Obat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.detailData, nil
}

func (m *mockRepository) DaftarJenis(ctx context.Context) ([]JenisObat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.jenisData, nil
}

func (m *mockRepository) DaftarGolongan(ctx context.Context) ([]GolonganObat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.golonganData, nil
}

func (m *mockRepository) DaftarKategori(ctx context.Context) ([]KategoriObat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.kategoriData, nil
}

func (m *mockRepository) CekKeberadaanObat(ctx context.Context, listKodeObat []string) (map[string]bool, error) {
	if m.err != nil {
		return nil, m.err
	}
	res := make(map[string]bool)
	for _, k := range listKodeObat {
		res[k] = true
	}
	return res, nil
}

func TestObatService_DaftarObat(t *testing.T) {
	repo := &mockRepository{
		obatData: []Obat{
			{
				KodeObat: "B001",
				NamaObat: "Paracetamol 500mg",
				Stok:     100,
				Harga:    "Rp 5.000,00",
			},
		},
		obatTotal: 1,
	}
	svc := NewService(repo, logger.New())

	data, meta, err := svc.DaftarObat(context.Background(), FilterDaftarObat{
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].KodeObat != "B001" {
		t.Errorf("unexpected obat data: %+v", data)
	}
	if meta.TotalRecords != 1 || meta.TotalPages != 1 {
		t.Errorf("unexpected pagination meta: %+v", meta)
	}

	// Test valid order_by options via Validate()
	validOrderBys := []string{
		"nama_obat",
		"stok",
		"nama_depo",
		"nama_jenis",
		"nama_golongan",
		"nama_kategori",
	}
	for _, ob := range validOrderBys {
		t.Run("ValidOrderBy_"+ob, func(t *testing.T) {
			f := FilterDaftarObat{
				OrderBy:   ob,
				SortOrder: "ASC",
				Page:      1,
				Limit:     20,
			}
			f.Sanitize()
			if errs := f.Validate(); errs != nil {
				t.Fatalf("expected nil validation errors for valid filter order_by=%s, got %v", ob, errs)
			}
		})
	}

	// Test invalid order_by
	invalidOrderFilter := FilterDaftarObat{OrderBy: "invalid_column", SortOrder: "ASC"}
	invalidOrderFilter.Sanitize()
	if errs := invalidOrderFilter.Validate(); errs == nil || errs["order_by"] == "" {
		t.Errorf("expected validation error for invalid order_by, got %v", errs)
	}

	// Test invalid sort_order
	invalidSortFilter := FilterDaftarObat{SortOrder: "SIDEWAYS"}
	invalidSortFilter.Sanitize()
	if errs := invalidSortFilter.Validate(); errs == nil || errs["sort_order"] == "" {
		t.Errorf("expected validation error for invalid sort_order, got %v", errs)
	}

	// Test min keyword length (< 3 chars) via Validate()
	filter := FilterDaftarObat{Keyword: "pa"}
	validationErr := filter.Validate()
	if validationErr == nil {
		t.Fatal("expected validation error for keyword < 3 chars, got nil")
	}
	if _, exists := validationErr["keyword"]; !exists {
		t.Error("expected validation error for 'keyword'")
	}

	repo.err = errors.New("db error")
	_, _, err = svc.DaftarObat(context.Background(), FilterDaftarObat{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}


func TestObatService_DetailObat(t *testing.T) {
	repo := &mockRepository{
		detailData: &Obat{
			KodeObat: "B001",
			NamaObat: "Paracetamol 500mg",
			Stok:     150,
			StokDepo: []StokDepo{
				{KodeDepo: "DPRJ", NamaDepo: "Depo Rawat Jalan", Stok: 100},
				{KodeDepo: "DPIGD", NamaDepo: "Depo IGD", Stok: 50},
			},
		},
	}
	svc := NewService(repo, logger.New())

	// Test Success
	detail, err := svc.DetailObat(context.Background(), "B001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if detail.KodeObat != "B001" || len(detail.StokDepo) != 2 {
		t.Errorf("unexpected detail: %+v", detail)
	}

	// Test Error DB
	repo.err = errors.New("db error")
	_, err = svc.DetailObat(context.Background(), "B001")
	if err == nil {
		t.Fatal("expected db error, got nil")
	}

	// Test Not Found
	repo.err = nil
	repo.detailData = nil
	_, err = svc.DetailObat(context.Background(), "B999")
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	var notFoundErr *apperror.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("expected *apperror.NotFoundError, got %T", err)
	}
}

func TestObatService_DaftarJenis(t *testing.T) {
	repo := &mockRepository{
		jenisData: []JenisObat{
			{Kode: "J01", Nama: "Tablet"},
		},
	}
	svc := NewService(repo, logger.New())

	data, err := svc.DaftarJenis(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "J01" {
		t.Errorf("unexpected jenis data: %+v", data)
	}
}

func TestObatService_DaftarGolongan(t *testing.T) {
	repo := &mockRepository{
		golonganData: []GolonganObat{
			{Kode: "G01", Nama: "Obat Bebas"},
		},
	}
	svc := NewService(repo, logger.New())

	data, err := svc.DaftarGolongan(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "G01" {
		t.Errorf("unexpected golongan data: %+v", data)
	}
}

func TestObatService_DaftarKategori(t *testing.T) {
	repo := &mockRepository{
		kategoriData: []KategoriObat{
			{Kode: "K01", Nama: "Antibiotik"},
		},
	}
	svc := NewService(repo, logger.New())

	data, err := svc.DaftarKategori(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "K01" {
		t.Errorf("unexpected kategori data: %+v", data)
	}
}

func TestObatService_CekKeberadaanObat(t *testing.T) {
	repo := &mockRepository{}
	svc := NewService(repo, logger.New())

	data, err := svc.CekKeberadaanObat(context.Background(), []string{"B001", "B002"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 2 || !data["B001"] || !data["B002"] {
		t.Errorf("unexpected cek keberadaan data: %+v", data)
	}
}

