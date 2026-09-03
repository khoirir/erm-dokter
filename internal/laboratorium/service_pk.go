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

func (s *service) SimpanPermintaanLabPK(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPKRequest) (*DetailPermintaanLabPK, error) {
	if !statusLanjut.IsValid() {
		return nil, apperror.NewBusinessError("Status lanjut tidak valid (pilihan: Ralan, Ranap)")
	}

	noRawat := req.NoRawat

	kunjungan, err := s.repo.GetKunjunganForPermintaanPK(ctx, noRawat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
		}
		s.log.Error("Gagal mengambil data kunjungan untuk permintaan lab PK %s: %v", noRawat, err)
		return nil, err
	}

	if kunjungan.StatusBayar == "Sudah Bayar" && kunjungan.KodePenjamin == "BPJ" {
		return nil, apperror.NewBusinessError("Pasien BPJS yang sudah menyelesaikan pembayaran / administrasi tidak dapat mengirim permintaan laboratorium baru")
	}

	tglRegStr, jamRegStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil data registrasi no_rawat %s: %v", noRawat, err)
		return nil, err
	}
	if !exists {
		return nil, apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

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
	templateMap := make(map[string][]int)
	seenTindakan := make(map[string]bool)

	for _, item := range req.Pemeriksaan {
		if !seenTindakan[item.KodeTindakan] {
			seenTindakan[item.KodeTindakan] = true
			kodeTindakanList = append(kodeTindakanList, item.KodeTindakan)
		}
		if len(item.KodeTemplate) > 0 {
			templateMap[item.KodeTindakan] = append(templateMap[item.KodeTindakan], item.KodeTemplate...)
		}
	}

	foundTindakan, err := s.tindakanService.CekKeberadaanTindakanLab(ctx, shared.KategoriLabPK, kodeTindakanList)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan tindakan lab: %v", err)
		return nil, err
	}

	valErrs := make(apperror.ValidationError)
	for i, item := range req.Pemeriksaan {
		if !foundTindakan[item.KodeTindakan] {
			valErrs[fmt.Sprintf("pemeriksaan[%d].id_tindakan", i)] = fmt.Sprintf("Pemeriksaan ke-%d: Data master tindakan laboratorium tidak ditemukan", i+1)
		}
	}
	if len(valErrs) > 0 {
		s.log.Warn("Validasi keberadaan tindakan lab gagal untuk no_rawat %s: %+v", noRawat, valErrs)
		return nil, valErrs
	}

	foundTemplates, err := s.tindakanService.CekKeberadaanTemplateLab(ctx, kodeTindakanList, templateMap)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan template lab: %v", err)
		return nil, err
	}

	for i, item := range req.Pemeriksaan {
		prefix := fmt.Sprintf("pemeriksaan[%d]", i)
		tmplMapForTindakan := foundTemplates[item.KodeTindakan]
		for j, idTmpl := range item.KodeTemplate {
			if tmplMapForTindakan == nil || !tmplMapForTindakan[idTmpl] {
				valErrs[fmt.Sprintf("%s.id_template[%d]", prefix, j)] = fmt.Sprintf("Pemeriksaan ke-%d parameter ke-%d: Template pengujian tidak terdaftar pada tindakan ini", i+1, j+1)
			}
		}
	}
	if len(valErrs) > 0 {
		s.log.Warn("Validasi template lab gagal untuk no_rawat %s: %+v", noRawat, valErrs)
		return nil, valErrs
	}

	noPermintaan, err := s.repo.SimpanPermintaanLabPK(ctx, noRawat, kodeDokterLogin, statusDB, req, kodeTindakanList, templateMap)
	if err != nil {
		s.log.Error("Gagal menyimpan permintaan laboratorium PK %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil membuat permintaan laboratorium PK %s untuk no_rawat %s oleh dokter %s", noPermintaan, noRawat, kodeDokterLogin)
	return s.GetDetailPermintaanLabPK(ctx, noRawat, noPermintaan, statusLanjut)
}

