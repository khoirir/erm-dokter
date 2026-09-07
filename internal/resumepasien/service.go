package resumepasien

import (
	"context"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
)

type Service interface {
	// Ralan
	DetailResumePasienRalan(ctx context.Context, noRawat string) (*ResumePasienRalan, error)
	RiwayatResumePasienRalanByNoRM(ctx context.Context, noRM string) ([]ResumePasienRalan, error)
	SimpanResumePasienRalan(ctx context.Context, noRawat, kodeDokter string, req SimpanResumePasienRalanRequest) (*ResumePasienRalan, error)
	UpdateResumePasienRalan(ctx context.Context, kodeDokterLogin, noRawat string, req UpdateResumePasienRalanRequest) (*ResumePasienRalan, error)
	HapusResumePasienRalan(ctx context.Context, kodeDokterLogin, noRawat string) error

	// Ranap
	DetailResumePasienRanap(ctx context.Context, noRawat string) (*ResumePasienRanap, error)
	RiwayatResumePasienRanapByNoRM(ctx context.Context, noRM string) ([]ResumePasienRanap, error)
	SimpanResumePasienRanap(ctx context.Context, noRawat, kodeDokter string, req SimpanResumePasienRanapRequest) (*ResumePasienRanap, error)
	UpdateResumePasienRanap(ctx context.Context, kodeDokterLogin, noRawat string, req UpdateResumePasienRanapRequest) (*ResumePasienRanap, error)
	HapusResumePasienRanap(ctx context.Context, kodeDokterLogin, noRawat string) error

	// Referensi Dropdown
	ReferensiRanap(ctx context.Context) ReferensiResumeRanap
	ReferensiRalan(ctx context.Context) ReferensiResumeRalan
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	rawatInapService  rawatinap.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(
	repo Repository,
	rawatJalanService rawatjalan.Service,
	rawatInapService rawatinap.Service,
	maxEditJam int,
	log *logger.Logger,
) Service {
	if maxEditJam <= 0 {
		maxEditJam = 48
	}
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		rawatInapService:  rawatInapService,
		maxEditJam:        maxEditJam,
		log:               log,
	}
}
