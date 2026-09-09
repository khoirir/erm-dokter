package rawatinap

import (
	"strings"
	"time"

	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
	"erm-dokter/internal/shared/formatter"
)

type StatusKamarInap struct {
	IsKamarAktif   bool
	HasRecordKamar bool
}

type OpsiReferensi struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type StatusKunjunganRanap string

const (
	StatusKunjunganBelumKRS   StatusKunjunganRanap = "belum_krs"
	StatusKunjunganTanggalMRS StatusKunjunganRanap = "tanggal_mrs"
	StatusKunjunganTanggalKRS StatusKunjunganRanap = "tanggal_krs"
)

func (s StatusKunjunganRanap) IsValid() bool {
	switch s {
	case StatusKunjunganBelumKRS, StatusKunjunganTanggalMRS, StatusKunjunganTanggalKRS:
		return true
	default:
		return false
	}
}

type ScopeDPJP string

const (
	ScopeDPJPDokter ScopeDPJP = "dpjp"
	ScopeDPJPSemua  ScopeDPJP = "semua"
)

func (s ScopeDPJP) IsValid() bool {
	switch s {
	case ScopeDPJPDokter, ScopeDPJPSemua:
		return true
	default:
		return false
	}
}

type OrderByRanap string

const (
	OrderByWaktuMRS     OrderByRanap = "waktu_mrs"
	OrderByWaktuKRS     OrderByRanap = "waktu_krs"
	OrderByNamaPasien   OrderByRanap = "nama_pasien"
	OrderByKamar        OrderByRanap = "kamar"
	OrderByKelas        OrderByRanap = "kelas"
	OrderByBangsal      OrderByRanap = "bangsal"
	OrderByPenjamin     OrderByRanap = "penjamin"
	OrderByStatusPulang OrderByRanap = "status_pulang"
)

func (o OrderByRanap) IsValid() bool {
	switch o {
	case OrderByWaktuMRS, OrderByWaktuKRS, OrderByNamaPasien,
		OrderByKamar, OrderByKelas, OrderByBangsal, OrderByPenjamin, OrderByStatusPulang:
		return true
	default:
		return false
	}
}

type KunjunganRawatInap struct {
	Id                   string              `json:"id"`
	IdKunjungan          string              `json:"id_kunjungan"`
	IdPasien             string              `json:"id_pasien"`
	NoRawat              string              `json:"no_rawat"`
	NoRegistrasi         string              `json:"no_registrasi"`
	TanggalRegistrasi    string              `json:"tanggal_registrasi"`
	JamRegistrasi        string              `json:"jam_registrasi"`
	NoRekamMedis         string              `json:"no_rekam_medis"`
	NamaPasien           string              `json:"nama_pasien"`
	JenisKelamin         string              `json:"jenis_kelamin"`
	TanggalLahir         string              `json:"tanggal_lahir"`
	Umur                 string              `json:"umur"`
	Alamat               string              `json:"alamat"`
	KodeKamar            string              `json:"kode_kamar"`
	KelasKamar           string              `json:"kelas"`
	KodeBangsal          string              `json:"kode_bangsal"`
	NamaBangsal          string              `json:"nama_bangsal"`
	TanggalMasuk         string              `json:"tanggal_masuk"`
	JamMasuk             string              `json:"jam_masuk"`
	TanggalKeluar        string              `json:"tanggal_keluar,omitempty"`
	JamKeluar            string              `json:"jam_keluar,omitempty"`
	StatusPulang         StatusPulang        `json:"status_pulang"`
	LamaInap             int                 `json:"lama_inap"`
	KodePenjamin         string              `json:"kode_penjamin"`
	NamaPenjamin         string              `json:"nama_penjamin"`
	StatusBayar          string              `json:"status_bayar"`
	StatusLanjut         shared.StatusLanjut `json:"status_lanjut"`
	KodeDokterRawatJalan string              `json:"kode_dokter_rawat_jalan,omitempty"`
	DokterRawatJalan     string              `json:"dokter_rawat_jalan"`
	DPJP                 []string            `json:"dokter_rawat_inap"`
	DiagnosaAwal         string              `json:"diagnosa_awal"`
	DiagnosaAkhir        string              `json:"diagnosa_akhir,omitempty"`
	GolonganDarah        string              `json:"golongan_darah"`
	Agama                string              `json:"agama"`
	NoTelepon            string              `json:"no_telepon"`
	NoPeserta            string              `json:"no_peserta"`
	NoKTP                string              `json:"no_ktp"`
	PenanggungJawab      string              `json:"penanggung_jawab"`
	HubunganPenanggungJawab string           `json:"hubungan_penanggung_jawab"`
	AlamatPenanggungJawab string             `json:"alamat_penanggung_jawab"`
}

