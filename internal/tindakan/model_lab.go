package tindakan

import (
	"strings"

	"erm-dokter/internal/shared/apperror"
)

type TindakanLab struct {
	Id           string  `json:"id"`
	KodeTindakan string  `json:"kode_tindakan"`
	NamaTindakan string  `json:"nama_tindakan"`
	Biaya        float64 `json:"biaya"`
}

type TemplateLab struct {
	IdTemplate      string `json:"id_template"`
	NamaPemeriksaan string `json:"nama_pemeriksaan"`
	Satuan          string `json:"satuan"`
	NilaiRujukanLD  string `json:"nilai_rujukan_ld"`
	NilaiRujukanLA  string `json:"nilai_rujukan_la"`
	NilaiRujukanPD  string `json:"nilai_rujukan_pd"`
	NilaiRujukanPA  string `json:"nilai_rujukan_pa"`
}

type DetailTindakanLab struct {
	TindakanLab
	Templates []TemplateLab `json:"templates"`
}

type FilterDaftarTindakanLab struct {
	Keyword string `json:"keyword,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

func (f *FilterDaftarTindakanLab) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	} else if f.Limit > 100 {
		f.Limit = 100
	}
	f.Keyword = strings.TrimSpace(f.Keyword)
}

func (f FilterDaftarTindakanLab) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f FilterDaftarTindakanLab) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if f.Keyword != "" && len(f.Keyword) < 3 {
		errs["keyword"] = "Pencarian tindakan laboratorium minimal 3 karakter"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
