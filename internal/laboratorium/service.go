package laboratorium

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/tindakan"
)

type Service interface {
	GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) (*HasilLaboratoriumKunjungan, shared.PaginationMeta, error)
	GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, shared.PaginationMeta, error)
	GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*HasilLaboratorium, error)

	SimpanPermintaanLabPK(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPKRequest) (*DetailPermintaanLabPK, error)
	GetDaftarPermintaanLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPK, error)
	GetRiwayatPermintaanLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPK, shared.PaginationMeta, error)
	GetDetailPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanLabPK, error)
	HapusPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error

	SimpanPermintaanLabPA(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPARequest) (*DetailPermintaanLabPA, error)
	GetDaftarPermintaanLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPA, error)
	GetRiwayatPermintaanLabPAByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPA, shared.PaginationMeta, error)
	GetDetailPermintaanLabPA(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanLabPA, error)
	HapusPermintaanLabPA(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error
}

type service struct {
	repo              Repository
	berkasSvc         berkasdigital.Service
	rawatJalanService rawatjalan.Service
	tindakanService   tindakan.Service
	maxEditJam        int
	urlBerkasDigital  string
	kodeBerkasPK      []string
	kodeBerkasPA      []string
	kodeBerkasMB      []string
	log               *logger.Logger
}

func NewService(
	repo Repository,
	berkasSvc berkasdigital.Service,
	rawatJalanService rawatjalan.Service,
	tindakanService tindakan.Service,
	maxEditJam int,
	urlBerkasDigital string,
	kodeBerkasPK, kodeBerkasPA, kodeBerkasMB []string,
	log *logger.Logger,
) Service {
	baseURL := ""
	if urlBerkasDigital != "" {
		baseURL = strings.TrimRight(urlBerkasDigital, "/") + "/"
	}
	return &service{
		repo:              repo,
		berkasSvc:         berkasSvc,
		rawatJalanService: rawatJalanService,
		tindakanService:   tindakanService,
		maxEditJam:        maxEditJam,
		urlBerkasDigital:  baseURL,
		kodeBerkasPK:      kodeBerkasPK,
		kodeBerkasPA:      kodeBerkasPA,
		kodeBerkasMB:      kodeBerkasMB,
		log:               log,
	}
}

func (s *service) GetRiwayatLabKunjungan(ctx context.Context, kategori shared.KategoriLab, noRawat string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) (*HasilLaboratoriumKunjungan, shared.PaginationMeta, error) {
	var items []HasilLaboratorium
	var total int
	var kodeBerkas []string
	var err error

	switch kategori {
	case shared.KategoriLabPK:
		items, total, err = s.repo.DaftarHasilLabPK(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasPK
	case shared.KategoriLabPA:
		items, total, err = s.repo.DaftarHasilLabPA(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasPA
	case shared.KategoriLabMB:
		items, total, err = s.repo.DaftarHasilLabMB(ctx, noRawat, statusLanjut, filter)
		kodeBerkas = s.kodeBerkasMB
	default:
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		s.log.Error("gagal mengambil riwayat lab kunjungan %s kategori %s: %v", noRawat, kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	berkasList := s.fetchBerkasDigitalKunjungan(ctx, noRawat, kodeBerkas)

	kunjunganData := &HasilLaboratoriumKunjungan{
		HasilPemeriksaan: items,
		BerkasDigital:    berkasList,
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return kunjunganData, meta, nil
}

func (s *service) GetRiwayatLabPasien(ctx context.Context, kategori shared.KategoriLab, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]HasilLaboratorium, shared.PaginationMeta, error) {
	var items []HasilLaboratorium
	var total int
	var err error

	switch kategori {
	case shared.KategoriLabPK:
		items, total, err = s.repo.DaftarHasilLabPKByRM(ctx, noRM, statusLanjut, filter)
	case shared.KategoriLabPA:
		items, total, err = s.repo.DaftarHasilLabPAByRM(ctx, noRM, statusLanjut, filter)
	case shared.KategoriLabMB:
		items, total, err = s.repo.DaftarHasilLabMBByRM(ctx, noRM, statusLanjut, filter)
	default:
		return nil, shared.PaginationMeta{}, apperror.NewBusinessError("Kategori laboratorium tidak valid")
	}

	if err != nil {
		s.log.Error("gagal mengambil riwayat lab pasien %s kategori %s: %v", noRM, kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return items, meta, nil
}

func (s *service) GetDetailHasilLab(ctx context.Context, kategori shared.KategoriLab, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa string) (*HasilLaboratorium, error) {
	var item *HasilLaboratorium
	var err error

	switch kategori {
	case shared.KategoriLabPK:
		item, err = s.repo.DetailHasilLabPK(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	case shared.KategoriLabPA:
		item, err = s.repo.DetailHasilLabPA(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
	case shared.KategoriLabMB:
		item, err = s.repo.DetailHasilLabMB(ctx, noRawat, kodeTindakan, tanggalPeriksa, jamPeriksa)
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

	return item, nil
}

func (s *service) fetchBerkasDigitalKunjungan(ctx context.Context, noRawat string, kodeBerkas []string) []berkasdigital.BerkasDigital {
	if len(kodeBerkas) == 0 {
		return make([]berkasdigital.BerkasDigital, 0)
	}

	berkasList, err := s.berkasSvc.GetBerkasByNoRawat(ctx, noRawat, kodeBerkas)
	if err != nil || len(berkasList) == 0 {
		return make([]berkasdigital.BerkasDigital, 0)
	}

	berkasDigital := make([]berkasdigital.BerkasDigital, 0, len(berkasList))
	for _, b := range berkasList {
		if b.LokasiFile == "" {
			continue
		}
		fullURL := s.berkasSvc.BuildFullURL(b.LokasiFile)
		berkasDigital = append(berkasDigital, berkasdigital.BerkasDigital{
			Kode:       b.Kode,
			NamaBerkas: b.NamaBerkas,
			IdBerkas:   fullURL,
		})
	}
	return berkasDigital
}
