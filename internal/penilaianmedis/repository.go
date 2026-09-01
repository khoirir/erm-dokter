package penilaianmedis

import (
	"context"
	"database/sql"
)

type Repository interface {
	DetailPenilaianMedisRalan(ctx context.Context, noRawat string) (*PenilaianMedisRalan, error)
	RiwayatPenilaianMedisRalanByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRalan, error)
	CekPenilaianMedisRalanAda(ctx context.Context, noRawat string) (bool, error)
	SimpanPenilaianMedisRalan(ctx context.Context, noRawat, kodeDokter string, req SimpanPenilaianMedisRalanRequest) error
	UpdatePenilaianMedisRalan(ctx context.Context, noRawat string, req UpdatePenilaianMedisRalanRequest) error
	HapusPenilaianMedisRalan(ctx context.Context, noRawat string) error

	DetailPenilaianMedisIGD(ctx context.Context, noRawat string) (*PenilaianMedisIGD, error)
	RiwayatPenilaianMedisIGDByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisIGD, error)
	CekPenilaianMedisIGDAda(ctx context.Context, noRawat string) (bool, error)
	SimpanPenilaianMedisIGD(ctx context.Context, noRawat, kodeDokter string, req SimpanPenilaianMedisIGDRequest) error
	UpdatePenilaianMedisIGD(ctx context.Context, noRawat string, req UpdatePenilaianMedisIGDRequest) error
	HapusPenilaianMedisIGD(ctx context.Context, noRawat string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
