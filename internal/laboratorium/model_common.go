package laboratorium

import (
	"strings"
	"time"

	"erm-dokter/internal/berkasdigital"
	"erm-dokter/internal/shared"
	"erm-dokter/internal/shared/apperror"
)

type PermintaanLabHeader struct {
	Id                string `json:"id"`
	NoPermintaan      string `json:"no_permintaan"`
	IdKunjungan       string `json:"id_kunjungan"`
	NoRawat           string `json:"no_rawat"`
	TanggalPermintaan string `json:"tanggal_permintaan"`
	JamPermintaan     string `json:"jam_permintaan"`
	TanggalSampel     string `json:"tanggal_sampel,omitempty"`
	JamSampel         string `json:"jam_sampel,omitempty"`
	TanggalHasil      string `json:"tanggal_hasil,omitempty"`
	JamHasil          string `json:"jam_hasil,omitempty"`
	KodeDokterPerujuk string `json:"kode_dokter_perujuk"`
	NamaDokterPerujuk string `json:"nama_dokter_perujuk"`
	Status            string `json:"status"`
	InformasiTambahan string `json:"informasi_tambahan"`
	DiagnosaKlinis    string `json:"diagnosa_klinis"`
	StatusProses      string `json:"status_proses"`
}

func (h *PermintaanLabHeader) FormatZeroDates() {
	if h.TanggalHasil != "0000-00-00" && strings.TrimSpace(h.TanggalHasil) != "" {
		h.StatusProses = "Selesai"
	} else if h.TanggalSampel != "0000-00-00" && strings.TrimSpace(h.TanggalSampel) != "" {
		h.StatusProses = "Sampel Diambil"
	} else {
		h.StatusProses = "Menunggu Sampel"
	}

	if h.TanggalSampel == "0000-00-00" {
		h.TanggalSampel = ""
	}
	if h.JamSampel == "00:00:00" {
		h.JamSampel = ""
	}
	if h.TanggalHasil == "0000-00-00" {
		h.TanggalHasil = ""
	}
	if h.JamHasil == "00:00:00" {
		h.JamHasil = ""
	}
}

type PermintaanLabHeaderRequest struct {
	NoRawat           string `json:"no_rawat"`
	TanggalPermintaan string `json:"tanggal_permintaan"`
	JamPermintaan     string `json:"jam_permintaan"`
	DiagnosaKlinis    string `json:"diagnosa_klinis"`
	InformasiTambahan string `json:"informasi_tambahan"`
}

func (h *PermintaanLabHeaderRequest) Sanitize() {
	h.NoRawat = strings.TrimSpace(h.NoRawat)
	h.TanggalPermintaan = strings.TrimSpace(h.TanggalPermintaan)
	h.JamPermintaan = strings.TrimSpace(h.JamPermintaan)
	h.DiagnosaKlinis = strings.TrimSpace(h.DiagnosaKlinis)
	h.InformasiTambahan = strings.TrimSpace(h.InformasiTambahan)
	if h.InformasiTambahan == "" {
		h.InformasiTambahan = "-"
	}

	if len(h.JamPermintaan) == 5 {
		if _, err := time.Parse("15:04", h.JamPermintaan); err == nil {
			h.JamPermintaan += ":00"
		}
	}
}

func (h PermintaanLabHeaderRequest) Validate(errs apperror.ValidationError) {
	if h.NoRawat == "" {
		errs["no_rawat"] = "Nomor rawat wajib diisi"
	}

	if h.TanggalPermintaan == "" {
		errs["tanggal_permintaan"] = "Tanggal permintaan wajib diisi"
	}
	tgl, errTgl := time.Parse("2006-01-02", h.TanggalPermintaan)
	if h.TanggalPermintaan != "" && errTgl != nil {
		errs["tanggal_permintaan"] = "Format tanggal permintaan harus YYYY-MM-DD"
	}

	if h.JamPermintaan == "" {
		errs["jam_permintaan"] = "Jam permintaan wajib diisi"
	}
	jam, errJam := time.Parse("15:04:05", h.JamPermintaan)
	if h.JamPermintaan != "" && errJam != nil {
		errs["jam_permintaan"] = "Format jam permintaan harus HH:mm:ss"
	}

	if errTgl == nil && errJam == nil {
		waktuPermintaan := time.Date(
			tgl.Year(), tgl.Month(), tgl.Day(),
			jam.Hour(), jam.Minute(), jam.Second(), 0,
			time.Local,
		)
		if waktuPermintaan.After(time.Now().Add(5 * time.Minute)) {
			errs["tanggal_permintaan"] = "Waktu permintaan laboratorium tidak boleh melebihi waktu saat ini"
		}
	}

	if h.DiagnosaKlinis == "" {
		errs["diagnosa_klinis"] = "Diagnosa / indikasi klinis wajib diisi"
	} else if len(h.DiagnosaKlinis) > 80 {
		errs["diagnosa_klinis"] = "Diagnosa klinis maksimal 80 karakter"
	}
	if len(h.InformasiTambahan) > 60 {
		errs["informasi_tambahan"] = "Informasi tambahan maksimal 60 karakter"
	}
}

type DokterInfo struct {
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}

type PetugasInfo struct {
	Nip  string `json:"nip"`
	Nama string `json:"nama"`
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
	HasilPemeriksaan []HasilLaboratorium           `json:"hasil_pemeriksaan"`
	BerkasDigital    []berkasdigital.BerkasDigital `json:"berkas_digital"`
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

func (f FilterRiwayatLab) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if f.Tanggal != "" {
		shared.ValidasiRentangTanggal(f.Tanggal, errs)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
