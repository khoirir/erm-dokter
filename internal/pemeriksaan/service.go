package pemeriksaan

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/go-sql-driver/mysql"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarPemeriksaan(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.MetaPaginasi, error)
	DaftarPemeriksaanByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterDaftarPemeriksaan) ([]Pemeriksaan, shared.MetaPaginasi, error)
	DetailPemeriksaan(ctx context.Context, id IdPemeriksaan, statusLanjut shared.StatusLanjut) (*Pemeriksaan, error)
	GetDaftarKesadaran(ctx context.Context) []OpsiReferensi
	SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPemeriksaanRequest) (*Pemeriksaan, error)
	HapusPemeriksaan(ctx context.Context, kodeDokter string, id IdPemeriksaan, statusLanjut shared.StatusLanjut) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, maxEditJam int, log *logger.Logger) Service {
	if maxEditJam <= 0 {
		maxEditJam = 48
	}
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		maxEditJam:        maxEditJam,
		log:               log,
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

	riwayatKunjungan, err := s.rawatJalanService.RiwayatKunjunganPasien(ctx, noRekamMedis)
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
		s.log.Error("Gagal query daftar pemeriksaan by RM %s: %v", noRekamMedis, err)
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
	if statusLanjut != shared.StatusLanjutRawatJalan && statusLanjut != shared.StatusLanjutRawatInap {
		return nil, apperror.NewBusinessError("status lanjut tidak valid (harus Ralan atau Ranap)")
	}

	if id.NoRawat == "" || id.TanggalPemeriksaan == "" || id.JamPemeriksaan == "" {
		return nil, apperror.NewBusinessError("parameter ID pemeriksaan tidak lengkap")
	}

	pemeriksaan, err := s.repo.DetailPemeriksaan(ctx, id, statusLanjut)
	if err != nil {
		s.log.Error("Gagal query detail pemeriksaan %+v (%s): %v", id, statusLanjut, err)
		return nil, err
	}

	return pemeriksaan, nil
}

func (s *service) GetDaftarKesadaran(ctx context.Context) []OpsiReferensi {
	return []OpsiReferensi{
		{Value: string(KesadaranComposMentis), Label: "Compos Mentis"},
		{Value: string(KesadaranSomnolen), Label: "Somnolen"},
		{Value: string(KesadaranSopor), Label: "Sopor"},
		{Value: string(KesadaranKoma), Label: "Koma"},
		{Value: string(KesadaranAlert), Label: "Alert"},
		{Value: string(KesadaranConfusion), Label: "Confusion"},
		{Value: string(KesadaranVoice), Label: "Voice"},
		{Value: string(KesadaranPain), Label: "Pain"},
		{Value: string(KesadaranUnresponsive), Label: "Unresponsive"},
		{Value: string(KesadaranApatis), Label: "Apatis"},
		{Value: string(KesadaranDelirium), Label: "Delirium"},
	}
}