func (k *KunjunganRawatInap) FormatJenisKelamin() string {
	return formatter.FormatJenisKelamin(k.JenisKelamin)
}

func (k *KunjunganRawatInap) FormatUmur() string {
	return formatter.FormatUmur(k.TanggalLahir)
}

func (k *KunjunganRawatInap) FormatNoRekamMedis() string {
	return formatter.FormatNoRekamMedis(k.NoRekamMedis)
}

type FilterPasienRawatInap struct {
	Bangsal         string               `json:"bangsal,omitempty"`
	Kelas           string               `json:"kelas,omitempty"`
	Penjamin        string               `json:"penjamin,omitempty"`
	ScopeDPJP       ScopeDPJP            `json:"scope_dpjp,omitempty"`
	StatusKunjungan StatusKunjunganRanap `json:"status_kunjungan,omitempty"`
	StatusPulang    StatusPulang         `json:"status_pulang,omitempty"`
	Tanggal         string               `json:"tanggal,omitempty"`
	Keyword         string               `json:"keyword,omitempty"`
	OrderBy         OrderByRanap         `json:"order_by,omitempty"`
	SortOrder       string               `json:"sort_order,omitempty"`
	Page            int                  `json:"page,omitempty"`
	Limit           int                  `json:"limit,omitempty"`
}

func (f *FilterPasienRawatInap) Sanitize() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	} else if f.Limit > 100 {
		f.Limit = 100
	}
	if f.StatusKunjungan == "" {
		f.StatusKunjungan = StatusKunjunganBelumKRS
	}
	if f.ScopeDPJP == "" {
		f.ScopeDPJP = ScopeDPJPSemua
	}
	if f.OrderBy == "" {
		f.OrderBy = OrderByWaktuMRS
	}
	if f.SortOrder == "" {
		f.SortOrder = "DESC"
	} else {
		f.SortOrder = strings.ToUpper(strings.TrimSpace(f.SortOrder))
	}
	f.StatusPulang = StatusPulang(strings.TrimSpace(string(f.StatusPulang)))
	f.Bangsal = strings.TrimSpace(f.Bangsal)
	f.Kelas = strings.TrimSpace(f.Kelas)
	f.Penjamin = strings.TrimSpace(f.Penjamin)
	f.Tanggal = strings.TrimSpace(f.Tanggal)
	f.Keyword = strings.TrimSpace(f.Keyword)

	if f.Tanggal == "" && (f.StatusKunjungan == StatusKunjunganTanggalMRS || f.StatusKunjungan == StatusKunjunganTanggalKRS) {
		today := time.Now().Format("2006-01-02")
		f.Tanggal = today + "," + today
	}
}

func (f FilterPasienRawatInap) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f FilterPasienRawatInap) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)
	if !shared.SortOrder(f.SortOrder).IsValid() {
		errs["sort_order"] = "Jenis pengurutan tidak valid"
	}
	if !f.OrderBy.IsValid() {
		errs["order_by"] = "Kolom pengurutan tidak valid"
	}
	if !f.StatusKunjungan.IsValid() {
		errs["status_kunjungan"] = "Status kunjungan tidak valid (pilihan: belum_krs, tanggal_mrs, tanggal_krs)"
	}
	if !f.ScopeDPJP.IsValid() {
		errs["scope_dpjp"] = "Scope DPJP tidak valid (pilihan: dpjp, semua)"
	}
	if f.StatusPulang != "" && !f.StatusPulang.IsValid() {
		errs["status_pulang"] = "Status pulang tidak valid"
	}
	if f.Tanggal != "" {
		shared.ValidasiRentangTanggal(f.Tanggal, errs)
	}
	if len(f.Keyword) > 0 && len(f.Keyword) < 3 {
		errs["keyword"] = "Kata kunci pencarian minimal 3 karakter"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
