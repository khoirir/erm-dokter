package laboratorium

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) (*HasilLaboratoriumKunjungan, shared.PaginationMeta, error)
	GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, encryptedIdPasien string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, shared.PaginationMeta, error)
	GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, encryptedIdHasil string) (*HasilLaboratorium, error)
}

type service struct {
	repo             Repository
	berkasSvc        berkasdigital.Service
	jwtSecret        string
	urlBerkasDigital string
	kodeBerkasPK     []string
	kodeBerkasPA     []string
	kodeBerkasMB     []string
	log              *logger.Logger
}

func NewService(repo Repository, berkasSvc berkasdigital.Service, jwtSecret string, urlBerkasDigital string, kodeBerkasPK, kodeBerkasPA, kodeBerkasMB []string, log *logger.Logger) Service {
	baseURL := ""
	if urlBerkasDigital != "" {
		baseURL = strings.TrimRight(urlBerkasDigital, "/") + "/"
	}
	return &service{
		repo:             repo,
		berkasSvc:        berkasSvc,
		jwtSecret:        jwtSecret,
		urlBerkasDigital: baseURL,
		kodeBerkasPK:     kodeBerkasPK,
		kodeBerkasPA:     kodeBerkasPA,
		kodeBerkasMB:     kodeBerkasMB,
		log:              log,
	}
}

