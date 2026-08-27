package resep

import (
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type Resep struct {
	Id                 string               `json:"id"`
	IdKunjungan        string               `json:"id_kunjungan"`
	NoResep            string               `json:"no_resep"`
	NoRawat            string               `json:"no_rawat"`
	TanggalPeresepan   string               `json:"tanggal_peresepan"`
	JamPeresepan       string               `json:"jam_peresepan"`
	Status             string               `json:"status"`
	KodeDokter         string               `json:"kode_dokter"`
	NamaDokter         string               `json:"nama_dokter"`
	ResepDokter        []ResepDokter        `json:"resep_dokter"`
	ResepDokterRacikan []ResepDokterRacikan `json:"resep_dokter_racikan"`
}

type ResepDokter struct {
	NoResep     string  `json:"-"`
	IdObat      string  `json:"id_obat"`
	KodeObat    string  `json:"kode_obat"`
	NamaObat    string  `json:"nama_obat"`
	Jumlah      float64 `json:"jumlah"`
	Satuan      string  `json:"satuan"`
	AturanPakai string  `json:"aturan_pakai"`
}

type ResepDokterRacikan struct {
	NoResep         string                     `json:"-"`
	NoRacik         string                     `json:"no_racik"`
	NamaRacik       string                     `json:"nama_racik"`
	KodeMetodeRacik string                     `json:"kode_racik"`
	NamaMetodeRacik string                     `json:"metode_racik"`
	JumlahRacikan   int                        `json:"jumlah_racikan"`
	AturanPakai     string                     `json:"aturan_pakai"`
	Keterangan      string                     `json:"keterangan"`
	DetailRacikan   []ResepDokterRacikanDetail `json:"detail_racikan"`
}

type ResepDokterRacikanDetail struct {
	NoResep   string  `json:"-"`
	NoRacik   string  `json:"-"`
	IdObat    string  `json:"id_obat"`
	KodeObat  string  `json:"kode_obat"`
	NamaObat  string  `json:"nama_obat"`
	Kandungan string  `json:"kandungan"`
	Jumlah    float64 `json:"jumlah"`
	Satuan    string  `json:"satuan"`
}

type FilterDaftarResep struct {
	Tanggal string `json:"tanggal,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

func (f *FilterDaftarResep) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 5
	} else if f.Limit > 100 {
		f.Limit = 100
	}
}

func (f FilterDaftarResep) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *FilterDaftarResep) Validate() apperror.ValidationError {
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

type AturanPakai struct {
	AturanPakai string `json:"aturan_pakai"`
}

type MetodeRacik struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}
