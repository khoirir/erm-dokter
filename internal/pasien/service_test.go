package pasien_test

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/pasien"
	"erm-dokter/internal/shared"
)

type mockRepository struct {
	cariByNoRMFn              func(ctx context.Context, noRM string) (*pasien.Pasien, error)
	cariByNIKFn               func(ctx context.Context, nik string) (*pasien.Pasien, error)
	cariPasienFn              func(ctx context.Context, kataKunci string, limit int) ([]pasien.Pasien, error)
	getNoRMByNoRawatFn        func(ctx context.Context, noRawat string) (string, error)
	riwayatKunjunganPasienFn  func(ctx context.Context, noRM string, filter pasien.FilterRiwayatKunjungan) ([]pasien.RiwayatKunjungan, int, error)
}

func (m *mockRepository) CariByNoRM(ctx context.Context, noRM string) (*pasien.Pasien, error) {
	if m.cariByNoRMFn != nil {
		return m.cariByNoRMFn(ctx, noRM)
	}
	return nil, nil
}

func (m *mockRepository) CariByNIK(ctx context.Context, nik string) (*pasien.Pasien, error) {
	if m.cariByNIKFn != nil {
		return m.cariByNIKFn(ctx, nik)
	}
	return nil, nil
}

func (m *mockRepository) CariPasien(ctx context.Context, kataKunci string, limit int) ([]pasien.Pasien, error) {
	if m.cariPasienFn != nil {
		return m.cariPasienFn(ctx, kataKunci, limit)
	}
	return nil, nil
}

func (m *mockRepository) GetNoRMByNoRawat(ctx context.Context, noRawat string) (string, error) {
	if m.getNoRMByNoRawatFn != nil {
		return m.getNoRMByNoRawatFn(ctx, noRawat)
	}
	return "", nil
}

func (m *mockRepository) RiwayatKunjunganPasien(ctx context.Context, noRM string, filter pasien.FilterRiwayatKunjungan) ([]pasien.RiwayatKunjungan, int, error) {
	if m.riwayatKunjunganPasienFn != nil {
		return m.riwayatKunjunganPasienFn(ctx, noRM, filter)
	}
	return nil, 0, nil
}

