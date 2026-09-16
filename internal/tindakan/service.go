package tindakan

import (
	"context"
	"database/sql"
	"errors"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	GetDaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter FilterDaftarTindakanLab) ([]TindakanLab, shared.PaginationMeta, error)
	GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*DetailTindakanLab, error)
	CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error)
	CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error)

	GetDaftarTindakanRadiologi(ctx context.Context, filter FilterDaftarTindakanRadiologi) ([]TindakanRadiologi, shared.PaginationMeta, error)
	GetDetailTindakanRadiologi(ctx context.Context, kodeTindakan string) (*TindakanRadiologi, error)
	CekKeberadaanTindakanRadiologi(ctx context.Context, listKodeTindakan []string) (map[string]bool, error)
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

func (s *service) GetDaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter FilterDaftarTindakanLab) ([]TindakanLab, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarTindakanLab(ctx, kategori, filter)
	if err != nil {
		s.log.Error("gagal mengambil daftar tindakan laboratorium kategori %s: %v", kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, kodeTindakan string) (*DetailTindakanLab, error) {
	tindakan, templates, err := s.repo.GetDetailTindakanLab(ctx, kategori, kodeTindakan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Tindakan laboratorium tidak ditemukan")
		}
		s.log.Error("gagal mengambil detail tindakan laboratorium %s kategori %s: %v", kodeTindakan, kategori, err)
		return nil, err
	}

	return &DetailTindakanLab{
		TindakanLab: *tindakan,
		Templates:   templates,
	}, nil
}

func (s *service) CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
	res, err := s.repo.CekKeberadaanTindakanLab(ctx, kategori, listKodeTindakan)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan tindakan lab di repository: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *service) CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
	res, err := s.repo.CekKeberadaanTemplateLab(ctx, listKodeTindakan, templateMap)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan template lab di repository: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *service) GetDaftarTindakanRadiologi(ctx context.Context, filter FilterDaftarTindakanRadiologi) ([]TindakanRadiologi, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarTindakanRadiologi(ctx, filter)
	if err != nil {
		s.log.Error("gagal mengambil daftar tindakan radiologi: %v", err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetDetailTindakanRadiologi(ctx context.Context, kodeTindakan string) (*TindakanRadiologi, error) {
	tindakan, err := s.repo.GetDetailTindakanRadiologi(ctx, kodeTindakan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Tindakan radiologi tidak ditemukan")
		}
		s.log.Error("gagal mengambil detail tindakan radiologi %s: %v", kodeTindakan, err)
		return nil, err
	}

	return tindakan, nil
}

func (s *service) CekKeberadaanTindakanRadiologi(ctx context.Context, listKodeTindakan []string) (map[string]bool, error) {
	res, err := s.repo.CekKeberadaanTindakanRadiologi(ctx, listKodeTindakan)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan tindakan radiologi di repository: %v", err)
		return nil, err
	}
	return res, nil
}


