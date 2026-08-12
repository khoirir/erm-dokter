package pemeriksaan

import (
	"context"
	"math"
	"strings"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.MetaPaginasi, error)
	DaftarPemeriksaanByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.MetaPaginasi, error)
	DetailPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) (*Pemeriksaan, error)
}

type service struct {
	repo           Repository
	rawatJalanRepo rawatjalan.Repository
	log            *logger.Logger
}

func NewService(repo Repository, rawatJalanRepo rawatjalan.Repository, log *logger.Logger) Service {
	return &service{
		repo:           repo,
		rawatJalanRepo: rawatJalanRepo,
		log:            log,
	}
}

func (s *service) DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.MetaPaginasi, error) {
	if noRawat == "" {
		return nil, shared.MetaPaginasi{}, apperror.NewBusinessError("nomor rawat tidak boleh kosong")
	}

	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		return nil, shared.MetaPaginasi{}, apperror.NewBusinessError("status lanjut tidak valid")
	}

	if errs := filter.Validate(); errs != nil {
		s.log.Warn("Filter validasi gagal: %+v", errs)
		return nil, shared.MetaPaginasi{}, errs
	}

	rawParts := strings.Split(noRawat, ",")
	var listNoRawat []string
	for _, p := range rawParts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			listNoRawat = append(listNoRawat, trimmed)
		}
	}
	if len(listNoRawat) == 0 {
		return nil, shared.MetaPaginasi{}, apperror.NewBusinessError("nomor rawat tidak boleh kosong")
	}

	daftarPemeriksaan, totalData, err := s.repo.DaftarPemeriksaan(ctx, listNoRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal query daftar pemeriksaan %s: %v", noRawat, err)
		return nil, shared.MetaPaginasi{}, err
	}

	totalHalaman := int(math.Ceil(float64(totalData) / float64(filter.Batas)))

	meta := shared.MetaPaginasi{
		TotalData:    totalData,
		TotalHalaman: totalHalaman,
		HalamanAktif: filter.Halaman,
		BatasData:    filter.Batas,
	}

	return daftarPemeriksaan, meta, nil
}

func (s *service) DaftarPemeriksaanByRM(ctx context.Context, noRekamMedis string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.MetaPaginasi, error) {
	if strings.TrimSpace(noRekamMedis) == "" {
		return nil, shared.MetaPaginasi{}, apperror.NewBusinessError("nomor rekam medis tidak boleh kosong")
	}

	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		return nil, shared.MetaPaginasi{}, apperror.NewBusinessError("status lanjut tidak valid")
	}

	if errs := filter.Validate(); errs != nil {
		s.log.Warn("Filter validasi gagal: %+v", errs)
		return nil, shared.MetaPaginasi{}, errs
	}

	riwayatKunjungan, err := s.rawatJalanRepo.RiwayatKunjunganPasien(ctx, noRekamMedis)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat kunjungan untuk RM %s: %v", noRekamMedis, err)
		return nil, shared.MetaPaginasi{}, err
	}

	if len(riwayatKunjungan) == 0 {
		return []Pemeriksaan{}, shared.MetaPaginasi{}, nil
	}

	listNoRawat := make([]string, len(riwayatKunjungan))
	for i, k := range riwayatKunjungan {
		listNoRawat[i] = k.NoRawat
	}

	daftarPemeriksaan, totalData, err := s.repo.DaftarPemeriksaan(ctx, listNoRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal query daftar pemeriksaan RM %s: %v", noRekamMedis, err)
		return nil, shared.MetaPaginasi{}, err
	}

	totalHalaman := int(math.Ceil(float64(totalData) / float64(filter.Batas)))

	meta := shared.MetaPaginasi{
		TotalData:    totalData,
		TotalHalaman: totalHalaman,
		HalamanAktif: filter.Halaman,
		BatasData:    filter.Batas,
	}

	return daftarPemeriksaan, meta, nil
}

func (s *service) DetailPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) (*Pemeriksaan, error) {
	if id.NoRawat == "" || id.TanggalPemeriksaan == "" || id.JamPemeriksaan == "" {
		return nil, apperror.NewBusinessError("ID detail pemeriksaan tidak lengkap")
	}

	if statusLanjut != "Semua" && !statusLanjut.IsValid() {
		return nil, apperror.NewBusinessError("status lanjut tidak valid")
	}

	pemeriksaan, err := s.repo.DetailPemeriksaan(ctx, id, statusLanjut)
	if err != nil {
		s.log.Error("Gagal query detail pemeriksaan %s (%s %s): %v", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, err)
		return nil, err
	}

	return pemeriksaan, nil
}
