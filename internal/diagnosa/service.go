package diagnosa

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"erm-dokter/internal/master"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarDiagnosaProsedur(ctx context.Context, noRawat string, status shared.StatusLanjut) (*DaftarDiagnosaProsedurResponse, error)
	RiwayatPasien(ctx context.Context, noRekamMedis string, status shared.StatusLanjut) (*RiwayatPasienResponse, error)
	TambahDiagnosa(ctx context.Context, noRawat string, status shared.StatusLanjut, req TambahDiagnosaRequest) (*DiagnosaPasien, error)
	UpdateDiagnosa(ctx context.Context, id IdDiagnosa, req UpdateDiagnosaRequest) error
	HapusDiagnosa(ctx context.Context, id IdDiagnosa) error
	ReorderDiagnosa(ctx context.Context, noRawat string, status shared.StatusLanjut, req ReorderRequest) error
	TambahProsedur(ctx context.Context, noRawat string, status shared.StatusLanjut, req TambahProsedurRequest) (*ProsedurPasien, error)
	UpdateProsedur(ctx context.Context, id IdProsedur, req UpdateProsedurRequest) error
	HapusProsedur(ctx context.Context, id IdProsedur) error
	ReorderProsedur(ctx context.Context, noRawat string, status shared.StatusLanjut, req ReorderRequest) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	rawatInapService  rawatinap.Service
	masterService     master.Service
	maxEditJam        int
	log               *logger.Logger
}

func NewService(
	repo Repository,
	rawatJalanService rawatjalan.Service,
	rawatInapService rawatinap.Service,
	masterService master.Service,
	maxEditJam int,
	log *logger.Logger,
) Service {
	if maxEditJam <= 0 {
		maxEditJam = 48
	}
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		rawatInapService:  rawatInapService,
		masterService:     masterService,
		maxEditJam:        maxEditJam,
		log:               log,
	}
}

func (s *service) DaftarDiagnosaProsedur(ctx context.Context, noRawat string, status shared.StatusLanjut) (*DaftarDiagnosaProsedurResponse, error) {
	diagnosaList, err := s.repo.GetDiagnosaByNoRawat(ctx, noRawat, status)
	if err != nil {
		s.log.Error("Gagal mengambil daftar diagnosa no_rawat '%s': %v", noRawat, err)
		return nil, err
	}

	prosedurList, err := s.repo.GetProsedurByNoRawat(ctx, noRawat, status)
	if err != nil {
		s.log.Error("Gagal mengambil daftar prosedur no_rawat '%s': %v", noRawat, err)
		return nil, err
	}

	return &DaftarDiagnosaProsedurResponse{
		Diagnosa: diagnosaList,
		Prosedur: prosedurList,
	}, nil
}

func (s *service) RiwayatPasien(ctx context.Context, noRekamMedis string, status shared.StatusLanjut) (*RiwayatPasienResponse, error) {
	riwayatKunjungan, err := s.rawatJalanService.RiwayatKunjunganPasien(ctx, noRekamMedis)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat kunjungan RM '%s': %v", noRekamMedis, err)
		return nil, err
	}

	if len(riwayatKunjungan) == 0 {
		return &RiwayatPasienResponse{
			Diagnosa: make([]DiagnosaPasien, 0),
			Prosedur: make([]ProsedurPasien, 0),
		}, nil
	}

	listNoRawat := make([]string, len(riwayatKunjungan))
	for i, k := range riwayatKunjungan {
		listNoRawat[i] = k.NoRawat
	}

	diagnosaList, err := s.repo.GetRiwayatDiagnosaByListNoRawat(ctx, listNoRawat, status)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat diagnosa RM '%s': %v", noRekamMedis, err)
		return nil, err
	}

	prosedurList, err := s.repo.GetRiwayatProsedurByListNoRawat(ctx, listNoRawat, status)
	if err != nil {
		s.log.Error("Gagal mengambil riwayat prosedur RM '%s': %v", noRekamMedis, err)
		return nil, err
	}

	return &RiwayatPasienResponse{
		Diagnosa: diagnosaList,
		Prosedur: prosedurList,
	}, nil
}

