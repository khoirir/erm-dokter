package resep

import (
	"context"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
)

type Service interface {
	DaftarResep(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error)
	DaftarResepByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error)
	DaftarAturanPakai(ctx context.Context, keyword string) ([]AturanPakai, error)
	DaftarMetodeRacik(ctx context.Context) ([]MetodeRacik, error)
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

func (s *service) DaftarResep(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error) {
	data, total, err := s.repo.DaftarResep(ctx, noRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil daftar resep untuk no_rawat %s: %v", noRawat, err)
		return nil, shared.PaginationMeta{}, err
	}

	return data, shared.NewPaginationMeta(total, filter.Page, filter.Limit), nil
}

func (s *service) DaftarResepByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error) {
	data, total, err := s.repo.DaftarResepByRM(ctx, noRM, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat resep untuk no_rm %s: %v", noRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	return data, shared.NewPaginationMeta(total, filter.Page, filter.Limit), nil
}

func (s *service) DaftarAturanPakai(ctx context.Context, keyword string) ([]AturanPakai, error) {
	data, err := s.repo.DaftarAturanPakai(ctx, keyword)
	if err != nil {
		s.log.Error("Gagal mengambil daftar aturan pakai: %v", err)
		return nil, err
	}

	return data, nil
}

func (s *service) DaftarMetodeRacik(ctx context.Context) ([]MetodeRacik, error) {
	data, err := s.repo.DaftarMetodeRacik(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar metode racik: %v", err)
		return nil, err
	}

	return data, nil
}
