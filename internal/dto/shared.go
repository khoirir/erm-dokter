package dto

import "strings"

type MetaPaginasi struct {
	TotalData    int `json:"total_data"`
	TotalHalaman int `json:"total_halaman"`
	HalamanAktif int `json:"halaman_aktif"`
	BatasData    int `json:"batas_data"`
}

type ValidationError map[string]string

func (v ValidationError) Error() string {
	var errs []string
	for field, msg := range v {
		errs = append(errs, field+": "+msg)
	}
	return strings.Join(errs, "; ")
}