func (s *service) GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) (*HasilLaboratoriumKunjungan, shared.PaginationMeta, error) {
	noRawat, err := crypto.Decrypt(encryptedIdKunjungan, s.jwtSecret)
	if err != nil {
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("ID kunjungan tidak valid")
	}

	var itemsDB []HasilLaboratoriumDB
	var total int
	var kodeBerkas []string

	switch kategori {
	case shared.KategoriLabPK:
		itemsDB, total, err = s.repo.DaftarHasilLabPK(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasPK
	case shared.KategoriLabPA:
		itemsDB, total, err = s.repo.DaftarHasilLabPA(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasPA
	case shared.KategoriLabMB:
		itemsDB, total, err = s.repo.DaftarHasilLabMB(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasMB
	default:
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		s.log.Error("gagal mengambil riwayat lab kunjungan %s kategori %s: %v", noRawat, kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	hasilPemeriksaan, err := s.mapListHasilLab(itemsDB)
	if err != nil {
		return nil, shared.PaginationMeta{}, err
	}

	berkasList := s.fetchBerkasDigitalKunjungan(ctx, noRawat, kodeBerkas)

	kunjunganData := &HasilLaboratoriumKunjungan{
		HasilPemeriksaan: hasilPemeriksaan,
		BerkasDigital:    berkasList,
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return kunjunganData, meta, nil
}

func (s *service) GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, encryptedIdPasien string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, shared.PaginationMeta, error) {
	noRM, err := crypto.Decrypt(encryptedIdPasien, s.jwtSecret)
	if err != nil {
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("ID pasien tidak valid")
	}

	var itemsDB []HasilLaboratoriumDB
	var total int

	switch kategori {
	case shared.KategoriLabPK:
		itemsDB, total, err = s.repo.DaftarHasilLabPKByRM(ctx, noRM, statusLanjut, filter)
	case shared.KategoriLabPA:
		itemsDB, total, err = s.repo.DaftarHasilLabPAByRM(ctx, noRM, statusLanjut, filter)
	case shared.KategoriLabMB:
		itemsDB, total, err = s.repo.DaftarHasilLabMBByRM(ctx, noRM, statusLanjut, filter)
	default:
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		s.log.Error("gagal mengambil riwayat lab pasien %s kategori %s: %v", noRM, kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	result, err := s.mapListHasilLab(itemsDB)
	if err != nil {
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return result, meta, nil
}

func (s *service) GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, encryptedIdKunjungan string, encryptedIdHasil string) (*HasilLaboratorium, error) {
	expectedNoRawat, err := crypto.Decrypt(encryptedIdKunjungan, s.jwtSecret)
	if err != nil {
		return nil, apperror.NewBusinessError("ID kunjungan tidak valid")
	}

	decryptedHasil, err := crypto.Decrypt(encryptedIdHasil, s.jwtSecret)
	if err != nil {
		return nil, apperror.NewBusinessError("ID hasil lab tidak valid")
	}

	parts := strings.Split(decryptedHasil, "~")
	if len(parts) != 4 {
		return nil, apperror.NewBusinessError("Format ID hasil lab tidak valid")
	}

	noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa := parts[0], parts[1], parts[2], parts[3]
	if noRawat != expectedNoRawat {
		return nil, apperror.NewBusinessError("ID hasil lab tidak sesuai dengan kunjungan pasien")
	}

	var itemDB *HasilLaboratoriumDB
	switch kategori {
	case shared.KategoriLabPK:
		itemDB, err = s.repo.DetailHasilLabPK(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	case shared.KategoriLabPA:
		itemDB, err = s.repo.DetailHasilLabPA(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	case shared.KategoriLabMB:
		itemDB, err = s.repo.DetailHasilLabMB(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	default:
		return nil, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data hasil laboratorium tidak ditemukan")
		}
		s.log.Error("gagal mengambil detail hasil lab %s %s: %v", noRawat, kodeTindakan, err)
		return nil, err
	}

	return s.mapSingleHasilLab(*itemDB)
}

func (s *service) fetchBerkasDigitalKunjungan(ctx context.Context, noRawat string, kodeBerkas []string) []BerkasDigital {
	if len(kodeBerkas) == 0 {
		return []BerkasDigital{}
	}

	berkasDBList, err := s.berkasSvc.GetBerkasByNoRawat(ctx, noRawat, kodeBerkas)
	if err != nil || len(berkasDBList) == 0 {
		return []BerkasDigital{}
	}

	berkasDigital := make([]BerkasDigital, 0, len(berkasDBList))
	for _, b := range berkasDBList {
		if b.LokasiFile != "" {
			fullURL := s.buildFullURLBerkas(b.LokasiFile)
			encIdBerkas, errEnc := crypto.Encrypt(fullURL, s.jwtSecret)
			if errEnc == nil {
				berkasDigital = append(berkasDigital, BerkasDigital{
					Kode:       b.Kode,
					NamaBerkas: b.NamaBerkas,
					IdBerkas:   encIdBerkas,
					UrlBerkas:  "/api/v1/berkas-digital/" + encIdBerkas,
				})
			}
		}
	}
	return berkasDigital
}

func (s *service) mapListHasilLab(items []HasilLaboratoriumDB) ([]HasilLaboratorium, error) {
	result := make([]HasilLaboratorium, 0, len(items))
	for _, item := range items {
		mapped, err := s.mapSingleHasilLab(item)
		if err != nil {
			return nil, err
		}
		result = append(result, *mapped)
	}
	return result, nil
}

func (s *service) mapSingleHasilLab(item HasilLaboratoriumDB) (*HasilLaboratorium, error) {
	encryptedIdHasil, err := crypto.Encrypt(item.NoRawat+"~"+item.KodeTindakan+"~"+item.TanggalPeriksa+"~"+item.JamPeriksa, s.jwtSecret)
	if err != nil {
		s.log.Error("gagal mengenkripsi id hasil lab: %v", err)
		return nil, err
	}

	encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, s.jwtSecret)
	if err != nil {
		s.log.Error("gagal mengenkripsi id kunjungan: %v", err)
		return nil, err
	}

	var detailPK []ItemHasilLabPK
	if item.DetailPK != nil {
		detailPK = make([]ItemHasilLabPK, 0, len(item.DetailPK))
		for _, d := range item.DetailPK {
			var encIdTemplate string
			if d.IdTemplate > 0 {
				encIdTemplate, _ = crypto.Encrypt(strconv.Itoa(d.IdTemplate), s.jwtSecret)
			}
			detailPK = append(detailPK, ItemHasilLabPK{
				IdTemplate:      encIdTemplate,
				NamaPemeriksaan: d.NamaPemeriksaan,
				Nilai:           d.Nilai,
				Satuan:          d.Satuan,
				NilaiRujukan:    d.NilaiRujukan,
				Keterangan:      d.Keterangan,
			})
		}
	}

	var detailPA *HasilLabPA
	if item.DetailPA != nil {
		detailPA = &HasilLabPA{
			DiagnosaKlinik: item.DetailPA.DiagnosaKlinik,
			Makroskopis:    item.DetailPA.Makroskopis,
			Mikroskopis:    item.DetailPA.Mikroskopis,
			Kesimpulan:     item.DetailPA.Kesimpulan,
			Kesan:          item.DetailPA.Kesan,
		}
	}

	return &HasilLaboratorium{
		Id:             encryptedIdHasil,
		IdKunjungan:    encryptedIdKunjungan,
		NoRawat:        item.NoRawat,
		KodeTindakan:   item.KodeTindakan,
		NamaTindakan:   item.NamaTindakan,
		Kategori:       item.Kategori,
		Status:         item.Status,
		TanggalPeriksa: item.TanggalPeriksa,
		JamPeriksa:     item.JamPeriksa,
		DokterPerujuk: DokterInfo{
			KodeDokter: item.KodeDokterPerujuk,
			NamaDokter: item.NamaDokterPerujuk,
		},
		DokterPJ: DokterInfo{
			KodeDokter: item.KodeDokterPJ,
			NamaDokter: item.NamaDokterPJ,
		},
		Petugas: PetugasInfo{
			Nip:  item.NipPetugas,
			Nama: item.NamaPetugas,
		},
		DetailPK: detailPK,
		DetailPA: detailPA,
	}, nil
}

func (s *service) buildFullURLBerkas(lokasiFile string) string {
	if strings.HasPrefix(lokasiFile, "http://") || strings.HasPrefix(lokasiFile, "https://") {
		return lokasiFile
	}
	cleanPath := strings.TrimPrefix(lokasiFile, "/")
	return s.urlBerkasDigital + cleanPath
}
