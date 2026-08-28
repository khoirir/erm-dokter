package resep

import (
	"context"
	"fmt"
	"time"

	"erm-dokter/internal/obat"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarResep(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error)
	DaftarResepByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error)
	DetailResep(ctx context.Context, noResep string) (*Resep, error)
	SimpanResep(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanResepRequest) (*Resep, error)
	HapusResep(ctx context.Context, kodeDokter, noRawat, noResep string, statusLanjut shared.StatusLanjut) error
	DaftarAturanPakai(ctx context.Context, keyword string) ([]AturanPakai, error)
	DaftarMetodeRacik(ctx context.Context) ([]MetodeRacik, error)
}


type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	obatService       obat.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, obatService obat.Service, maxEditJam int, log *logger.Logger) Service {
	if maxEditJam <= 0 {
		maxEditJam = 48
	}
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		obatService:       obatService,
		maxEditJam:        maxEditJam,
		log:               log,
	}
}

func (s *service) DaftarResep(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error) {
	data, total, err := s.repo.DaftarResep(ctx, noRawat, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil daftar resep untuk no_rawat %s: %v", noRawat, err)
		return nil, shared.PaginationMeta{}, err
	}

	return data, shared.NewPaginationMeta(total, filter.Page, filter.Limit), nil
}

