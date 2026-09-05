package tindakan

import (
	"strings"

	"erm-dokter/internal/shared/apperror"
)

type TindakanRadiologi struct {
	Id           string  `json:"id"`
	KodeTindakan string  `json:"kode_tindakan"`
	NamaTindakan string  `json:"nama_tindakan"`
	Biaya        float64 `json:"biaya"`
}

type FilterDaftarTindakanRadiologi struct {
	Keyword string `json:"keyword,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

func (f *FilterDaftarTindakanRadiologi) Sanitize() {
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

func (f FilterDaftarTindakanRadiologi) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f FilterDaftarTindakanRadiologi) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if f.Keyword != "" && len(f.Keyword) < 3 {
		errs["keyword"] = "Pencarian tindakan radiologi minimal 3 karakter"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
