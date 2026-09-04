package penilaianmedis

import (
	"context"
	"fmt"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	Referensi(ctx context.Context) ReferensiPenilaianMedis

	DetailPenilaianMedisRalan(ctx context.Context, noRawat string) (*PenilaianMedisRalan, error)
	RiwayatPenilaianMedisRalanByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRalan, error)
	SimpanPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisRalanRequest) (*PenilaianMedisRalan, error)
	UpdatePenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisRalanRequest) (*PenilaianMedisRalan, error)
	HapusPenilaianMedisRalan(ctx context.Context, kodeDokterLogin, noRawat string) error

	DetailPenilaianMedisIGD(ctx context.Context, noRawat string) (*PenilaianMedisIGD, error)
	RiwayatPenilaianMedisIGDByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisIGD, error)
	SimpanPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisIGDRequest) (*PenilaianMedisIGD, error)
	UpdatePenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisIGDRequest) (*PenilaianMedisIGD, error)
	HapusPenilaianMedisIGD(ctx context.Context, kodeDokterLogin, noRawat string) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	log               *logger.Logger
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, log *logger.Logger) Service {
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		log:               log,
	}
}

func (s *service) Referensi(ctx context.Context) ReferensiPenilaianMedis {
	return GetReferensiPenilaianMedis()
}

func (s *service) validasiWaktuPenilaianMedis(ctx context.Context, noRawat, tanggalPeriksa, aksi string) error {
	tanggalRegistrasiStr, jamRegistrasiStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	if err := shared.ValidasiBatasWaktuRekamMedis(tanggalRegistrasiStr, jamRegistrasiStr, 48, aksi); err != nil {
		s.log.Warn("Validasi batas waktu 48 jam penilaian medis gagal untuk no_rawat %s (%s %s): %v", noRawat, tanggalRegistrasiStr, jamRegistrasiStr, err)
		return err
	}

	if tanggalPeriksa != "" {
		waktuRegistrasi, errReg := shared.ParseWaktu(tanggalRegistrasiStr, jamRegistrasiStr)
		waktuPemeriksaan, errPer := time.ParseInLocation("2006-01-02 15:04:05", tanggalPeriksa, time.Local)
		if errReg == nil && errPer == nil && waktuPemeriksaan.Before(waktuRegistrasi) {
			errs := apperror.ValidationError{
				"tanggal_penilaian": fmt.Sprintf("Waktu penilaian medis (%s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tanggalPeriksa, tanggalRegistrasiStr, jamRegistrasiStr),
			}
			s.log.Warn("Validasi waktu penilaian medis gagal untuk no_rawat %s: %+v", noRawat, errs)
			return errs
		}
	}

	return nil
}