func (s *service) DaftarResepByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterDaftarResep) ([]Resep, shared.PaginationMeta, error) {
	data, total, err := s.repo.DaftarResepByRM(ctx, noRM, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat resep untuk no_rm %s: %v", noRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	return data, shared.NewPaginationMeta(total, filter.Page, filter.Limit), nil
}

func (s *service) DetailResep(ctx context.Context, noResep string) (*Resep, error) {
	resep, err := s.repo.DetailResep(ctx, noResep)
	if err != nil {
		s.log.Error("Gagal mengambil detail resep %s: %v", noResep, err)
		return nil, err
	}
	if resep == nil {
		return nil, apperror.NewNotFoundError("Data resep obat tidak ditemukan")
	}

	return resep, nil
}

func (s *service) SimpanResep(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanResepRequest) (*Resep, error) {
	if err := s.validasiWaktuRegistrasi(ctx, req.NoRawat, req.TanggalPeresepan, req.JamPeresepan, statusLanjut); err != nil {
		return nil, err
	}

	if err := s.validasiStatusKamarInap(ctx, req.NoRawat, statusLanjut); err != nil {
		return nil, err
	}

	if err := s.validasiKeberadaanObat(ctx, req); err != nil {
		return nil, err
	}

	if err := s.validasiMetodeRacik(ctx, req); err != nil {
		return nil, err
	}

	resep, err := s.repo.SimpanResep(ctx, kodeDokter, statusLanjut, req)
	if err != nil {
		s.log.Error("Gagal menyimpan resep obat no_rawat %s (%s): %v", req.NoRawat, statusLanjut, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan resep obat no_resep %s untuk no_rawat %s oleh dokter %s", resep.NoResep, req.NoRawat, kodeDokter)
	return resep, nil
}

func (s *service) HapusResep(ctx context.Context, kodeDokter, noRawat, noResep string, statusLanjut shared.StatusLanjut) error {
	resep, err := s.repo.DetailResep(ctx, noResep)
	if err != nil {
		s.log.Error("Gagal mengambil detail resep %s untuk hapus: %v", noResep, err)
		return err
	}
	if resep == nil || resep.NoRawat != noRawat {
		return apperror.NewNotFoundError("Data resep obat tidak ditemukan pada kunjungan ini")
	}

	if resep.KodeDokter != kodeDokter {
		s.log.Warn("Percobaan menghapus resep no_resep %s oleh dokter %s ditolak: diresepkan oleh %s (%s)", noResep, kodeDokter, resep.KodeDokter, resep.NamaDokter)
		return apperror.NewForbiddenError(fmt.Sprintf("Anda tidak memiliki hak akses untuk menghapus resep ini karena diresepkan oleh dokter lain (%s)", resep.NamaDokter))
	}

	if resep.TanggalPerawatan != "" && resep.TanggalPerawatan != "0000-00-00" && resep.JamPerawatan != "" && resep.JamPerawatan != "00:00:00" {
		s.log.Warn("Percobaan menghapus resep no_resep %s ditolak: telah divalidasi farmasi pada %s %s", noResep, resep.TanggalPerawatan, resep.JamPerawatan)
		return apperror.NewForbiddenError("Resep obat telah divalidasi oleh pihak farmasi dan tidak dapat dihapus")
	}

	if resep.TanggalPenyerahan != "" && resep.TanggalPenyerahan != "0000-00-00" && resep.JamPenyerahan != "" && resep.JamPenyerahan != "00:00:00" {
		s.log.Warn("Percobaan menghapus resep no_resep %s ditolak: telah diserahkan ke pasien pada %s %s", noResep, resep.TanggalPenyerahan, resep.JamPenyerahan)
		return apperror.NewForbiddenError("Resep obat telah diserahkan ke pasien dan tidak dapat dihapus")
	}

	if statusLanjut == shared.StatusLanjutRawatJalan {
		tglRegStr, jamRegStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
		if err != nil {
			s.log.Error("Gagal mengambil data registrasi no_rawat %s: %v", noRawat, err)
			return err
		}
		if !exists {
			return apperror.NewNotFoundError("Data registrasi kunjungan pasien tidak ditemukan")
		}

		waktuRegistrasi, err := shared.ParseWaktu(tglRegStr, jamRegStr)
		if err != nil {
			s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tglRegStr, jamRegStr, err)
			return err
		}

		batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
		if time.Now().After(batasWaktu) {
			errMsg := fmt.Sprintf("Batas waktu penghapusan resep obat untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tglRegStr, jamRegStr)
			s.log.Warn("Batas waktu hapus resep kadaluarsa untuk no_rawat %s: %s", noRawat, errMsg)
			return apperror.NewForbiddenError(errMsg)
		}
	}

	if err := s.repo.HapusResep(ctx, noResep); err != nil {
		s.log.Error("Gagal menghapus resep obat no_resep %s untuk no_rawat %s: %v", noResep, noRawat, err)
		return err
	}

	s.log.Info("Berhasil menghapus resep obat no_resep %s untuk no_rawat %s oleh dokter %s", noResep, noRawat, kodeDokter)
	return nil
}

func (s *service) validasiKeberadaanObat(ctx context.Context, req SimpanResepRequest) error {

	listKode := make([]string, 0)
	for _, rd := range req.ResepDokter {
		if rd.KodeObat != "" {
			listKode = append(listKode, rd.KodeObat)
		}
	}
	for _, rr := range req.ResepRacikan {
		for _, d := range rr.Detail {
			if d.KodeObat != "" {
				listKode = append(listKode, d.KodeObat)
			}
		}
	}

	if len(listKode) == 0 {
		return nil
	}

	foundMap, err := s.obatService.CekKeberadaanObat(ctx, listKode)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan obat: %v", err)
		return err
	}

	valErrs := make(apperror.ValidationError)
	for i, rd := range req.ResepDokter {
		if !foundMap[rd.KodeObat] {
			valErrs[fmt.Sprintf("resep_dokter[%d].id_obat", i)] = "Data obat tidak ditemukan"
		}
	}

	for i, rr := range req.ResepRacikan {
		for j, d := range rr.Detail {
			if !foundMap[d.KodeObat] {
				valErrs[fmt.Sprintf("resep_racikan[%d].detail[%d].id_obat", i, j)] = fmt.Sprintf("Bahan racikan ke-%d: Data obat tidak ditemukan", j+1)
			}
		}
	}

	if len(valErrs) > 0 {
		s.log.Warn("Validasi keberadaan obat gagal untuk no_rawat %s: %+v", req.NoRawat, valErrs)
		return valErrs
	}

	return nil
}

func (s *service) validasiMetodeRacik(ctx context.Context, req SimpanResepRequest) error {
	listKodeRacik := make([]string, 0)
	for _, rr := range req.ResepRacikan {
		if rr.KodeRacik != "" {
			listKodeRacik = append(listKodeRacik, rr.KodeRacik)
		}
	}

	if len(listKodeRacik) == 0 {
		return nil
	}

	foundMap, err := s.repo.CekKeberadaanMetodeRacik(ctx, listKodeRacik)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan metode racik: %v", err)
		return err
	}

	valErrs := make(apperror.ValidationError)
	for i, rr := range req.ResepRacikan {
		if !foundMap[rr.KodeRacik] {
			valErrs[fmt.Sprintf("resep_racikan[%d].kode_racik", i)] = "Metode racik tidak ditemukan"
		}
	}

	if len(valErrs) > 0 {
		s.log.Warn("Validasi metode racik gagal untuk no_rawat %s: %+v", req.NoRawat, valErrs)
		return valErrs
	}

	return nil
}



func (s *service) validasiWaktuRegistrasi(ctx context.Context, noRawat, tglPeresepan, jamPeresepan string, statusLanjut shared.StatusLanjut) error {
	tglRegStr, jamRegStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil data registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data registrasi kunjungan pasien tidak ditemukan")
	}

	waktuRegistrasi, err := shared.ParseWaktu(tglRegStr, jamRegStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tglRegStr, jamRegStr, err)
		return err
	}

	waktuPeresepan, err := shared.ParseWaktu(tglPeresepan, jamPeresepan)
	if err != nil {
		return apperror.NewBusinessError(err.Error())
	}

	if waktuPeresepan.Before(waktuRegistrasi) {
		errs := apperror.ValidationError{
			"tanggal_peresepan": fmt.Sprintf("Waktu peresepan (%s %s) tidak boleh lebih awal dari waktu registrasi pasien (%s %s)", tglPeresepan, jamPeresepan, tglRegStr, jamRegStr),
		}
		s.log.Warn("Validasi waktu peresepan gagal untuk no_rawat %s: %+v", noRawat, errs)
		return errs
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if statusLanjut == shared.StatusLanjutRawatJalan && time.Now().After(batasWaktu) {
		errMsg := fmt.Sprintf("Batas waktu peresepan obat untuk kunjungan rawat jalan ini telah berakhir (maksimal %d jam dari waktu registrasi: %s %s)", s.maxEditJam, tglRegStr, jamRegStr)
		s.log.Warn("Peresepan ditolak karena lewat batas %d jam untuk no_rawat %s: %s", s.maxEditJam, noRawat, errMsg)
		return apperror.NewForbiddenError(errMsg)
	}

	return nil
}

