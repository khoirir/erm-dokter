package usecase

import (
	"context"

	"erm-dokter/internal/domain"
	"erm-dokter/pkg/logger"
)

type penjaminUsecase struct {
	penjaminRepo domain.PenjaminRepository
	Log          *logger.Logger
}

func NewPenjaminUsecase(repo domain.PenjaminRepository, log *logger.Logger) domain.PenjaminUsecase {
	return &penjaminUsecase{
		penjaminRepo: repo,
		Log:          log,
	}
}

func (u *penjaminUsecase) DaftarPenjamin(ctx context.Context) ([]domain.Penjamin, error) {
	data, err := u.penjaminRepo.DaftarPenjamin(ctx)
	if err != nil {
		u.Log.Error("Gagal query daftar penjamin: %v", err)
		return nil, err
	}
	return data, nil
}