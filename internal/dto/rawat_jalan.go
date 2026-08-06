package dto

import (
	"strings"
	"time"
)

type OpsiReferensi struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ReferensiFilterRawatJalan struct {
	StatusPemeriksaan []OpsiReferensi `json:"status_pemeriksaan"`
	StatusLanjut      []OpsiReferensi `json:"status_lanjut"`
	StatusBayar       []OpsiReferensi `json:"status_bayar"`
	JenisAntrean      []OpsiReferensi `json:"jenis_antrean"`
}

type FilterAntreanDokter struct {
	KodeDokter        string `json:"kode_dokter"`
	Tanggal           string `json:"tanggal"`
	KodePenjamin      string `json:"kode_penjamin"`
	StatusPemeriksaan string `json:"status_pemeriksaan"`
	JenisAntrean      string `json:"jenis_antrean"`
	KataKunci         string `json:"keyword"`
	OrderBy           string `json:"order_by"`
	SortOrder         string `json:"sort_order"`
	Halaman           int    `json:"halaman"`
	Batas             int    `json:"batas"`
}

func (f *FilterAntreanDokter) Sanitize() {
	if f.Halaman <= 0 {
		f.Halaman = 1
	}
	if f.Batas <= 0 {
		f.Batas = 10
	} else if f.Batas > 2000 {
		f.Batas = 2000
	}
	if f.OrderBy == "" {
		f.OrderBy = "waktu_registrasi"
	}
	if f.SortOrder == "" {
		f.SortOrder = "ASC"
	} else {
		f.SortOrder = strings.ToUpper(strings.TrimSpace(f.SortOrder))
	}
}

var validStatusPemeriksaan = map[string]bool{
	"Belum": true, "Sudah": true, "Batal": true, "Berkas Diterima": true,
	"Dirujuk": true, "Meninggal": true, "Dirawat": true, "Pulang Paksa": true,
}

var validJenisAntrean = map[string]bool{
	"Rujukan": true, "Bukan Rujukan": true,
}

var validOrderBy = map[string]bool{
	"waktu_registrasi": true, "nama_pasien": true,
}

var validSortOrder = map[string]bool{
	"ASC": true, "DESC": true,
}

func (f *FilterAntreanDokter) Validate() ValidationError {
	f.Sanitize()
	errs := make(ValidationError)
	if !validSortOrder[f.SortOrder] {
		errs["sort_order"] = "sort_order harus 'ASC' atau 'DESC'"
	}
	if !validOrderBy[f.OrderBy] {
		errs["order_by"] = "order_by hanya boleh 'waktu_registrasi' atau 'nama_pasien'"
	}
	if f.JenisAntrean != "" && !validJenisAntrean[f.JenisAntrean] {
		errs["jenis_antrean"] = "jenis_antrean tidak valid (hanya 'Rujukan' atau 'Bukan Rujukan')"
	}
	if f.StatusPemeriksaan != "" && !validStatusPemeriksaan[f.StatusPemeriksaan] {
		errs["status_pemeriksaan"] = "status_pemeriksaan tidak valid"
	}
	if f.Tanggal != "" {
		validasiTanggal(f.Tanggal, errs)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validasiTanggal(tanggal string, errs ValidationError) {
	const layout = "2006-01-02"
	parts := strings.Split(tanggal, ",")

	if len(parts) > 2 {
		errs["tanggal"] = "format tanggal tidak valid (gunakan YYYY-MM-DD atau YYYY-MM-DD,YYYY-MM-DD)"
		return
	}

	tglAwal, err := time.Parse(layout, strings.TrimSpace(parts[0]))
	if err != nil {
		errs["tanggal"] = "format tanggal awal tidak valid (gunakan YYYY-MM-DD)"
		return
	}

	if len(parts) == 2 {
		tglAkhir, err := time.Parse(layout, strings.TrimSpace(parts[1]))
		if err != nil {
			errs["tanggal"] = "format tanggal akhir tidak valid (gunakan YYYY-MM-DD)"
			return
		}
		if tglAwal.After(tglAkhir) {
			errs["tanggal"] = "tanggal awal harus lebih kecil atau sama dengan tanggal akhir"
		}
	}
}
