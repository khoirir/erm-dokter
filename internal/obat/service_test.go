package obat_test

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/obat"
	"erm-dokter/internal/pkg/logger"
)

type mockRepository struct {
	obatData     []obat.Obat
	obatTotal    int64
	detailData   *obat.Obat
	jenisData    []obat.JenisObat
	golonganData []obat.GolonganObat
	kategoriData []obat.KategoriObat
	err          error
}

func (m *mockRepository) DaftarObat(ctx context.Context, filter obat.FilterDaftarObat) ([]obat.Obat, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.obatData, m.obatTotal, nil
}

func (m *mockRepository) DetailObat(ctx context.Context, kodeObat string) (*obat.Obat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.detailData, nil
}

func (m *mockRepository) DaftarJenis(ctx context.Context) ([]obat.JenisObat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.jenisData, nil
}

func (m *mockRepository) DaftarGolongan(ctx context.Context) ([]obat.GolonganObat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.golonganData, nil
}

func (m *mockRepository) DaftarKategori(ctx context.Context) ([]obat.KategoriObat, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.kategoriData, nil
}

func TestObatService_DaftarObat(t *testing.T) {
	repo := &mockRepository{
		obatData: []obat.Obat{
			{
				KodeObat: "B001",
				NamaObat: "Paracetamol 500mg",
				Stok:     100,
				Harga:    "Rp 5.000,00",
			},
		},
		obatTotal: 1,
	}
	svc := obat.NewService(repo, logger.New())

	data, meta, err := svc.DaftarObat(context.Background(), obat.FilterDaftarObat{
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

	// Test min keyword length (< 3 chars) via Validate()
	filter := obat.FilterDaftarObat{Keyword: "pa"}
	validationErr := filter.Validate()
	if validationErr == nil {
		t.Fatal("expected validation error for keyword < 3 chars, got nil")
	}
	if _, exists := validationErr["keyword"]; !exists {
		t.Error("expected validation error for 'keyword'")
	}

	repo.err = errors.New("db error")
	_, _, err = svc.DaftarObat(context.Background(), obat.FilterDaftarObat{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestObatService_DetailObat(t *testing.T) {
	repo := &mockRepository{
		detailData: &obat.Obat{
			KodeObat: "B001",
			NamaObat: "Paracetamol 500mg",
			Stok:     150,
			StokDepo: []obat.StokDepo{
				{KodeDepo: "DPRJ", NamaDepo: "Depo Rawat Jalan", Stok: 100},
				{KodeDepo: "DPIGD", NamaDepo: "Depo IGD", Stok: 50},
			},
		},
	}
	svc := obat.NewService(repo, logger.New())

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
}

func TestObatService_DaftarJenis(t *testing.T) {
	repo := &mockRepository{
		jenisData: []obat.JenisObat{
			{Kode: "J01", Nama: "Tablet"},
		},
	}
	svc := obat.NewService(repo, logger.New())

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
		golonganData: []obat.GolonganObat{
			{Kode: "G01", Nama: "Obat Bebas"},
		},
	}
	svc := obat.NewService(repo, logger.New())

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
		kategoriData: []obat.KategoriObat{
			{Kode: "K01", Nama: "Antibiotik"},
		},
	}
	svc := obat.NewService(repo, logger.New())

	data, err := svc.DaftarKategori(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "K01" {
		t.Errorf("unexpected kategori data: %+v", data)
	}
}