func (s *service) TambahDiagnosa(ctx context.Context, noRawat string, status shared.StatusLanjut, req TambahDiagnosaRequest) (*DiagnosaPasien, error) {
	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, status, "menambah diagnosa"); err != nil {
		return nil, err
	}

	existsMap, err := s.masterService.CekKeberadaanICD10(ctx, []string{req.Kode})
	if err != nil {
		s.log.Error("Gagal cek keberadaan master ICD-10 '%s': %v", req.Kode, err)
		return nil, err
	}
	if !existsMap[req.Kode] {
		return nil, apperror.NewNotFoundError(fmt.Sprintf("Kode ICD-10 '%s' tidak ditemukan", req.Kode))
	}

	isDup, err := s.repo.CekDuplikasiDiagnosa(ctx, noRawat, req.Kode, status)
	if err != nil {
		s.log.Error("Gagal cek duplikasi diagnosa no_rawat '%s', kode '%s': %v", noRawat, req.Kode, err)
		return nil, err
	}
	if isDup {
		return nil, apperror.NewBusinessError(fmt.Sprintf("Diagnosa dengan kode '%s' sudah diinput", req.Kode))
	}

	prioritas := 1
	if req.Prioritas != nil {
		prioritas = *req.Prioritas
	} else {
		maxPrio, err := s.repo.GetMaxPrioritasDiagnosa(ctx, noRawat, status)
		if err != nil {
			s.log.Error("Gagal mengambil max prioritas diagnosa: %v", err)
			return nil, err
		}
		prioritas = maxPrio + 1
	}

	statusPenyakit := req.StatusPenyakit
	if statusPenyakit == "" {
		statusPenyakit = s.tentukanStatusPenyakit(ctx, noRawat, req.Kode, status)
	}

	diagnosa := DiagnosaPasien{
		NoRawat:        noRawat,
		Kode:           req.Kode,
		Status:         status,
		Prioritas:      prioritas,
		StatusPenyakit: statusPenyakit,
	}

	if err := s.repo.TambahDiagnosa(ctx, diagnosa); err != nil {
		s.log.Error("Gagal menambah diagnosa no_rawat '%s', kode '%s': %v", noRawat, req.Kode, err)
		return nil, err
	}

	listDiagnosa, err := s.repo.GetDiagnosaByNoRawat(ctx, noRawat, status)
	if err != nil {
		return &diagnosa, nil
	}
	for _, item := range listDiagnosa {
		if item.Kode == req.Kode {
			return &item, nil
		}
	}

	return &diagnosa, nil
}

func (s *service) UpdateDiagnosa(ctx context.Context, id IdDiagnosa, req UpdateDiagnosaRequest) error {
	if err := s.validasiRegistrasiDanStatus(ctx, id.NoRawat, id.Status, "mengubah diagnosa"); err != nil {
		return err
	}

	err := s.repo.UpdateDiagnosa(ctx, id.NoRawat, id.Kode, id.Status, req.Prioritas, req.StatusPenyakit)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError("Data diagnosa tidak ditemukan")
		}
		s.log.Error("Gagal update diagnosa no_rawat '%s', kode '%s': %v", id.NoRawat, id.Kode, err)
		return err
	}
	return nil
}

func (s *service) HapusDiagnosa(ctx context.Context, id IdDiagnosa) error {
	if err := s.validasiRegistrasiDanStatus(ctx, id.NoRawat, id.Status, "menghapus diagnosa"); err != nil {
		return err
	}

	err := s.repo.HapusDiagnosa(ctx, id.NoRawat, id.Kode, id.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError("Data diagnosa tidak ditemukan")
		}
		s.log.Error("Gagal menghapus diagnosa no_rawat '%s', kode '%s': %v", id.NoRawat, id.Kode, err)
		return err
	}
	return nil
}

func (s *service) ReorderDiagnosa(ctx context.Context, noRawat string, status shared.StatusLanjut, req ReorderRequest) error {
	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, status, "mengurutkan diagnosa"); err != nil {
		return err
	}

	err := s.repo.ReorderDiagnosa(ctx, noRawat, status, req.Items)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError("Salah satu item diagnosa dalam urutan tidak ditemukan")
		}
		s.log.Error("Gagal reorder diagnosa no_rawat '%s': %v", noRawat, err)
		return err
	}
	return nil
}

func (s *service) TambahProsedur(ctx context.Context, noRawat string, status shared.StatusLanjut, req TambahProsedurRequest) (*ProsedurPasien, error) {
	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, status, "menambah prosedur"); err != nil {
		return nil, err
	}

	existsMap, err := s.masterService.CekKeberadaanICD9(ctx, []string{req.Kode})
	if err != nil {
		s.log.Error("Gagal cek keberadaan master ICD-9 '%s': %v", req.Kode, err)
		return nil, err
	}
	if !existsMap[req.Kode] {
		return nil, apperror.NewNotFoundError(fmt.Sprintf("Kode ICD-9 '%s' tidak ditemukan", req.Kode))
	}

	isDup, err := s.repo.CekDuplikasiProsedur(ctx, noRawat, req.Kode, status)
	if err != nil {
		s.log.Error("Gagal cek duplikasi prosedur no_rawat '%s', kode '%s': %v", noRawat, req.Kode, err)
		return nil, err
	}
	if isDup {
		return nil, apperror.NewBusinessError(fmt.Sprintf("Prosedur dengan kode '%s' sudah diinput", req.Kode))
	}

	prioritas := 1
	if req.Prioritas != nil {
		prioritas = *req.Prioritas
	} else {
		maxPrio, err := s.repo.GetMaxPrioritasProsedur(ctx, noRawat, status)
		if err != nil {
			s.log.Error("Gagal mengambil max prioritas prosedur: %v", err)
			return nil, err
		}
		prioritas = maxPrio + 1
	}

	prosedur := ProsedurPasien{
		NoRawat:   noRawat,
		Kode:      req.Kode,
		Status:    status,
		Prioritas: prioritas,
	}

	if err := s.repo.TambahProsedur(ctx, prosedur); err != nil {
		s.log.Error("Gagal menambah prosedur no_rawat '%s', kode '%s': %v", noRawat, req.Kode, err)
		return nil, err
	}

	listProsedur, err := s.repo.GetProsedurByNoRawat(ctx, noRawat, status)
	if err != nil {
		return &prosedur, nil
	}
	for _, item := range listProsedur {
		if item.Kode == req.Kode {
			return &item, nil
		}
	}

	return &prosedur, nil
}

