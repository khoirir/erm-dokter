package rawatjalan

import (
	"strings"
	"time"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/shared/formatter"
)

type KunjunganRawatJalan struct {
	NoRawat           string                   `json:"no_rawat"`
	DetailKey         string                   `json:"detail_key,omitempty"`
	NoRegistrasi      string                   `json:"no_registrasi"`
	TanggalRegistrasi string                   `json:"tanggal_registrasi"`
	JamRegistrasi     string                   `json:"jam_registrasi"`
	NoRekamMedis      string                   `json:"no_rekam_medis"`
	NamaPasien        string                   `json:"nama_pasien"`
	JenisKelamin      string                   `json:"jenis_kelamin"`
	TanggalLahir      string                   `json:"tanggal_lahir"`
	Umur              string                   `json:"umur"`
	Alamat            string                   `json:"alamat"`
	KodePoliAsal      string                   `json:"kode_poli_asal"`
	NamaPoliAsal      string                   `json:"nama_poli_asal"`
	KodeDokterAsal    string                   `json:"kode_dokter_asal"`
	NamaDokterAsal    string                   `json:"nama_dokter_asal"`
	KodePoliRujukan   string                   `json:"kode_poli_rujukan,omitempty"`
	NamaPoliRujukan   string                   `json:"nama_poli_rujukan,omitempty"`
	KodeDokterRujukan string                   `json:"kode_dokter_rujukan,omitempty"`
	NamaDokterRujukan string                   `json:"nama_dokter_rujukan,omitempty"`
	KodePenjamin      string                   `json:"kode_penjamin"`
	NamaPenjamin      string                   `json:"nama_penjamin"`
	StatusPemeriksaan shared.StatusPemeriksaan `json:"status_pemeriksaan"`
	StatusLanjut      shared.StatusLanjut      `json:"status_lanjut"`
	StatusBayar       shared.StatusBayar       `json:"status_bayar"`
	JenisAntrean      shared.JenisAntrean      `json:"jenis_antrean"`
}

func (k *KunjunganRawatJalan) FormatJenisKelamin() string {
	return formatter.FormatJenisKelamin(k.JenisKelamin)
}

func (k *KunjunganRawatJalan) FormatUmur() string {
	return formatter.FormatUmur(k.TanggalLahir)
}

func (k *KunjunganRawatJalan) FormatNoRekamMedis() string {
	return formatter.FormatNoRekamMedis(k.NoRekamMedis)
}

type MetaPaginasi struct {
	TotalData    int `json:"total_data"`
	TotalHalaman int `json:"total_halaman"`
	HalamanAktif int `json:"halaman_aktif"`
	BatasData    int `json:"batas_data"`
}

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
	if strings.TrimSpace(f.Tanggal) == "" {
		today := time.Now().Format("2006-01-02")
		f.Tanggal = today + "," + today
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

func (f *FilterAntreanDokter) Validate() apperror.ValidationError {
	f.Sanitize()
	errs := make(apperror.ValidationError)
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

func validasiTanggal(tanggal string, errs apperror.ValidationError) {
	const layout = "2006-01-02"
	parts := strings.Split(tanggal, ",")

	if len(parts) != 2 {
		errs["tanggal"] = "format tanggal wajib menggunakan rentang 2 tanggal (contoh: YYYY-MM-DD,YYYY-MM-DD)"
		return
	}

	tglAwal, err := time.Parse(layout, strings.TrimSpace(parts[0]))
	if err != nil {
		errs["tanggal"] = "format tanggal awal tidak valid (gunakan YYYY-MM-DD)"
		return
	}

	tglAkhir, err := time.Parse(layout, strings.TrimSpace(parts[1]))
	if err != nil {
		errs["tanggal"] = "format tanggal akhir tidak valid (gunakan YYYY-MM-DD)"
		return
	}

	if tglAwal.After(tglAkhir) {
		errs["tanggal"] = "tanggal awal harus lebih kecil atau sama dengan tanggal akhir"
	}
}
