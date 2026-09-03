package laboratorium

import (
	"fmt"
	"strings"
	"time"

	"erm-dokter/internal/shared/apperror"
)

type ItemHasilLabPK struct {
	IdTemplate      string `json:"id_template,omitempty"`
	NamaPemeriksaan string `json:"nama_pemeriksaan"`
	Nilai           string `json:"nilai"`
	Satuan          string `json:"satuan"`
	NilaiRujukan    string `json:"nilai_rujukan"`
	Keterangan      string `json:"keterangan"`
}

type PermintaanLabPK struct {
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

type DetailTemplateLabPKItem struct {
	IdTemplate      string `json:"id_template"`
	NamaPemeriksaan string `json:"nama_pemeriksaan"`
	Satuan          string `json:"satuan"`
	NilaiRujukan    string `json:"nilai_rujukan"`
	StatusBayar     string `json:"status_bayar"`
}

type PemeriksaanLabPKItem struct {
	IdTindakan     string                    `json:"id_tindakan"`
	KodeTindakan   string                    `json:"kode_tindakan"`
	NamaTindakan   string                    `json:"nama_tindakan"`
	StatusBayar    string                    `json:"status_bayar"`
	DetailTemplate []DetailTemplateLabPKItem `json:"detail_template,omitempty"`
}

type DetailPermintaanLabPK struct {
	PermintaanLabPK
	Pemeriksaan []PemeriksaanLabPKItem `json:"pemeriksaan"`
}

type ItemPemeriksaanLabPKRequest struct {
	IdTindakan   string   `json:"id_tindakan"`
	KodeTindakan string   `json:"-"`
	IdTemplate   []string `json:"id_template,omitempty"`
	KodeTemplate []int    `json:"-"`
}

type SimpanPermintaanLabPKRequest struct {
	NoRawat           string                        `json:"no_rawat"`
	TanggalPermintaan string                        `json:"tanggal_permintaan"`
	JamPermintaan     string                        `json:"jam_permintaan"`
	DiagnosaKlinis    string                        `json:"diagnosa_klinis"`
	InformasiTambahan string                        `json:"informasi_tambahan"`
	Pemeriksaan       []ItemPemeriksaanLabPKRequest `json:"pemeriksaan"`
}

func (req *SimpanPermintaanLabPKRequest) Sanitize() {
	req.NoRawat = strings.TrimSpace(req.NoRawat)
	req.TanggalPermintaan = strings.TrimSpace(req.TanggalPermintaan)
	req.JamPermintaan = strings.TrimSpace(req.JamPermintaan)
	req.DiagnosaKlinis = strings.TrimSpace(req.DiagnosaKlinis)
	req.InformasiTambahan = strings.TrimSpace(req.InformasiTambahan)
	if req.InformasiTambahan == "" {
		req.InformasiTambahan = "-"
	}

	if len(req.JamPermintaan) == 5 {
		if _, err := time.Parse("15:04", req.JamPermintaan); err == nil {
			req.JamPermintaan += ":00"
		}
	}

	for i := range req.Pemeriksaan {
		req.Pemeriksaan[i].IdTindakan = strings.TrimSpace(req.Pemeriksaan[i].IdTindakan)
		for j := range req.Pemeriksaan[i].IdTemplate {
			req.Pemeriksaan[i].IdTemplate[j] = strings.TrimSpace(req.Pemeriksaan[i].IdTemplate[j])
		}
	}
}

func (req *SimpanPermintaanLabPKRequest) Validate() apperror.ValidationError {
	req.Sanitize()
	errs := make(apperror.ValidationError)

	if req.NoRawat == "" {
		errs["no_rawat"] = "Nomor rawat wajib diisi"
	}

	if req.TanggalPermintaan == "" {
		errs["tanggal_permintaan"] = "Tanggal permintaan wajib diisi"
	}
	tgl, errTgl := time.Parse("2006-01-02", req.TanggalPermintaan)
	if req.TanggalPermintaan != "" && errTgl != nil {
		errs["tanggal_permintaan"] = "Format tanggal permintaan harus YYYY-MM-DD"
	}

	if req.JamPermintaan == "" {
		errs["jam_permintaan"] = "Jam permintaan wajib diisi"
	}
	jam, errJam := time.Parse("15:04:05", req.JamPermintaan)
	if req.JamPermintaan != "" && errJam != nil {
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

	if req.DiagnosaKlinis == "" {
		errs["diagnosa_klinis"] = "Diagnosa / indikasi klinis wajib diisi"
	}

	if len(req.Pemeriksaan) == 0 {
		errs["pemeriksaan"] = "Pemeriksaan laboratorium minimal harus memilih 1 tindakan"
	}

	for i, item := range req.Pemeriksaan {
		prefix := fmt.Sprintf("pemeriksaan[%d]", i)
		if item.IdTindakan == "" {
			errs[prefix+".id_tindakan"] = fmt.Sprintf("Pemeriksaan ke-%d: ID tindakan wajib diisi", i+1)
		}
		for j, idTemplate := range item.IdTemplate {
			if idTemplate == "" {
				errs[fmt.Sprintf("%s.id_template[%d]", prefix, j)] = fmt.Sprintf("Pemeriksaan ke-%d parameter ke-%d: ID template pengujian wajib diisi", i+1, j+1)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type KunjunganInfoLabPK struct {
	NoRawat           string
	NoRkmMedis        string
	TanggalRegistrasi string
	JamRegistrasi     string
	KodePenjamin      string
	NamaPenjamin      string
	StatusBayar       string
	StatusLanjut      string
}