func (s *service) SimpanPemeriksaan(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanPemeriksaanRequest) (*Pemeriksaan, error) {
	if statusLanjut != shared.StatusLanjutRawatJalan && statusLanjut != shared.StatusLanjutRawatInap {
		return nil, apperror.NewBusinessError("status lanjut tidak valid (harus Ralan atau Ranap)")
	}

	if errs := req.Validate(); errs != nil {
		s.log.Warn("Validasi simpan pemeriksaan gagal: %+v", errs)
		return nil, errs
	}

	tglRegStr, jamRegStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, req.NoRawat)
	if err != nil {
		s.log.Error("Gagal mengambil data registrasi no_rawat %s: %v", req.NoRawat, err)
		return nil, err
	}
	if !exists {
		return nil, apperror.NewNotFoundError("Data registrasi kunjungan pasien tidak ditemukan")
	}

	waktuRegistrasi, err := shared.ParseWaktu(tglRegStr, jamRegStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", req.NoRawat, tglRegStr, jamRegStr, err)
		return nil, err
	}

	waktuPemeriksaan, err := shared.ParseWaktu(req.TanggalPemeriksaan, req.JamPemeriksaan)
	if err != nil {
		return nil, apperror.NewBusinessError(err.Error())
	}

	if waktuPemeriksaan.Before(waktuRegistrasi) {
		errs := apperror.ValidationError{
			"tanggal_pemeriksaan": fmt.Sprintf("Waktu pemeriksaan (%s %s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", req.TanggalPemeriksaan, req.JamPemeriksaan, tglRegStr, jamRegStr),
		}
		s.log.Warn("Validasi waktu pemeriksaan gagal untuk no_rawat %s: %+v", req.NoRawat, errs)
		return nil, errs
	}

	if err := s.repo.SimpanPemeriksaan(ctx, kodeDokter, statusLanjut, req); err != nil {
		var mysqlErr *mysql.MySQLError
		if (errors.As(err, &mysqlErr) && mysqlErr.Number == 1062) || strings.Contains(err.Error(), "1062") || strings.Contains(err.Error(), "Duplicate entry") {
			errMsg := fmt.Sprintf("Data pemeriksaan pada tanggal %s jam %s sudah pernah disimpan sebelumnya. Silakan sesuaikan jam pemeriksaan.", req.TanggalPemeriksaan, req.JamPemeriksaan)
			s.log.Warn("Penyimpanan pemeriksaan duplikat untuk no_rawat %s (%s %s): %s", req.NoRawat, req.TanggalPemeriksaan, req.JamPemeriksaan, errMsg)
			return nil, apperror.ValidationError{
				"jam_pemeriksaan": errMsg,
			}
		}
		s.log.Error("Gagal menyimpan pemeriksaan no_rawat %s (%s): %v", req.NoRawat, statusLanjut, err)
		return nil, err
	}

	idPemeriksaan := IdPemeriksaan{
		NoRawat:            req.NoRawat,
		TanggalPemeriksaan: req.TanggalPemeriksaan,
		JamPemeriksaan:     req.JamPemeriksaan,
	}

	detail, err := s.repo.DetailPemeriksaan(ctx, idPemeriksaan, statusLanjut)
	if err != nil || detail == nil {
		detail = &Pemeriksaan{
			NoRawat:             req.NoRawat,
			TanggalPemeriksaan:  req.TanggalPemeriksaan,
			JamPemeriksaan:      req.JamPemeriksaan,
			SuhuTubuh:           req.SuhuTubuh,
			Tensi:               req.Tensi,
			Nadi:                req.Nadi,
			Respirasi:           req.Respirasi,
			TinggiBadan:         req.TinggiBadan,
			BeratBadan:          req.BeratBadan,
			SpO2:                req.SpO2,
			Gcs:                 req.Gcs,
			Kesadaran:           req.Kesadaran,
			Keluhan:             req.Keluhan,
			Pemeriksaan:         req.Pemeriksaan,
			Alergi:              req.Alergi,
			LingkarPerut:        req.LingkarPerut,
			RencanaTindakLanjut: req.RencanaTindakLanjut,
			Penilaian:           req.Penilaian,
			Instruksi:           req.Instruksi,
			Evaluasi:            req.Evaluasi,
			KodeDokterPetugas:   kodeDokter,
		}
	}

	s.log.Info("Berhasil menyimpan pemeriksaan no_rawat %s (%s) oleh dokter %s", req.NoRawat, statusLanjut, kodeDokter)
	return detail, nil
}

func (s *service) HapusPemeriksaan(ctx context.Context, kodeDokter string, id IdPemeriksaan, statusLanjut shared.StatusLanjut) error {
	if statusLanjut != shared.StatusLanjutRawatJalan && statusLanjut != shared.StatusLanjutRawatInap {
		return apperror.NewBusinessError("status lanjut tidak valid (harus Ralan atau Ranap)")
	}

	if id.NoRawat == "" || id.TanggalPemeriksaan == "" || id.JamPemeriksaan == "" {
		return apperror.NewBusinessError("parameter ID pemeriksaan tidak lengkap")
	}

	pemeriksaan, err := s.repo.DetailPemeriksaan(ctx, id, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil detail pemeriksaan untuk hapus %+v (%s): %v", id, statusLanjut, err)
		return err
	}
	if pemeriksaan == nil {
		return apperror.NewNotFoundError("Data pemeriksaan tidak ditemukan")
	}

	if pemeriksaan.KodeDokterPetugas != kodeDokter {
		s.log.Warn("Percobaan menghapus pemeriksaan no_rawat %s (%s %s) oleh dokter %s ditolak: diinput oleh %s (%s)", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, kodeDokter, pemeriksaan.KodeDokterPetugas, pemeriksaan.NamaDokterPetugas)
		return apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk menghapus data pemeriksaan ini karena diinput oleh dokter/petugas lain (%s)", pemeriksaan.NamaDokterPetugas))
	}

	if err := shared.ValidasiBatasWaktuRekamMedis(pemeriksaan.TanggalPemeriksaan, pemeriksaan.JamPemeriksaan, s.maxEditJam, "dihapus"); err != nil {
		s.log.Warn("Percobaan menghapus pemeriksaan no_rawat %s (%s %s) oleh dokter %s ditolak: %v", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, kodeDokter, err)
		return err
	}

	if err := s.repo.HapusPemeriksaan(ctx, id, statusLanjut); err != nil {
		s.log.Error("Gagal menghapus pemeriksaan no_rawat %s (%s): %v", id.NoRawat, statusLanjut, err)
		return err
	}

	s.log.Info("Berhasil menghapus data pemeriksaan no_rawat %s (%s %s) oleh dokter %s", id.NoRawat, id.TanggalPemeriksaan, id.JamPemeriksaan, kodeDokter)
	return nil
}



