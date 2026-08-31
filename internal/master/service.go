package master

import (
	"context"

	"erm-dokter/internal/pkg/logger"
)

type Service interface {
	DaftarPenjamin(ctx context.Context) ([]Penjamin, error)
	DaftarDepo(ctx context.Context) ([]Depo, error)
	DaftarPoliklinik(ctx context.Context) ([]Poliklinik, error)
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

func (s *service) DaftarPenjamin(ctx context.Context) ([]Penjamin, error) {
	data, err := s.repo.DaftarPenjamin(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar penjamin: %v", err)
		return nil, err
	}
	return data, nil
}

func (s *service) DaftarDepo(ctx context.Context) ([]Depo, error) {
	data, err := s.repo.DaftarDepo(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar depo: %v", err)
		return nil, err
	}
	return data, nil
}

func (s *service) DaftarPoliklinik(ctx context.Context) ([]Poliklinik, error) {
	data, err := s.repo.DaftarPoliklinik(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar poliklinik: %v", err)
		return nil, err
	}
	return data, nil
}

