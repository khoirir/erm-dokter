package usecase

import (
	"context"
	"erm-dokter/internal/domain"
)

type penjaminUsecase struct {
	penjaminRepo domain.PenjaminRepository
}

func NewPenjaminUsecase(repo domain.PenjaminRepository) domain.PenjaminUsecase {
	return &penjaminUsecase{
		penjaminRepo: repo,
	}
}

func (u *penjaminUsecase) DaftarPenjamin(ctx context.Context) ([]domain.Penjamin, error) {
	return u.penjaminRepo.DaftarPenjamin(ctx)
}