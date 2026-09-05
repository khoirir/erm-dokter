package radiologi

import (
	"errors"
	"fmt"
	"strings"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type IdHasilRadiologi struct {
	NoRawat        string
	KodeTindakan   string
	TanggalPeriksa string
	JamPeriksa     string
}

func (id IdHasilRadiologi) CompositeKey() string {
	return fmt.Sprintf("%s~%s~%s~%s", id.NoRawat, id.KodeTindakan, id.TanggalPeriksa, id.JamPeriksa)
}

func ParseIdHasilRadiologi(decryptedKey string) (IdHasilRadiologi, error) {
	parts := strings.Split(decryptedKey, "~")
	if len(parts) != 4 {
		return IdHasilRadiologi{}, errors.New("format ID hasil radiologi tidak valid")
	}
	return IdHasilRadiologi{
		NoRawat:        parts[0],
		KodeTindakan:   parts[1],
		TanggalPeriksa: parts[2],
		JamPeriksa:     parts[3],
	}, nil
}

type DokterInfo struct {
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}

type PetugasInfo struct {
	Nip  string `json:"nip"`
	Nama string `json:"nama"`
}

type HasilRadiologi struct {
	Id              string      `json:"id"`
	IdKunjungan     string      `json:"id_kunjungan"`
	NoRawat         string      `json:"no_rawat"`
	KodeTindakan    string      `json:"kode_tindakan"`
	NamaTindakan    string      `json:"nama_tindakan"`
	Status          string      `json:"status"`
	TanggalPeriksa  string      `json:"tanggal_periksa"`
	JamPeriksa      string      `json:"jam_periksa"`
	DokterPerujuk   DokterInfo  `json:"dokter_perujuk"`
	DokterRadiologi DokterInfo  `json:"dokter_radiologi"`
	Petugas         PetugasInfo `json:"petugas"`
	Hasil           string      `json:"hasil"`
	GambarPACS      []string    `json:"gambar_pacs"`
}

func (h *HasilRadiologi) ToIdHasilRadiologi() IdHasilRadiologi {
	return IdHasilRadiologi{
		NoRawat:        h.NoRawat,
		KodeTindakan:   h.KodeTindakan,
		TanggalPeriksa: h.TanggalPeriksa,
		JamPeriksa:     h.JamPeriksa,
	}
}

func (h *HasilRadiologi) CompositeKey() string {
	return h.ToIdHasilRadiologi().CompositeKey()
}

type FilterRiwayatRadiologi struct {
	Tanggal string `json:"tanggal,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

func (f *FilterRiwayatRadiologi) Sanitize() {
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

func (f FilterRiwayatRadiologi) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f FilterRiwayatRadiologi) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if f.Tanggal != "" {
		shared.ValidasiRentangTanggal(f.Tanggal, errs)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
