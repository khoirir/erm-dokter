package penjamin

import (
	"context"
	"fmt"

	"erm-dokter/internal/pkg/logger"
)

type Service interface {
	DaftarPenjamin(ctx context.Context) ([]Penjamin, error)
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
	penjamin, err := s.repo.DaftarPenjamin(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar penjamin: %v", err)
		return nil, fmt.Errorf("gagal mengambil daftar penjamin: %w", err)
	}

	return penjamin, nil
}
