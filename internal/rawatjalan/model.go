package rawatjalan

import (
	"strings"
	"time"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/shared/formatter"
)

type KunjunganRawatJalan struct {
	Id                string              `json:"id"`
	IdPasien          string              `json:"id_pasien"`
	NoRawat           string              `json:"no_rawat"`
	NoRegistrasi      string              `json:"no_registrasi"`
	TanggalRegistrasi string              `json:"tanggal_registrasi"`
	JamRegistrasi     string              `json:"jam_registrasi"`
	NoRekamMedis      string              `json:"no_rekam_medis"`
	NamaPasien        string              `json:"nama_pasien"`
	JenisKelamin      string              `json:"jenis_kelamin"`
	TanggalLahir      string              `json:"tanggal_lahir"`
	Umur              string              `json:"umur"`
	Alamat            string              `json:"alamat"`
	KodePoliAsal      string              `json:"kode_poli_asal"`
	NamaPoliAsal      string              `json:"nama_poli_asal"`
	KodeDokterAsal    string              `json:"kode_dokter_asal"`
	NamaDokterAsal    string              `json:"nama_dokter_asal"`
	KodePoliRujukan   string              `json:"kode_poli_rujukan,omitempty"`
	NamaPoliRujukan   string              `json:"nama_poli_rujukan,omitempty"`
	KodeDokterRujukan string              `json:"kode_dokter_rujukan,omitempty"`
	NamaDokterRujukan string              `json:"nama_dokter_rujukan,omitempty"`
	KodePenjamin      string              `json:"kode_penjamin"`
	NamaPenjamin      string              `json:"nama_penjamin"`
	StatusPemeriksaan StatusPemeriksaan   `json:"status_pemeriksaan"`
	StatusLanjut      shared.StatusLanjut `json:"status_lanjut"`
	StatusBayar       StatusBayar         `json:"status_bayar"`
	JenisAntrean      JenisAntrean        `json:"jenis_antrean"`
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
	Tanggal           string              `json:"tanggal,omitempty"`
	KodePenjamin      string              `json:"kode_penjamin,omitempty"`
	StatusPemeriksaan StatusPemeriksaan   `json:"status_pemeriksaan,omitempty"`
	JenisAntrean      JenisAntrean        `json:"jenis_antrean,omitempty"`
	StatusLanjut      shared.StatusLanjut `json:"status_lanjut,omitempty"`
	KataKunci         string              `json:"keyword,omitempty"`
	OrderBy           string              `json:"order_by,omitempty"`
	SortOrder         string              `json:"sort_order,omitempty"`
	Halaman           int                 `json:"halaman,omitempty"`
	Batas             int                 `json:"batas,omitempty"`
}

func (f *FilterAntreanDokter) Sanitize() {
	if f.Halaman <= 0 {
		f.Halaman = 1
	}
	if f.Batas <= 0 {
		f.Batas = 20
	} else if f.Batas > 1000 {
		f.Batas = 1000
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

func (f FilterAntreanDokter) Offset() int {
	return (f.Halaman - 1) * f.Batas
}

func (f *FilterAntreanDokter) Validate() apperror.ValidationError {
	f.Sanitize()
	errs := make(apperror.ValidationError)
	if !shared.SortOrder(f.SortOrder).IsValid() {
		errs["sort_order"] = "Jenis pengurutan tidak valid"
	}
	if !OrderBy(f.OrderBy).IsValid() {
		errs["order_by"] = "Jenis pengurutan tidak valid"
	}
	if f.JenisAntrean != "" && !JenisAntrean(f.JenisAntrean).IsValid() {
		errs["jenis_antrean"] = "Jenis antrean tidak valid"
	}
	if f.StatusPemeriksaan != "" && !StatusPemeriksaan(f.StatusPemeriksaan).IsValid() {
		errs["status_pemeriksaan"] = "Status pemeriksaan tidak valid"
	}
	if f.StatusLanjut != "" && !shared.StatusLanjut(f.StatusLanjut).IsValid() {
		errs["status_lanjut"] = "Status lanjut tidak valid"
	}
	if f.Tanggal != "" {
		shared.ValidasiRentangTanggal(f.Tanggal, errs)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
