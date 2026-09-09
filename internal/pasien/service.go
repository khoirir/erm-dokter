package pasien

import (
	"context"
	"strings"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DetailPasien(ctx context.Context, noRM string) (*Pasien, error)
	CariPasien(ctx context.Context, kataKunci string) ([]Pasien, error)
	GetNoRMByNoRawat(ctx context.Context, noRawat string) (string, error)
	RiwayatKunjunganPasien(ctx context.Context, noRM string, filter FilterRiwayatKunjungan) ([]RiwayatKunjungan, shared.PaginationMeta, error)
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

func (s *service) DetailPasien(ctx context.Context, noRM string) (*Pasien, error) {
	cleanNoRM := strings.TrimSpace(noRM)
	if cleanNoRM == "" {
		return nil, apperror.NewBusinessError("Nomor rekam medis wajib diisi")
	}

	p, err := s.repo.CariByNoRM(ctx, cleanNoRM)
	if err != nil {
		s.log.Error("Gagal mengambil detail pasien %s: %v", cleanNoRM, err)
		return nil, err
	}
	return p, nil
}

func (s *service) CariPasien(ctx context.Context, kataKunci string) ([]Pasien, error) {
	cleanKunci := strings.TrimSpace(kataKunci)
	if len(cleanKunci) < 3 {
		return nil, apperror.NewBusinessError("Kata kunci pencarian minimal 3 karakter")
	}

	list, err := s.repo.CariPasien(ctx, cleanKunci, 20)
	if err != nil {
		s.log.Error("Gagal mencari pasien dengan kata kunci %s: %v", cleanKunci, err)
		return nil, err
	}
	return list, nil
}

func (s *service) GetNoRMByNoRawat(ctx context.Context, noRawat string) (string, error) {
	cleanNoRawat := strings.TrimSpace(noRawat)
	if cleanNoRawat == "" {
		return "", apperror.NewBusinessError("Nomor rawat wajib diisi")
	}

	noRM, err := s.repo.GetNoRMByNoRawat(ctx, cleanNoRawat)
	if err != nil {
		s.log.Error("Gagal mengambil no_rkm_medis dari no_rawat %s: %v", cleanNoRawat, err)
		return "", err
	}
	return noRM, nil
}

func (s *service) RiwayatKunjunganPasien(ctx context.Context, noRM string, filter FilterRiwayatKunjungan) ([]RiwayatKunjungan, shared.PaginationMeta, error) {
	cleanNoRM := strings.TrimSpace(noRM)
	if cleanNoRM == "" {
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("Nomor rekam medis wajib diisi")
	}

	list, total, err := s.repo.RiwayatKunjunganPasien(ctx, cleanNoRM, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat kunjungan pasien %s: %v", cleanNoRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}
