package laboratorium

import (
	"strings"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type DokterInfo struct {
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}

type PetugasInfo struct {
	Nip  string `json:"nip"`
	Nama string `json:"nama"`
}

type BerkasDigital struct {
	Kode       string `json:"kode"`
	NamaBerkas string `json:"nama_berkas"`
	IdBerkas   string `json:"id_berkas"`
	UrlBerkas  string `json:"url_berkas"`
}

type HasilLaboratorium struct {
	Id             string           `json:"id"`
	IdKunjungan    string           `json:"id_kunjungan"`
	NoRawat        string           `json:"no_rawat"`
	KodeTindakan   string           `json:"kode_tindakan"`
	NamaTindakan   string           `json:"nama_tindakan"`
	Kategori       string           `json:"kategori"`
	Status         string           `json:"status"`
	TanggalPeriksa string           `json:"tanggal_periksa"`
	JamPeriksa     string           `json:"jam_periksa"`
	DokterPerujuk  DokterInfo       `json:"dokter_perujuk"`
	DokterPJ       DokterInfo       `json:"dokter_pj"`
	Petugas        PetugasInfo      `json:"petugas"`
	DetailPK       []ItemHasilLabPK `json:"detail_pk,omitempty"`
	DetailPA       *HasilLabPA      `json:"detail_pa,omitempty"`
}

type HasilLaboratoriumKunjungan struct {
	HasilPemeriksaan []HasilLaboratorium `json:"hasil_pemeriksaan"`
	BerkasDigital    []BerkasDigital     `json:"berkas_digital"`
}

type FilterRiwayatLab struct {
	Tanggal string `json:"tanggal,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

func (f *FilterRiwayatLab) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 5
	} else if f.Limit > 100 {
		f.Limit = 100
	}
	f.Tanggal = strings.TrimSpace(f.Tanggal)
}

func (f FilterRiwayatLab) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *FilterRiwayatLab) Validate() apperror.ValidationError {
	f.Sanitize()
	errs := make(apperror.ValidationError)

	if f.Tanggal != "" {
		shared.ValidasiRentangTanggal(f.Tanggal, errs)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
