package laboratorium

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

func (s *service) SimpanPermintaanLabPA(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPARequest) (*DetailPermintaanLabPA, error) {
	if !statusLanjut.IsValid() {
		return nil, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)")
	}

	noRawat := req.NoRawat

	infoReg, err := s.rawatJalanService.GetInfoRegistrasi(ctx, noRawat)
	if err != nil {
		return nil, err
	}

	if infoReg.StatusBayar == "Sudah Bayar" && infoReg.KodePenjamin == "BPJ" {
		return nil, apperror.NewBusinessError("Pasien BPJS yang sudah menyelesaikan pembayaran / administrasi tidak dapat mengirim permintaan laboratorium baru")
	}

	tglRegStr := infoReg.TanggalRegistrasi
	jamRegStr := infoReg.JamRegistrasi
	waktuRegistrasi, err := shared.ParseWaktu(tglRegStr, jamRegStr)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat %s (%s %s): %v", noRawat, tglRegStr, jamRegStr, err)
		return nil, err
	}

	waktuPermintaan, err := shared.ParseWaktu(req.TanggalPermintaan, req.JamPermintaan)
	if err != nil {
		return nil, apperror.NewBusinessError(err.Error())
	}

	if waktuPermintaan.Before(waktuRegistrasi) {
		return nil, apperror.NewBusinessError(fmt.Sprintf("Waktu permintaan laboratorium (%s %s) tidak boleh mendahului waktu registrasi pasien (%s %s)", req.TanggalPermintaan, req.JamPermintaan, tglRegStr, jamRegStr))
	}

	isKamarAktif, hasRecordKamar, err := s.repo.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal memeriksa status kamar inap pasien %s: %v", noRawat, err)
		return nil, err
	}

	var statusDB string
	if statusLanjut == shared.StatusLanjutRawatInap {
		if !hasRecordKamar {
			return nil, apperror.NewBusinessError("Pasien tidak memiliki data kamar inap untuk order Ranap")
		}
		if !isKamarAktif {
			return nil, apperror.NewBusinessError("Pasien rawat inap sudah keluar / checkout dari kamar inap")
		}
		statusDB = "ranap"
	} else {
		if isKamarAktif {
			return nil, apperror.NewBusinessError("Pasien saat ini berstatus rawat inap aktif, permintaan laboratorium harus berstatus rawat inap")
		}
		batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
		if time.Now().After(batasWaktu) {
			return nil, apperror.NewBusinessError("Permintaan laboratorium rawat jalan telah melewati batas waktu 48 jam sejak waktu registrasi")
		}
		statusDB = "ralan"
	}

	var kodeTindakanList []string
	seenTindakan := make(map[string]bool)

	for _, item := range req.Pemeriksaan {
		if !seenTindakan[item.KodeTindakan] {
			seenTindakan[item.KodeTindakan] = true
			kodeTindakanList = append(kodeTindakanList, item.KodeTindakan)
		}
	}

	foundTindakan, err := s.tindakanService.CekKeberadaanTindakanLab(ctx, shared.KategoriLabPA, kodeTindakanList)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan tindakan lab PA: %v", err)
		return nil, err
	}

	valErrs := make(apperror.ValidationError)
	for i, item := range req.Pemeriksaan {
		if !foundTindakan[item.KodeTindakan] {
			valErrs[fmt.Sprintf("pemeriksaan[%d].id_tindakan", i)] = fmt.Sprintf("Pemeriksaan ke-%d: Data master tindakan laboratorium tidak ditemukan", i+1)
		}
	}
	if len(valErrs) > 0 {
		s.log.Warn("Validasi keberadaan tindakan lab PA gagal untuk no_rawat %s: %+v", noRawat, valErrs)
		return nil, valErrs
	}

	noPermintaan, err := s.repo.SimpanPermintaanLabPA(ctx, noRawat, kodeDokterLogin, statusDB, req, kodeTindakanList)
	if err != nil {
		s.log.Error("Gagal menyimpan permintaan laboratorium PA %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil membuat permintaan laboratorium PA %s untuk no_rawat %s oleh dokter %s", noPermintaan, noRawat, kodeDokterLogin)
	return s.GetDetailPermintaanLabPA(ctx, noRawat, noPermintaan, statusLanjut)
}

func (s *service) GetDaftarPermintaanLabPA(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPA, error) {
	items, err := s.repo.DaftarPermintaanLabPA(ctx, noRawat, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat permintaan lab PA %s: %v", noRawat, err)
		return nil, err
	}

	result := make([]PermintaanLabPA, 0, len(items))
	for _, item := range items {
		result = append(result, s.formatPermintaanLabPA(item))
	}

	return result, nil
}

func (s *service) GetRiwayatPermintaanLabPAByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPA, shared.PaginationMeta, error) {
	items, total, err := s.repo.DaftarPermintaanLabPAByRM(ctx, noRM, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat permintaan lab PA by RM %s: %v", noRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	result := make([]PermintaanLabPA, 0, len(items))
	for _, item := range items {
		result = append(result, s.formatPermintaanLabPA(item))
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return result, meta, nil
}

func (s *service) GetDetailPermintaanLabPA(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanLabPA, error) {
	detail, err := s.repo.DetailPermintaanLabPA(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail permintaan lab PA %s: %v", noPermintaan, err)
		return nil, err
	}

	if detail.NoRawat != noRawat {
		return nil, apperror.NewBusinessError("Permintaan laboratorium tidak sesuai dengan kunjungan pasien")
	}

	if statusLanjut != "Semua" && statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return nil, apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
	}

	detail.PermintaanLabPA = s.formatPermintaanLabPA(detail.PermintaanLabPA)
	return detail, nil
}

func (s *service) HapusPermintaanLabPA(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error {
	detail, err := s.repo.DetailPermintaanLabPA(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
		}
		s.log.Error("Gagal mengambil data permintaan lab PA untuk dihapus %s: %v", noPermintaan, err)
		return err
	}

	if detail.NoRawat != noRawat {
		return apperror.NewBusinessError("Permintaan laboratorium tidak sesuai dengan kunjungan pasien")
	}

	if statusLanjut != "Semua" && statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
	}

	if detail.KodeDokterPerujuk != kodeDokterLogin {
		return apperror.NewForbiddenError("Hanya dokter pemohon yang berhak membatalkan permintaan laboratorium ini")
	}

	isSampelDiambil := detail.TanggalSampel != "0000-00-00" && strings.TrimSpace(detail.TanggalSampel) != ""
	isHasilKeluar := detail.TanggalHasil != "0000-00-00" && strings.TrimSpace(detail.TanggalHasil) != ""

	if isSampelDiambil || isHasilKeluar {
		return apperror.NewBusinessError("Permintaan laboratorium sudah diproses (sudah diambil sampel atau hasil sudah keluar) dan tidak dapat dibatalkan")
	}

	if err := s.repo.HapusPermintaanLabPA(ctx, noPermintaan); err != nil {
		s.log.Error("Gagal menghapus permintaan lab PA di repository %s: %v", noPermintaan, err)
		return err
	}

	s.log.Info("Berhasil menghapus permintaan lab PA %s untuk no_rawat %s oleh dokter %s", noPermintaan, noRawat, kodeDokterLogin)
	return nil
}

func (s *service) formatPermintaanLabPA(item PermintaanLabPA) PermintaanLabPA {
	item.FormatZeroDates()
	if item.PengambilanBahan == "0000-00-00" {
		item.PengambilanBahan = ""
	}
	if item.TanggalPASebelumnya == "0000-00-00" {
		item.TanggalPASebelumnya = ""
	}
	return item
}
