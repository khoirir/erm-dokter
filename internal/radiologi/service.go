package radiologi

import (
	"context"
	"database/sql"
	"errors"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	GetRiwayatRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error)
	GetRiwayatRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error)
	GetDetailHasilRadiologi(ctx context.Context, idHasil IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*HasilRadiologi, error)
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

func (s *service) GetRiwayatRadiologiKunjungan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarHasilRadiologiKunjungan(ctx, noRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat radiologi kunjungan %s: %v", noRawat, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetRiwayatRadiologiPasien(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatRadiologi) ([]HasilRadiologi, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarHasilRadiologiPasien(ctx, noRM, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat radiologi pasien %s: %v", noRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetDetailHasilRadiologi(ctx context.Context, idHasil IdHasilRadiologi, statusLanjut shared.StatusLanjut) (*HasilRadiologi, error) {
	item, err := s.repo.DetailHasilRadiologi(ctx, idHasil, statusLanjut)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data hasil pemeriksaan radiologi tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail hasil radiologi no_rawat %s kode_tindakan %s: %v", idHasil.NoRawat, idHasil.KodeTindakan, err)
		return nil, err
	}

	return item, nil
}
