package usecase

import (
	"context"
	"errors"
	"math"

	"erm-dokter/internal/domain"
	"erm-dokter/pkg/crypto"
)

type rawatJalanUsecase struct {
	rawatJalanRepo domain.RawatJalanRepository
	encryptionKey  string
}

func NewRawatJalanUsecase(repo domain.RawatJalanRepository, encryptionKey string) domain.RawatJalanUsecase {
	return &rawatJalanUsecase{
		rawatJalanRepo: repo,
		encryptionKey:  encryptionKey,
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

	for i := range daftarAntrean {
		encrypted, err := crypto.Encrypt(daftarAntrean[i].NoRawat, u.encryptionKey)
		if err == nil {
			daftarAntrean[i].DetailURL = "/api/v1/rawat-jalan/detail/" + encrypted
		}
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

func (u *rawatJalanUsecase) DetailKunjungan(ctx context.Context, encryptedNoRawat string, kodeDokter string) (*domain.KunjunganRawatJalan, error) {
	if encryptedNoRawat == "" {
		return nil, errors.New("nomor rawat tidak boleh kosong")
	}

	noRawat, err := crypto.Decrypt(encryptedNoRawat, u.encryptionKey)
	if err != nil {
		return nil, errors.New("nomor rawat tidak valid")
	}

	return u.rawatJalanRepo.DetailKunjungan(ctx, noRawat, kodeDokter)
}

func (u *rawatJalanUsecase) GetReferensiFilter(ctx context.Context) domain.ReferensiFilterRawatJalan {
	return domain.ReferensiFilterRawatJalan{
		StatusPemeriksaan: []domain.OpsiReferensi{
			{Value: string(domain.StatusBelum), Label: "Belum Periksa"},
			{Value: string(domain.StatusSudah), Label: "Sudah Periksa"},
			{Value: string(domain.StatusBatal), Label: "Batal Periksa"},
			{Value: string(domain.StatusBerkasDiterima), Label: "Berkas Diterima"},
			{Value: string(domain.StatusDirujuk), Label: "Dirujuk"},
			{Value: string(domain.StatusMeninggal), Label: "Meninggal"},
			{Value: string(domain.StatusDirawat), Label: "Dirawat"},
			{Value: string(domain.StatusPulangPaksa), Label: "Pulang Paksa"},
		},
		StatusLanjut: []domain.OpsiReferensi{
			{Value: string(domain.StatusLanjutRawatJalan), Label: "Rawat Jalan"},
			{Value: string(domain.StatusLanjutRawatInap), Label: "Rawat Inap"},
		},
		StatusBayar: []domain.OpsiReferensi{
			{Value: string(domain.StatusBayarSudah), Label: "Sudah Bayar"},
			{Value: string(domain.StatusBayarBelum), Label: "Belum Bayar"},
		},
		JenisAntrean: []domain.OpsiReferensi{
			{Value: string(domain.JenisAntreanRujukan), Label: "Rujukan"},
			{Value: string(domain.JenisAntreanTidakRujukan), Label: "Bukan Rujukan"},
		},
	}
}