func TestService_DetailPasien(t *testing.T) {
	log := logger.New()

	t.Run("Empty noRM returns error", func(t *testing.T) {
		svc := pasien.NewService(&mockRepository{}, log)
		_, err := svc.DetailPasien(context.Background(), "   ")
		if err == nil || err.Error() != "Nomor rekam medis wajib diisi" {
			t.Fatalf("expected Nomor rekam medis wajib diisi, got %v", err)
		}
	})

	t.Run("Success returns Pasien", func(t *testing.T) {
		mockRepo := &mockRepository{
			cariByNoRMFn: func(ctx context.Context, noRM string) (*pasien.Pasien, error) {
				return &pasien.Pasien{NoRM: noRM, Nama: "Budi"}, nil
			},
		}
		svc := pasien.NewService(mockRepo, log)
		p, err := svc.DetailPasien(context.Background(), "123456")
		if err != nil || p.Nama != "Budi" {
			t.Fatalf("unexpected result: p=%+v, err=%v", p, err)
		}
	})

	t.Run("Repo error propagated", func(t *testing.T) {
		repoErr := errors.New("db error")
		mockRepo := &mockRepository{
			cariByNoRMFn: func(ctx context.Context, noRM string) (*pasien.Pasien, error) {
				return nil, repoErr
			},
		}
		svc := pasien.NewService(mockRepo, log)
		_, err := svc.DetailPasien(context.Background(), "123456")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestService_CariPasien(t *testing.T) {
	log := logger.New()

	t.Run("Short keyword returns error", func(t *testing.T) {
		svc := pasien.NewService(&mockRepository{}, log)
		_, err := svc.CariPasien(context.Background(), "bu")
		if err == nil || err.Error() != "Kata kunci pencarian minimal 3 karakter" {
			t.Fatalf("expected error minimal 3 karakter, got %v", err)
		}
	})

	t.Run("Success returns list", func(t *testing.T) {
		mockRepo := &mockRepository{
			cariPasienFn: func(ctx context.Context, kataKunci string, limit int) ([]pasien.Pasien, error) {
				return []pasien.Pasien{{Nama: "Budi"}}, nil
			},
		}
		svc := pasien.NewService(mockRepo, log)
		list, err := svc.CariPasien(context.Background(), "budi")
		if err != nil || len(list) != 1 {
			t.Fatalf("unexpected result: len=%d, err=%v", len(list), err)
		}
	})
}

func TestService_GetNoRMByNoRawat(t *testing.T) {
	log := logger.New()

	t.Run("Empty noRawat returns error", func(t *testing.T) {
		svc := pasien.NewService(&mockRepository{}, log)
		_, err := svc.GetNoRMByNoRawat(context.Background(), "")
		if err == nil || err.Error() != "Nomor rawat wajib diisi" {
			t.Fatalf("expected error Nomor rawat wajib diisi, got %v", err)
		}
	})

	t.Run("Success returns noRM", func(t *testing.T) {
		mockRepo := &mockRepository{
			getNoRMByNoRawatFn: func(ctx context.Context, noRawat string) (string, error) {
				return "00123456", nil
			},
		}
		svc := pasien.NewService(mockRepo, log)
		noRM, err := svc.GetNoRMByNoRawat(context.Background(), "2026/09/01/000001")
		if err != nil || noRM != "00123456" {
			t.Fatalf("unexpected noRM %s, err: %v", noRM, err)
		}
	})
}

func TestService_RiwayatKunjunganPasien(t *testing.T) {
	log := logger.New()

	t.Run("Empty noRM returns error", func(t *testing.T) {
		svc := pasien.NewService(&mockRepository{}, log)
		_, _, err := svc.RiwayatKunjunganPasien(context.Background(), "", pasien.FilterRiwayatKunjungan{})
		if err == nil || err.Error() != "Nomor rekam medis wajib diisi" {
			t.Fatalf("expected error Nomor rekam medis wajib diisi, got %v", err)
		}
	})

	t.Run("Success returns list and pagination", func(t *testing.T) {
		mockRepo := &mockRepository{
			riwayatKunjunganPasienFn: func(ctx context.Context, noRM string, filter pasien.FilterRiwayatKunjungan) ([]pasien.RiwayatKunjungan, int, error) {
				return []pasien.RiwayatKunjungan{
					{
						NoRawat:           "2026/09/01/000001",
						TanggalRegistrasi: "2026-09-01",
						StatusLanjut:      shared.StatusLanjutRawatJalan,
						NamaPoli:          "Poli Penyakit Dalam",
					},
					{
						NoRawat:           "2026/08/15/000002",
						TanggalRegistrasi: "2026-08-15",
						StatusLanjut:      shared.StatusLanjutRawatInap,
						KamarInap: &pasien.RiwayatRawatInap{
							DPJP: []string{"dr. Budi, Sp.PD"},
							Kamar: []pasien.RiwayatKamarInap{
								{
									KodeKamar:   "K01",
									NamaBangsal: "Bangsal Melati",
									Kelas:       "Kelas 1",
								},
							},
						},
					},
				}, 12, nil
			},
		}

		svc := pasien.NewService(mockRepo, log)
		filter := pasien.FilterRiwayatKunjungan{Page: 1, Limit: 5}
		list, meta, err := svc.RiwayatKunjunganPasien(context.Background(), "00123456", filter)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if meta.TotalRecords != 12 {
			t.Errorf("expected 12 total records, got %d", meta.TotalRecords)
		}
		if meta.TotalPages != 3 {
			t.Errorf("expected 3 total pages, got %d", meta.TotalPages)
		}
	})

	t.Run("Repository error propagated", func(t *testing.T) {
		repoErr := errors.New("db error")
		mockRepo := &mockRepository{
			riwayatKunjunganPasienFn: func(ctx context.Context, noRM string, filter pasien.FilterRiwayatKunjungan) ([]pasien.RiwayatKunjungan, int, error) {
				return nil, 0, repoErr
			},
		}
		svc := pasien.NewService(mockRepo, log)
		_, _, err := svc.RiwayatKunjunganPasien(context.Background(), "00123456", pasien.FilterRiwayatKunjungan{Page: 1, Limit: 5})
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}