func (s *service) validasiStatusKamarInap(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) error {
	isAktifRanap, isPernahRanap, err := s.repo.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat %s: %v", noRawat, err)
		return err
	}

	if isAktifRanap && statusLanjut == shared.StatusLanjutRawatJalan {
		return apperror.NewBusinessError("Pasien sedang dirawat inap aktif. Peresepan obat wajib menggunakan status 'Ranap'.")
	}

	if !isAktifRanap && isPernahRanap && statusLanjut == shared.StatusLanjutRawatInap {
		return apperror.NewBusinessError("Pasien telah checkout / keluar dari rawat inap. Tidak dapat membuat resep baru untuk kunjungan ini.")
	}

	if !isAktifRanap && !isPernahRanap && statusLanjut == shared.StatusLanjutRawatInap {
		return apperror.NewBusinessError("Pasien belum/tidak terdaftar di kamar inap. Peresepan obat harus menggunakan status 'Ralan'.")
	}

	return nil
}

func (s *service) DaftarAturanPakai(ctx context.Context, keyword string) ([]AturanPakai, error) {
	data, err := s.repo.DaftarAturanPakai(ctx, keyword)
	if err != nil {
		s.log.Error("Gagal mengambil daftar aturan pakai: %v", err)
		return nil, err
	}

	return data, nil
}

func (s *service) DaftarMetodeRacik(ctx context.Context) ([]MetodeRacik, error) {
	data, err := s.repo.DaftarMetodeRacik(ctx)
	if err != nil {
		s.log.Error("Gagal mengambil daftar metode racik: %v", err)
		return nil, err
	}

	return data, nil
}
