package rujukaninternal

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	DaftarOpsiPoliDokter(ctx context.Context, kodeDokterLogin, keyword string) ([]OpsiPoliDokter, error)
	DaftarRujukanInternal(ctx context.Context, noRawat string) ([]RujukanInternal, error)
	SimpanRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanRujukanRequest) (*RujukanInternal, error)
	HapusRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat, idRujukan string) error
}

type service struct {
	repo              Repository
	rawatJalanService rawatjalan.Service
	log               *logger.Logger
	encryptionKey     string
}

func NewService(repo Repository, rawatJalanService rawatjalan.Service, log *logger.Logger, encryptionKey string) Service {
	return &service{
		repo:              repo,
		rawatJalanService: rawatJalanService,
		log:               log,
		encryptionKey:     encryptionKey,
	}
}

func (s *service) DaftarOpsiPoliDokter(ctx context.Context, kodeDokterLogin, keyword string) ([]OpsiPoliDokter, error) {
	cleanKeyword := strings.TrimSpace(keyword)
	list, err := s.repo.DaftarOpsiPoliDokter(ctx, kodeDokterLogin, cleanKeyword)
	if err != nil {
		s.log.Error("Gagal mengambil opsi poli dokter untuk dokter %s (keyword: %s): %v", kodeDokterLogin, cleanKeyword, err)
		return nil, err
	}

	for i := range list {
		item := &list[i]
		if encrypted, err := crypto.Encrypt(item.CompositeKey(), s.encryptionKey); err == nil {
			item.Id = encrypted
		}
	}

	return list, nil
}

func (s *service) DaftarRujukanInternal(ctx context.Context, noRawat string) ([]RujukanInternal, error) {
	list, err := s.repo.DaftarRujukanInternalByNoRawat(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil daftar rujukan internal no_rawat %s: %v", noRawat, err)
		return nil, err
	}

	for i := range list {
		item := &list[i]
		if encrypted, err := crypto.Encrypt(item.CompositeKey(), s.encryptionKey); err == nil {
			item.Id = encrypted
		}
		if encryptedIdKunjungan, err := crypto.Encrypt(item.NoRawat, s.encryptionKey); err == nil {
			item.IdKunjungan = encryptedIdKunjungan
		}
	}

	return list, nil
}

func (s *service) SimpanRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat string, req SimpanRujukanRequest) (*RujukanInternal, error) {
	if errs := req.Validate(); errs != nil {
		return nil, errs
	}

	decryptedTujuan, err := crypto.Decrypt(req.IdTujuan, s.encryptionKey)
	if err != nil {
		s.log.Warn("Gagal dekripsi id_tujuan rujukan internal %s: %v", req.IdTujuan, err)
		return nil, apperror.NewBusinessError("ID tujuan rujukan tidak valid")
	}

	target, err := ParseIdOpsiPoliDokter(decryptedTujuan)
	if err != nil {
		return nil, apperror.NewBusinessError(err.Error())
	}

	if target.KodeDokter == kodeDokterLogin {
		s.log.Warn("Dokter %s mencoba membuat rujukan internal ke dirinya sendiri pada no_rawat %s", kodeDokterLogin, noRawat)
		return nil, apperror.ValidationError{
			"id_tujuan": "Tidak dapat membuat rujukan internal ke diri sendiri",
		}
	}

	if err := s.validasiBatasWaktu(ctx, noRawat, "dibuat"); err != nil {
		return nil, err
	}

	ada, err := s.repo.CekRujukanInternalAda(ctx, noRawat, target.KodeDokter)
	if err != nil {
		s.log.Error("Gagal memeriksa duplikasi rujukan internal no_rawat %s dokter %s: %v", noRawat, target.KodeDokter, err)
		return nil, err
	}
	if ada {
		s.log.Warn("Percobaan duplikasi rujukan internal no_rawat %s ke dokter %s", noRawat, target.KodeDokter)
		return nil, apperror.ValidationError{
			"id_tujuan": "Pasien sudah pernah dirujuk ke dokter ini pada kunjungan yang sama",
		}
	}

	if err := s.repo.SimpanRujukanInternal(ctx, noRawat, target.KodePoli, target.KodeDokter); err != nil {
		s.log.Error("Gagal menyimpan rujukan internal no_rawat %s poli %s dokter %s: %v", noRawat, target.KodePoli, target.KodeDokter, err)
		return nil, err
	}

	s.log.Info("Berhasil menyimpan rujukan internal no_rawat %s ke poli %s dokter %s oleh dokter %s", noRawat, target.KodePoli, target.KodeDokter, kodeDokterLogin)

	list, err := s.repo.DaftarRujukanInternalByNoRawat(ctx, noRawat)
	if err != nil {
		return nil, nil
	}

	var hasil *RujukanInternal
	for i := range list {
		if list[i].KodeDokter == target.KodeDokter && list[i].KodePoli == target.KodePoli {
			hasil = &list[i]
			if encrypted, err := crypto.Encrypt(hasil.CompositeKey(), s.encryptionKey); err == nil {
				hasil.Id = encrypted
			}
			if encryptedIdKunjungan, err := crypto.Encrypt(hasil.NoRawat, s.encryptionKey); err == nil {
				hasil.IdKunjungan = encryptedIdKunjungan
			}
			break
		}
	}

	return hasil, nil
}

func (s *service) HapusRujukanInternal(ctx context.Context, kodeDokterLogin, noRawat, idRujukan string) error {
	decryptedId, err := crypto.Decrypt(idRujukan, s.encryptionKey)
	if err != nil {
		s.log.Warn("Gagal dekripsi id_rujukan %s: %v", idRujukan, err)
		return apperror.NewBusinessError("ID rujukan tidak valid")
	}

	parsedId, err := ParseIdRujukanInternal(decryptedId)
	if err != nil {
		return apperror.NewBusinessError(err.Error())
	}

	if parsedId.NoRawat != noRawat {
		s.log.Warn("Percobaan menghapus rujukan internal no_rawat mismatch: URL %s vs ID %s", noRawat, parsedId.NoRawat)
		return apperror.NewForbiddenError("Data rujukan internal tidak sesuai dengan kunjungan pasien")
	}

	if err := s.validasiBatasWaktu(ctx, noRawat, "menghapus"); err != nil {
		return err
	}

	if err := s.repo.HapusRujukanInternal(ctx, noRawat, parsedId.KodeDokter); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFoundError("Data rujukan internal tidak ditemukan")
		}
		s.log.Error("Gagal menghapus rujukan internal no_rawat %s dokter %s: %v", noRawat, parsedId.KodeDokter, err)
		return err
	}

	s.log.Info("Berhasil menghapus rujukan internal no_rawat %s ke dokter %s oleh dokter %s", noRawat, parsedId.KodeDokter, kodeDokterLogin)
	return nil
}

func (s *service) validasiBatasWaktu(ctx context.Context, noRawat, aksi string) error {
	tglRegStr, jamRegStr, exists, err := s.rawatJalanService.GetWaktuRegistrasi(ctx, noRawat)
	if err != nil {
		s.log.Error("Gagal mengambil waktu registrasi no_rawat %s: %v", noRawat, err)
		return err
	}
	if !exists {
		return apperror.NewNotFoundError("Data kunjungan pasien tidak ditemukan")
	}

	if err := shared.ValidasiBatasWaktuRekamMedis(tglRegStr, jamRegStr, 48, aksi); err != nil {
		s.log.Warn("Validasi batas waktu 48 jam rujukan internal gagal untuk no_rawat %s (%s %s): %v", noRawat, tglRegStr, jamRegStr, err)
		return err
	}

	return nil
}
