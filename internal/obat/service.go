package obat

import (
	"context"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
)

type Service interface {
	DaftarObat(ctx context.Context, filter FilterDaftarObat) ([]Obat, shared.PaginationMeta, error)
	DetailObat(ctx context.Context, kodeObat string) (*Obat, error)
	DaftarJenis(ctx context.Context) ([]JenisObat, error)
	DaftarGolongan(ctx context.Context) ([]GolonganObat, error)
	DaftarKategori(ctx context.Context) ([]KategoriObat, error)
	CekKeberadaanObat(ctx context.Context, listKodeObat []string) (map[string]bool, error)
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

func (s *service) DaftarObat(ctx context.Context, filter FilterDaftarObat) ([]Obat, shared.PaginationMeta, error) {
	data, total, err := s.repo.DaftarObat(ctx, filter)
	if err != nil {
		s.log.Error("Gagal mengambil daftar obat: %v", err)
		return nil, shared.PaginationMeta{}, err
	}

	return data, shared.NewPaginationMeta(int(total), filter.Page, filter.Limit), nil
}

func (s *service) DetailObat(ctx context.Context, kodeObat string) (*Obat, error) {
	detail, err := s.repo.DetailObat(ctx, kodeObat)
	if err != nil {
		s.log.Error("Gagal mengambil detail obat %s: %v", kodeObat, err)
		return nil, err
	}

	return detail, nil
}

func (s *service) DaftarJenis(ctx context.Context) ([]JenisObat, error) {
	data, err := s.repo.DaftarJenis(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar jenis obat: %v", err)
		return nil, err
	}
	return data, nil
}

func (s *service) DaftarGolongan(ctx context.Context) ([]GolonganObat, error) {
	data, err := s.repo.DaftarGolongan(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar golongan obat: %v", err)
		return nil, err
	}
	return data, nil
}

func (s *service) DaftarKategori(ctx context.Context) ([]KategoriObat, error) {
	data, err := s.repo.DaftarKategori(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar kategori obat: %v", err)
		return nil, err
	}
	return data, nil
}

func (s *service) CekKeberadaanObat(ctx context.Context, listKodeObat []string) (map[string]bool, error) {
	res, err := s.repo.CekKeberadaanObat(ctx, listKodeObat)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan obat di repository: %v", err)
		return nil, err
	}
	return res, nil
}
