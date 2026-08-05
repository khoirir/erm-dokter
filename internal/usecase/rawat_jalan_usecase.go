package usecase

import (
	"context"
	"errors"
	"math"

	"erm-dokter/internal/domain"
)

type rawatJalanUsecase struct {
	rawatJalanRepo domain.RawatJalanRepository
}

func NewRawatJalanUsecase(repo domain.RawatJalanRepository) domain.RawatJalanUsecase {
	return &rawatJalanUsecase{
		rawatJalanRepo: repo,
	}
}

type ResponAntreanDokter struct {
	Data []domain.KunjunganRawatJalan `json:"data"`
	Meta domain.MetaPaginasi          `json:"meta"`
}

func (u *rawatJalanUsecase) DaftarAntreanDokter(ctx context.Context, filter domain.FilterAntreanDokter) ([]domain.KunjunganRawatJalan, domain.MetaPaginasi, error) {
	if len(filter.KataKunci) > 0 && len(filter.KataKunci) < 3 {
		return nil, domain.MetaPaginasi{}, errors.New("kata kunci pencarian minimal 3 karakter")
	}

	if filter.Halaman <= 0 {
		filter.Halaman = 1
	}
	if filter.Batas <= 0 {
		filter.Batas = 10
	}

	daftarAntrean, totalData, err := u.rawatJalanRepo.DaftarAntreanDokter(ctx, filter)
	if err != nil {
		return nil, domain.MetaPaginasi{}, err
	}

	totalHalaman := int(math.Ceil(float64(totalData) / float64(filter.Batas)))

	meta := domain.MetaPaginasi{
		TotalData:    totalData,
		TotalHalaman: totalHalaman,
		HalamanAktif: filter.Halaman,
		BatasData:    filter.Batas,
	}

	return daftarAntrean, meta, nil
}

func (u *rawatJalanUsecase) DetailKunjungan(ctx context.Context, noRawat string) (*domain.KunjunganRawatJalan, error) {
	if noRawat == "" {
		return nil, errors.New("nomor rawat tidak boleh kosong")
	}

	return u.rawatJalanRepo.DetailKunjungan(ctx, noRawat)
}
