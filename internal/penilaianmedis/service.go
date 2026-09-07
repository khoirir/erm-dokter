package penilaianmedis

import (
	"context"
	"fmt"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
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

	DetailPenilaianMedisRanap(ctx context.Context, noRawat string) (*PenilaianMedisRanap, error)
	RiwayatPenilaianMedisRanapByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRanap, error)
	SimpanPenilaianMedisRanap(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisRanapRequest) (*PenilaianMedisRanap, error)
	UpdatePenilaianMedisRanap(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisRanapRequest) (*PenilaianMedisRanap, error)
	HapusPenilaianMedisRanap(ctx context.Context, kodeDokterLogin, noRawat string) error

	DetailPenilaianMedisRalanKandungan(ctx context.Context, noRawat string) (*PenilaianMedisRalanKandungan, error)
	RiwayatPenilaianMedisRalanKandunganByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRalanKandungan, error)
	SimpanPenilaianMedisRalanKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisRalanKandunganRequest) (*PenilaianMedisRalanKandungan, error)
	UpdatePenilaianMedisRalanKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisRalanKandunganRequest) (*PenilaianMedisRalanKandungan, error)
	HapusPenilaianMedisRalanKandungan(ctx context.Context, kodeDokterLogin, noRawat string) error

	DetailPenilaianMedisRanapKandungan(ctx context.Context, noRawat string) (*PenilaianMedisRanapKandungan, error)
	RiwayatPenilaianMedisRanapKandunganByNoRM(ctx context.Context, noRM string) ([]PenilaianMedisRanapKandungan, error)
	SimpanPenilaianMedisRanapKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanPenilaianMedisRanapKandunganRequest) (*PenilaianMedisRanapKandungan, error)
	UpdatePenilaianMedisRanapKandungan(ctx context.Context, kodeDokterLogin, noRawat string, req UpdatePenilaianMedisRanapKandunganRequest) (*PenilaianMedisRanapKandungan, error)
	HapusPenilaianMedisRanapKandungan(ctx context.Context, kodeDokterLogin, noRawat string) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	rawatInapService  rawatinap.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, rawatInapService rawatinap.Service, maxEditJam int, log *logger.Logger) Service {
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

func (s *service) Referensi(ctx context.Context) ReferensiPenilaianMedis {
	return GetReferensiPenilaianMedis()
}

func (s *service) validasiRegistrasiDanStatus(ctx context.Context, noRawat, tglPenilaian string) error {
	tanggalRegistrasiStr, jamRegistrasiStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	waktuRegistrasi, err := shared.ParseWaktu(tanggalRegistrasiStr, jamRegistrasiStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tanggalRegistrasiStr, jamRegistrasiStr, err)
		return err
	}

	if tglPenilaian != "" {
		waktuPenilaian, errPer := time.ParseInLocation("2006-01-02 15:04:05", tglPenilaian, time.Local)
		if errPer == nil && waktuPenilaian.Before(waktuRegistrasi) {
			errs := apperror.ValidationError{
				"tanggal_penilaian": fmt.Sprintf("Waktu penilaian medis (%s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tglPenilaian, tanggalRegistrasiStr, jamRegistrasiStr),
			}
			s.log.Warn("Validasi waktu penilaian medis gagal untuk no_rawat %s: %+v", noRawat, errs)
			return errs
		}
	}

	return s.validasiStatusKamarDanBatasWaktu(ctx, noRawat, waktuRegistrasi, tanggalRegistrasiStr, jamRegistrasiStr)
}

func (s *service) validasiStatusKamarDanBatasWaktu(ctx context.Context, noRawat string, waktuRegistrasi time.Time, tglRegStr, jamRegStr string) error {
	isAktifRanap, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat %s: %v", noRawat, err)
		return err
	}

	if hasRecordKamar {
		if !isAktifRanap {
			return apperror.NewBusinessError("Pasien rawat inap sudah keluar / checkout dari kamar inap")
		}
		return nil
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		errMsg := fmt.Sprintf("Batas waktu penilaian medis untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tglRegStr, jamRegStr)
		s.log.Warn("Penilaian medis ditolak karena lewat batas %d jam untuk no_rawat %s: %s", s.maxEditJam, noRawat, errMsg)
		return apperror.NewForbiddenError(errMsg)
	}

	return nil
}

func (s *service) validasiRegistrasiDanStatusRanap(ctx context.Context, noRawat, tglPenilaian string) error {
	tanggalRegistrasiStr, jamRegistrasiStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	waktuRegistrasi, err := shared.ParseWaktu(tanggalRegistrasiStr, jamRegistrasiStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tanggalRegistrasiStr, jamRegistrasiStr, err)
		return err
	}

	if tglPenilaian != "" {
		waktuPenilaian, errPer := time.ParseInLocation("2006-01-02 15:04:05", tglPenilaian, time.Local)
		if errPer == nil && waktuPenilaian.Before(waktuRegistrasi) {
			errs := apperror.ValidationError{
				"tanggal_penilaian": fmt.Sprintf("Waktu penilaian medis (%s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tglPenilaian, tanggalRegistrasiStr, jamRegistrasiStr),
			}
			s.log.Warn("Validasi waktu penilaian medis ranap gagal untuk no_rawat %s: %+v", noRawat, errs)
			return errs
		}
	}

	isAktifRanap, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat %s: %v", noRawat, err)
		return err
	}

	if !hasRecordKamar {
		return apperror.NewBusinessError("Pasien belum/tidak terdaftar di kamar inap. Penilaian medis ranap hanya untuk pasien rawat inap.")
	}

	if !isAktifRanap {
		return apperror.NewBusinessError("Pasien rawat inap sudah keluar / checkout dari kamar inap")
	}

	return nil
}
