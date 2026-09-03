package tindakan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"erm-dokter/internal/pkg/crypto"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	GetDaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter FilterDaftarTindakanLab) ([]TindakanLab, shared.PaginationMeta, error)
	GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, encryptedId string) (*DetailTindakanLab, error)
	CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error)
	CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error)
}

type service struct {
	repo      Repository
	jwtSecret string
	log       *logger.Logger
}

func NewService(repo Repository, jwtSecret string, log *logger.Logger) Service {
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}

func (s *service) GetDaftarTindakanLab(ctx context.Context, kategori shared.KategoriLab, filter FilterDaftarTindakanLab) ([]TindakanLab, shared.PaginationMeta, error) {
	list, total, err := s.repo.DaftarTindakanLab(ctx, kategori, filter)
	if err != nil {
		s.log.Error("gagal mengambil daftar tindakan laboratorium kategori %s: %v", kategori, err)
		return nil, shared.PaginationMeta{}, err
	}

	for i := range list {
		encId, errEnc := crypto.Encrypt(list[i].KodeTindakan, s.jwtSecret)
		if errEnc != nil {
			s.log.Error("gagal mengenkripsi id tindakan lab %s: %v", list[i].KodeTindakan, errEnc)
			return nil, shared.PaginationMeta{}, fmt.Errorf("gagal memproses data tindakan: %w", errEnc)
		}
		list[i].Id = encId
	}

	meta := shared.NewPaginationMeta(total, filter.Page, filter.Limit)
	return list, meta, nil
}

func (s *service) GetDetailTindakanLab(ctx context.Context, kategori shared.KategoriLab, encryptedId string) (*DetailTindakanLab, error) {
	kodeTindakan, err := crypto.Decrypt(encryptedId, s.jwtSecret)
	if err != nil {
		return nil, apperror.NewBusinessError("ID tindakan lab tidak valid")
	}

	tindakan, templatesDB, err := s.repo.GetDetailTindakanLab(ctx, kategori, kodeTindakan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFoundError("Data tindakan laboratorium tidak ditemukan")
		}
		s.log.Error("gagal mengambil detail tindakan laboratorium %s kategori %s: %v", kodeTindakan, kategori, err)
		return nil, err
	}

	tindakan.Id = encryptedId

	templates := make([]TemplateLab, 0, len(templatesDB))
	for _, t := range templatesDB {
		encTemplateId, errEnc := crypto.Encrypt(strconv.Itoa(t.IdTemplate), s.jwtSecret)
		if errEnc != nil {
			s.log.Error("gagal mengenkripsi id template lab %d: %v", t.IdTemplate, errEnc)
			return nil, fmt.Errorf("gagal memproses data template: %w", errEnc)
		}

		templates = append(templates, TemplateLab{
			IdTemplate:      encTemplateId,
			NamaPemeriksaan: t.NamaPemeriksaan,
			Satuan:          t.Satuan,
			NilaiRujukanLD:  t.NilaiRujukanLD,
			NilaiRujukanLA:  t.NilaiRujukanLA,
			NilaiRujukanPD:  t.NilaiRujukanPD,
			NilaiRujukanPA:  t.NilaiRujukanPA,
		})
	}

	return &DetailTindakanLab{
		TindakanLab: *tindakan,
		Templates:   templates,
	}, nil
}

func (s *service) CekKeberadaanTindakanLab(ctx context.Context, kategori shared.KategoriLab, listKodeTindakan []string) (map[string]bool, error) {
	res, err := s.repo.CekKeberadaanTindakanLab(ctx, kategori, listKodeTindakan)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan tindakan lab di repository: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *service) CekKeberadaanTemplateLab(ctx context.Context, listKodeTindakan []string, templateMap map[string][]int) (map[string]map[int]bool, error) {
	res, err := s.repo.CekKeberadaanTemplateLab(ctx, listKodeTindakan, templateMap)
	if err != nil {
		s.log.Error("Gagal memeriksa keberadaan template lab di repository: %v", err)
		return nil, err
	}
	return res, nil
}

