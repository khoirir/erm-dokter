package laboratorium

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

func (s *service) SimpanPermintaanLabPK(ctx context.Context, kodeDokterLogin string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPKRequest) (*DetailPermintaanLabPK, error) {
	noRawat := req.NoRawat

	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, req.TanggalPermintaan, req.JamPermintaan, statusLanjut, "membuat"); err != nil {
		return nil, err
	}

	kodeTindakanList, templateMap, err := s.validasiTindakanDanTemplateLabPK(ctx, noRawat, req)
	if err != nil {
		return nil, err
	}

	noPermintaan, err := s.repo.SimpanPermintaanLabPK(ctx, noRawat, kodeDokterLogin, statusLanjut, req, kodeTindakanList, templateMap)
	if err != nil {
		s.log.Error("Gagal menyimpan permintaan laboratorium PK %s: %v", noRawat, err)
		return nil, err
	}

	s.log.Info("Berhasil membuat permintaan laboratorium PK %s untuk no_rawat %s oleh dokter %s", noPermintaan, noRawat, kodeDokterLogin)
	return s.GetDetailPermintaanLabPK(ctx, noRawat, noPermintaan, statusLanjut)
}

func (s *service) UpdatePermintaanLabPK(ctx context.Context, kodeDokterLogin, noRawat, noPermintaan string, statusLanjut shared.StatusLanjut, req SimpanPermintaanLabPKRequest) (*DetailPermintaanLabPK, error) {
	detail, err := s.repo.DetailPermintaanLabPK(ctx, noPermintaan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
		}
		s.log.Error("Gagal mengambil data permintaan lab PK untuk diubah %s: %v", noPermintaan, err)
		return nil, err
	}

	if detail.NoRawat != noRawat {
		return nil, apperror.NewBusinessError("Permintaan laboratorium tidak sesuai dengan kunjungan pasien")
	}

	if statusLanjut != "Semua" && statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return nil, apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
	}

	if detail.KodeDokterPerujuk != kodeDokterLogin {
		s.log.Warn("Percobaan mengubah permintaan lab PK %s oleh dokter %s ditolak: dibuat oleh %s (%s)", noPermintaan, kodeDokterLogin, detail.KodeDokterPerujuk, detail.NamaDokterPerujuk)
		return nil, apperror.NewForbiddenError("Hanya dokter pemohon yang berhak mengubah permintaan laboratorium ini")
	}

	isSampelDiambil := detail.TanggalSampel != "0000-00-00"
	isHasilKeluar := detail.TanggalHasil != "0000-00-00"

	if isSampelDiambil || isHasilKeluar {
		return nil, apperror.NewBusinessError("Permintaan laboratorium sudah diproses (sudah diambil sampel atau hasil sudah keluar) dan tidak dapat diubah")
	}

	if err := s.validasiRegistrasiDanStatus(ctx, req.NoRawat, req.TanggalPermintaan, req.JamPermintaan, statusLanjut, "mengubah"); err != nil {
		return nil, err
	}

	kodeTindakanList, templateMap, err := s.validasiTindakanDanTemplateLabPK(ctx, noRawat, req)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePermintaanLabPK(ctx, noPermintaan, req, kodeTindakanList, templateMap); err != nil {
		s.log.Error("Gagal memperbarui permintaan lab PK %s: %v", noPermintaan, err)
		return nil, err
	}

	s.log.Info("Berhasil memperbarui permintaan lab PK %s untuk no_rawat %s oleh dokter %s", noPermintaan, noRawat, kodeDokterLogin)

	updatedDetail, err := s.repo.DetailPermintaanLabPK(ctx, noPermintaan)
	if err != nil {
		s.log.Error("Gagal mengambil detail setelah update permintaan lab PK %s: %v", noPermintaan, err)
		return nil, err
	}

	updatedDetail.PermintaanLabPK = s.formatPermintaanLabPK(updatedDetail.PermintaanLabPK)
	return updatedDetail, nil
}

func (s *service) validasiTindakanDanTemplateLabPK(ctx context.Context, noRawat string, req SimpanPermintaanLabPKRequest) ([]string, map[string][]int, error) {
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
		return nil, nil, err
	}

	valErrs := make(apperror.ValidationError)
	for i, item := range req.Pemeriksaan {
		if !foundTindakan[item.KodeTindakan] {
			valErrs[fmt.Sprintf("pemeriksaan[%d].id_tindakan", i)] = fmt.Sprintf("Pemeriksaan ke-%d: Data master tindakan laboratorium tidak ditemukan", i+1)
		}
	}
	if len(valErrs) > 0 {
		s.log.Warn("Validasi keberadaan tindakan lab gagal untuk no_rawat %s: %+v", noRawat, valErrs)
		return nil, nil, valErrs
	}

	foundTemplates, err := s.tindakanService.CekKeberadaanTemplateLab(ctx, kodeTindakanList, templateMap)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan template lab: %v", err)
		return nil, nil, err
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
		return nil, nil, valErrs
	}

	return kodeTindakanList, templateMap, nil
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

	if statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
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

	if statusLanjut != "" && !strings.EqualFold(detail.Status, string(statusLanjut)) {
		return apperror.NewNotFoundError("Data permintaan laboratorium tidak ditemukan")
	}

	if detail.KodeDokterPerujuk != kodeDokterLogin {
		return apperror.NewForbiddenError("Hanya dokter pemohon yang berhak membatalkan permintaan laboratorium ini")
	}

	isSampelDiambil := detail.TanggalSampel != "0000-00-00"
	isHasilKeluar := detail.TanggalHasil != "0000-00-00"

	if isSampelDiambil || isHasilKeluar {
		return apperror.NewBusinessError("Permintaan laboratorium sudah diproses (sudah diambil sampel atau hasil sudah keluar) dan tidak dapat dibatalkan")
	}

	orderStatus := statusLanjut
	if orderStatus == "" || orderStatus == "Semua" {
		if strings.EqualFold(detail.Status, string(shared.StatusLanjutRawatInap)) {
			orderStatus = shared.StatusLanjutRawatInap
		} else {
			orderStatus = shared.StatusLanjutRawatJalan
		}
	}

	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, "", "", orderStatus, "menghapus"); err != nil {
		return err
	}

	if err := s.repo.HapusPermintaanLabPK(ctx, noPermintaan); err != nil {
		s.log.Error("Gagal menghapus permintaan lab PK di repository %s: %v", noPermintaan, err)
		return err
	}

	s.log.Info("Berhasil menghapus permintaan lab PK %s untuk no_rawat %s oleh dokter %s", noPermintaan, noRawat, kodeDokterLogin)
	return nil
}

func (s *service) formatPermintaanLabPK(item PermintaanLabPK) PermintaanLabPK {
	item.FormatZeroDates()
	return item
}
