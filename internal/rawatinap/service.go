package rawatinap

import (
	"context"
	"strings"

	"erm-dokter/internal/pkg/logger"
)

type Service interface {
	CekStatusKamarInap(ctx context.Context, noRawat string) (isKamarAktif bool, hasRecordKamar bool, err error)
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

func (s *service) CekStatusKamarInap(ctx context.Context, noRawat string) (bool, bool, error) {
	cleanNoRawat := strings.TrimSpace(noRawat)
	if cleanNoRawat == "" {
		return false, false, nil
	}

	isAktif, hasRecord, err := s.repo.CekStatusKamarInap(ctx, cleanNoRawat)
	if err != nil {
		s.log.Error("Gagal memeriksa status kamar inap pasien %s: %v", cleanNoRawat, err)
		return false, false, err
	}

	return isAktif, hasRecord, nil
}
