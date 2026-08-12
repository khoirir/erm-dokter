package pemeriksaan

import (
	"context"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]Pemeriksaan, error)
	// DetailKunjungan(ctx context.Context, noRawat string, kodeDokter string) (*KunjunganRawatJalan, error)
	// GetReferensiFilter(ctx context.Context) ReferensiFilterRawatJalan
}

type service struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]Pemeriksaan, error) {
	if noRawat == "" {
		return nil, apperror.NewBusinessError("nomor rawat tidak boleh kosong")
	}

	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		return nil, apperror.NewBusinessError("status lanjut tidak valid")
	}

	daftarPemeriksaan, err := s.repo.DaftarPemeriksaan(ctx, noRawat, statusLanjut)
	if err != nil {
		s.log.Error("Gagal query daftar pemeriksaan %s: %v", noRawat, err)
		return nil, err
	}

	return daftarPemeriksaan, nil
}