func (s *service) UpdateProsedur(ctx context.Context, id IdProsedur, req UpdateProsedurRequest) error {
	if err := s.validasiRegistrasiDanStatus(ctx, id.NoRawat, id.Status, "mengubah prosedur"); err != nil {
		return err
	}

	err := s.repo.UpdateProsedur(ctx, id.NoRawat, id.Kode, id.Status, *req.Prioritas)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError("Data prosedur tidak ditemukan")
		}
		s.log.Error("Gagal update prosedur no_rawat '%s', kode '%s': %v", id.NoRawat, id.Kode, err)
		return err
	}
	return nil
}

func (s *service) HapusProsedur(ctx context.Context, id IdProsedur) error {
	if err := s.validasiRegistrasiDanStatus(ctx, id.NoRawat, id.Status, "menghapus prosedur"); err != nil {
		return err
	}

	err := s.repo.HapusProsedur(ctx, id.NoRawat, id.Kode, id.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError("Data prosedur tidak ditemukan")
		}
		s.log.Error("Gagal menghapus prosedur no_rawat '%s', kode '%s': %v", id.NoRawat, id.Kode, err)
		return err
	}
	return nil
}

func (s *service) ReorderProsedur(ctx context.Context, noRawat string, status shared.StatusLanjut, req ReorderRequest) error {
	if err := s.validasiRegistrasiDanStatus(ctx, noRawat, status, "mengurutkan prosedur"); err != nil {
		return err
	}

	err := s.repo.ReorderProsedur(ctx, noRawat, status, req.Items)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFoundError("Salah satu item prosedur dalam urutan tidak ditemukan")
		}
		s.log.Error("Gagal reorder prosedur no_rawat '%s': %v", noRawat, err)
		return err
	}
	return nil
}

func (s *service) tentukanStatusPenyakit(ctx context.Context, noRawat, kode string, status shared.StatusLanjut) string {
	kunjungan, err := s.rawatJalanService.DetailKunjungan(ctx, noRawat, "")
	if err != nil || kunjungan == nil || kunjungan.NoRekamMedis == "" {
		return "Baru"
	}

	riwayat, err := s.rawatJalanService.RiwayatKunjunganPasien(ctx, kunjungan.NoRekamMedis)
	if err != nil || len(riwayat) == 0 {
		return "Baru"
	}

	var listPastNoRawat []string
	for _, r := range riwayat {
		if r.NoRawat != noRawat {
			listPastNoRawat = append(listPastNoRawat, r.NoRawat)
		}
	}

	if len(listPastNoRawat) == 0 {
		return "Baru"
	}

	exists, err := s.repo.CekRiwayatPenyakit(ctx, listPastNoRawat, kode, status)
	if err != nil || !exists {
		return "Baru"
	}
	return "Lama"
}

func (s *service) validasiRegistrasiDanStatus(ctx context.Context, noRawat string, statusLanjut shared.StatusLanjut, aksi string) error {
	tanggalRegistrasi, jamRegistrasi, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil data registrasi no_rawat '%s': %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data registrasi kunjungan pasien tidak ditemukan")
	}

	waktuRegistrasi, err := shared.ParseWaktu(tanggalRegistrasi, jamRegistrasi)
	if err != nil {
		s.log.Error("Gagal parse waktu registrasi no_rawat '%s' (%s %s): %v", noRawat, tanggalRegistrasi, jamRegistrasi, err)
		return err
	}

	isAktifRanap, hasRecordKamar, err := s.rawatInapService.CekStatusKamarInap(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal cek status kamar inap untuk no_rawat '%s': %v", noRawat, err)
		return err
	}

	if hasRecordKamar {
		if !isAktifRanap {
			return apperror.NewBusinessError("Pasien sudah keluar dari kamar inap")
		}
		return nil
	}

	if strings.EqualFold(string(statusLanjut), string(shared.StatusLanjutRawatInap)) {
		return apperror.NewBusinessError("Pasien belum terdaftar di kamar inap, gunakan status 'ralan'")
	}

	batasWaktu := waktuRegistrasi.Add(time.Duration(s.maxEditJam) * time.Hour)
	if time.Now().After(batasWaktu) {
		return apperror.NewForbiddenError(
			fmt.Sprintf("Data diagnosa/prosedur tidak dapat %s, melebihi batas %d jam", aksi, s.maxEditJam),
			fmt.Sprintf("Kunjungan no_rawat '%s' ditolak untuk %s karena melewati batas %d jam dari registrasi (%s)", noRawat, aksi, s.maxEditJam, waktuRegistrasi.Format("2006-01-02 15:04:05")),
		)
	}

	return nil
}
