package rawatjalan_test

import (
	"context"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
)

type mockRepository struct {
	daftarAntreanFunc          func(ctx context.Context, kodeDokter string, filter rawatjalan.FilterAntreanDokter) ([]rawatjalan.KunjunganRawatJalan, int, error)
	detailKunjunganFunc        func(ctx context.Context, noRawat string, kodeDokter string) (*rawatjalan.KunjunganRawatJalan, error)
	riwayatKunjunganPasienFunc func(ctx context.Context, noRM string) ([]rawatjalan.KunjunganRawatJalan, error)
	getInfoRegistrasiFunc      func(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error)
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

func (m *mockRepository) RiwayatKunjunganPasien(ctx context.Context, noRM string) ([]rawatjalan.KunjunganRawatJalan, error) {
	if m.riwayatKunjunganPasienFunc != nil {
		return m.riwayatKunjunganPasienFunc(ctx, noRM)
	}
	return []rawatjalan.KunjunganRawatJalan{}, nil
}

func (m *mockRepository) GetWaktuRegistrasi(ctx context.Context, noRawat string) (string, string, bool, error) {
	return "2020-01-01", "00:00:00", true, nil
}

func (m *mockRepository) GetInfoRegistrasi(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error) {
	if m.getInfoRegistrasiFunc != nil {
		return m.getInfoRegistrasiFunc(ctx, noRawat)
	}
	return &rawatjalan.InfoRegistrasiPasien{
		TanggalRegistrasi: "2020-01-01",
		JamRegistrasi:     "00:00:00",
		KodePenjamin:      "UMU",
		StatusBayar:       "Belum Bayar",
	}, nil
}

func TestDaftarAntreanDokter_ValidationError(t *testing.T) {
	filter := rawatjalan.FilterAntreanDokter{
		OrderBy:   "kolom_tidak_ada",
		SortOrder: "SALAH",
	}

	validationErr := filter.Validate()
	if validationErr == nil {
		t.Fatal("expected validation error, got nil")
	}

	if _, exists := validationErr["order_by"]; !exists {
		t.Error("expected validation error for 'order_by'")
	}
	if _, exists := validationErr["sort_order"]; !exists {
		t.Error("expected validation error for 'sort_order'")
	}
}

func TestDaftarAntreanDokter_ValidOrderByOptions(t *testing.T) {
	validOrderBys := []string{
		"waktu_registrasi",
		"nama_pasien",
		"penjamin",
		"status_pemeriksaan",
		"status_lanjut",
		"status_bayar",
		"jenis_antrean",
	}

	for _, ob := range validOrderBys {
		t.Run(ob, func(t *testing.T) {
			filter := rawatjalan.FilterAntreanDokter{
				OrderBy:   ob,
				SortOrder: "ASC",
			}
			filter.Sanitize()
			if errs := filter.Validate(); errs != nil {
				t.Fatalf("expected valid filter for order_by=%s, got %v", ob, errs)
			}
		})
	}
}


func TestDaftarAntreanDokter_MinKeywordLength(t *testing.T) {
	filter := rawatjalan.FilterAntreanDokter{Keyword: "ab"}
	filter.Sanitize()

	validationErr := filter.Validate()
	if validationErr == nil {
		t.Fatal("expected validation error for short keyword, got nil")
	}

	if _, exists := validationErr["keyword"]; !exists {
		t.Error("expected validation error for 'keyword'")
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
	if meta.TotalRecords != 2 {
		t.Errorf("expected TotalRecords=2, got %d", meta.TotalRecords)
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

	filter := rawatjalan.FilterAntreanDokter{Page: 2, Limit: 5}

	_, meta, err := uc.DaftarAntreanDokter(context.Background(), "DK001", filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.TotalPages != 3 {
		t.Errorf("expected TotalPages=3, got %d", meta.TotalPages)
	}
	if meta.CurrentPage != 2 {
		t.Errorf("expected CurrentPage=2, got %d", meta.CurrentPage)
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
	if err.Error() != "Nomor rawat wajib diisi" {
		t.Errorf("expected 'Nomor rawat wajib diisi', got %q", err.Error())
	}
}

func TestRiwayatKunjunganPasien_EmptyNoRM(t *testing.T) {
	repo := &mockRepository{}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)

	_, err := uc.RiwayatKunjunganPasien(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty noRM, got nil")
	}
	if err.Error() != "Nomor rekam medis wajib diisi" {
		t.Errorf("expected 'Nomor rekam medis wajib diisi', got %q", err.Error())
	}
}

func TestGetWaktuRegistrasi_EmptyNoRawat(t *testing.T) {
	repo := &mockRepository{}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)

	_, _, _, err := uc.GetWaktuRegistrasi(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty noRawat, got nil")
	}
	if err.Error() != "Nomor rawat wajib diisi" {
		t.Errorf("expected 'Nomor rawat wajib diisi', got %q", err.Error())
	}
}

func TestRawatJalan_StaticReferences(t *testing.T) {
	repo := &mockRepository{}
	log := logger.New()
	uc := rawatjalan.NewService(repo, log)
	ctx := context.Background()

	statusPeriksa := uc.DaftarStatusPemeriksaan(ctx)
	if len(statusPeriksa) != 8 {
		t.Errorf("expected 8 status periksa items, got %d", len(statusPeriksa))
	}

	statusLanjut := uc.DaftarStatusLanjut(ctx)
	if len(statusLanjut) != 2 {
		t.Errorf("expected 2 status lanjut items, got %d", len(statusLanjut))
	}

	statusBayar := uc.DaftarStatusBayar(ctx)
	if len(statusBayar) != 2 {
		t.Errorf("expected 2 status bayar items, got %d", len(statusBayar))
	}

	jenisAntrean := uc.DaftarJenisAntrean(ctx)
	if len(jenisAntrean) != 2 {
		t.Errorf("expected 2 jenis antrean items, got %d", len(jenisAntrean))
	}
}

func TestGetInfoRegistrasi(t *testing.T) {
	ctx := context.Background()
	log := logger.New()

	t.Run("Validasi noRawat kosong", func(t *testing.T) {
		repo := &mockRepository{}
		svc := rawatjalan.NewService(repo, log)

		_, err := svc.GetInfoRegistrasi(ctx, "")
		if err == nil {
			t.Fatal("harus error jika noRawat kosong")
		}
		if err.Error() != "Nomor rawat wajib diisi" {
			t.Errorf("pesan error tidak sesuai: %v", err)
		}
	})

	t.Run("Data tidak ditemukan", func(t *testing.T) {
		repo := &mockRepository{
			getInfoRegistrasiFunc: func(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error) {
				return nil, nil
			},
		}
		svc := rawatjalan.NewService(repo, log)

		_, err := svc.GetInfoRegistrasi(ctx, "2026/09/04/000001")
		if err == nil {
			t.Fatal("harus error not found jika data nil")
		}
		if err.Error() != "Data kunjungan pasien tidak ditemukan" {
			t.Errorf("pesan error tidak sesuai: %v", err)
		}
	})

	t.Run("Sukses", func(t *testing.T) {
		repo := &mockRepository{
			getInfoRegistrasiFunc: func(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error) {
				return &rawatjalan.InfoRegistrasiPasien{
					TanggalRegistrasi: "2026-09-04",
					JamRegistrasi:     "08:00:00",
					KodePenjamin:      "BPJ",
					StatusBayar:       "Sudah Bayar",
				}, nil
			},
		}
		svc := rawatjalan.NewService(repo, log)

		info, err := svc.GetInfoRegistrasi(ctx, "2026/09/04/000001")
		if err != nil {
			t.Fatalf("tidak diharapkan error: %v", err)
		}
		if info.KodePenjamin != "BPJ" || info.StatusBayar != "Sudah Bayar" {
			t.Errorf("data info registrasi tidak cocok: %+v", info)
		}
	})
}