func (s *service) GetDaftarPermintaanLabPK(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut) ([]PermintaanLabPK, error) {
	items, err := s.repo.DaftarPermintaanLabPK(ctx, noRawat, statusLanjut)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat permintaan lab PK %s: %v", noRawat, err)
		return nil, err
	}

	result := make([]PermintaanLabPK, 0, len(items))
	for _, item := range items {
		result = append(result, s.formatPermintaanLabPK(item))
	}

	return result, nil
}

func (s *service) GetRiwayatPermintaanLabPKByRM(ctx context.Context, noRM string, statusLanjut shared.StatusLanjut, filter FilterRiwayatLab) ([]PermintaanLabPK, shared.PaginationMeta, error) {
	items, total, err := s.repo.DaftarPermintaanLabPKByRM(ctx, noRM, statusLanjut, filter)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat permintaan lab PK by RM %s: %v", noRM, err)
		return nil, shared.PaginationMeta{}, err
	}

	result := make([]PermintaanLabPK, 0, len(items))
	for _, item := range items {
		result = append(result, s.formatPermintaanLabPK(item))
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return result, meta, nil
}

func (s *service) GetDetailPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut) (*DetailPermintaanLabPK, error) {
	detail, err := s.repo.DetailPermintaanLabPK(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
		}
		s.log.Error("Gagal mengambil detail permintaan lab PK %s: %v", noPermintaan, err)
		return nil, err
	}

	if detail.NoRawat != noRawat {
		return nil, apperror.NewBusinessError("Permintaan laboratorium tidak sesuai dengan kunjungan pasien")
	}

	if statusLanjut != "Semua" && statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return nil, apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
	}

	detail.PermintaanLabPK = s.formatPermintaanLabPK(detail.PermintaanLabPK)
	return detail, nil
}

func (s *service) HapusPermintaanLabPK(ctx context.Context, noRawat string, noPermintaan string, statusLanjut shared.StatusLanjut, kodeDokterLogin string) error {
	detail, err := s.repo.DetailPermintaanLabPK(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
		}
		s.log.Error("Gagal mengambil data permintaan lab PK untuk dihapus %s: %v", noPermintaan, err)
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

	if err := s.repo.HapusPermintaanLabPK(ctx, noPermintaan); err != nil {
		s.log.Error("Gagal menghapus permintaan lab PK di repository %s: %v", noPermintaan, err)
		return err
	}

	s.log.Info("Berhasil menghapus permintaan lab PK %s untuk no_rawat %s oleh dokter %s", noPermintaan, noRawat, kodeDokterLogin)
	return nil
}

func (s *service) formatPermintaanLabPK(item PermintaanLabPK) PermintaanLabPK {
	statusProses := "Menunggu Sampel"
	if item.TanggalHasil != "0000-00-00" && strings.TrimSpace(item.TanggalHasil) != "" {
		statusProses = "Selesai"
	} else if item.TanggalSampel != "0000-00-00" && strings.TrimSpace(item.TanggalSampel) != "" {
		statusProses = "Sampel Diambil"
	}

	tglSampel := item.TanggalSampel
	if tglSampel == "0000-00-00" {
		tglSampel = ""
	}
	jamSampel := item.JamSampel
	if jamSampel == "00:00:00" {
		jamSampel = ""
	}
	tglHasil := item.TanggalHasil
	if tglHasil == "0000-00-00" {
		tglHasil = ""
	}
	jamHasil := item.JamHasil
	if jamHasil == "00:00:00" {
		jamHasil = ""
	}

	return PermintaanLabPK{
		NoPermintaan:      item.NoPermintaan,
		NoRawat:           item.NoRawat,
		TanggalPermintaan: item.TanggalPermintaan,
		JamPermintaan:     item.JamPermintaan,
		TanggalSampel:     tglSampel,
		JamSampel:         jamSampel,
		TanggalHasil:      tglHasil,
		JamHasil:          jamHasil,
		KodeDokterPerujuk: item.KodeDokterPerujuk,
		NamaDokterPerujuk: item.NamaDokterPerujuk,
		Status:            item.Status,
		InformasiTambahan: item.InformasiTambahan,
		DiagnosaKlinis:    item.DiagnosaKlinis,
		StatusProses:      statusProses,
	}
}
