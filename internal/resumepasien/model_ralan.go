package resumepasien

import (
	"strings"

	"erm-dokter/internal/shared/apperror"
)

type DataResumePasienRalan struct {
	KeluhanUtama                 string        `json:"keluhan_utama"`
	JalannyaPenyakit             string        `json:"jalannya_penyakit"`
	HasilPemeriksaanRadiologi    string        `json:"hasil_pemeriksaan_radiologi"`
	HasilPemeriksaanLaboratorium string        `json:"hasil_pemeriksaan_laboratorium"`
	DiagnosaUtama                string        `json:"diagnosa_utama"`
	KodeDiagnosaUtama            string        `json:"kode_diagnosa_utama,omitempty"`
	DiagnosaSekunder             string        `json:"diagnosa_sekunder"`
	KodeDiagnosaSekunder         string        `json:"kode_diagnosa_sekunder,omitempty"`
	DiagnosaSekunder2            string        `json:"diagnosa_sekunder2"`
	KodeDiagnosaSekunder2        string        `json:"kode_diagnosa_sekunder2,omitempty"`
	DiagnosaSekunder3            string        `json:"diagnosa_sekunder3"`
	KodeDiagnosaSekunder3        string        `json:"kode_diagnosa_sekunder3,omitempty"`
	DiagnosaSekunder4            string        `json:"diagnosa_sekunder4"`
	KodeDiagnosaSekunder4        string        `json:"kode_diagnosa_sekunder4,omitempty"`
	ProsedurUtama                string        `json:"prosedur_utama"`
	KodeProsedurUtama            string        `json:"kode_prosedur_utama,omitempty"`
	ProsedurSekunder             string        `json:"prosedur_sekunder"`
	KodeProsedurSekunder         string        `json:"kode_prosedur_sekunder,omitempty"`
	ProsedurSekunder2            string        `json:"prosedur_sekunder2"`
	KodeProsedurSekunder2        string        `json:"kode_prosedur_sekunder2,omitempty"`
	ProsedurSekunder3            string        `json:"prosedur_sekunder3"`
	KodeProsedurSekunder3        string        `json:"kode_prosedur_sekunder3,omitempty"`
	KeadaanPulang                KeadaanPulang `json:"keadaan_pulang"`
	ObatAtauInstruksi            string        `json:"obat_atau_instruksi"`
}

func (d *DataResumePasienRalan) Sanitize() {
	d.KeluhanUtama = strings.TrimSpace(d.KeluhanUtama)
	d.JalannyaPenyakit = strings.TrimSpace(d.JalannyaPenyakit)
	d.HasilPemeriksaanRadiologi = strings.TrimSpace(d.HasilPemeriksaanRadiologi)
	d.HasilPemeriksaanLaboratorium = strings.TrimSpace(d.HasilPemeriksaanLaboratorium)
	d.DiagnosaUtama = strings.TrimSpace(d.DiagnosaUtama)
	d.KodeDiagnosaUtama = strings.TrimSpace(d.KodeDiagnosaUtama)
	d.DiagnosaSekunder = strings.TrimSpace(d.DiagnosaSekunder)
	d.KodeDiagnosaSekunder = strings.TrimSpace(d.KodeDiagnosaSekunder)
	d.DiagnosaSekunder2 = strings.TrimSpace(d.DiagnosaSekunder2)
	d.KodeDiagnosaSekunder2 = strings.TrimSpace(d.KodeDiagnosaSekunder2)
	d.DiagnosaSekunder3 = strings.TrimSpace(d.DiagnosaSekunder3)
	d.KodeDiagnosaSekunder3 = strings.TrimSpace(d.KodeDiagnosaSekunder3)
	d.DiagnosaSekunder4 = strings.TrimSpace(d.DiagnosaSekunder4)
	d.KodeDiagnosaSekunder4 = strings.TrimSpace(d.KodeDiagnosaSekunder4)
	d.ProsedurUtama = strings.TrimSpace(d.ProsedurUtama)
	d.KodeProsedurUtama = strings.TrimSpace(d.KodeProsedurUtama)
	d.ProsedurSekunder = strings.TrimSpace(d.ProsedurSekunder)
	d.KodeProsedurSekunder = strings.TrimSpace(d.KodeProsedurSekunder)
	d.ProsedurSekunder2 = strings.TrimSpace(d.ProsedurSekunder2)
	d.KodeProsedurSekunder2 = strings.TrimSpace(d.KodeProsedurSekunder2)
	d.ProsedurSekunder3 = strings.TrimSpace(d.ProsedurSekunder3)
	d.KodeProsedurSekunder3 = strings.TrimSpace(d.KodeProsedurSekunder3)
	d.KeadaanPulang = KeadaanPulang(strings.TrimSpace(string(d.KeadaanPulang)))
	if d.KeadaanPulang == "" {
		d.KeadaanPulang = KeadaanPulangHidup
	}
	d.ObatAtauInstruksi = strings.TrimSpace(d.ObatAtauInstruksi)
}

func (d *DataResumePasienRalan) Validate(errs apperror.ValidationError) {
	if d.KeluhanUtama == "" {
		errs["keluhan_utama"] = "Keluhan utama wajib diisi"
	}
	if d.DiagnosaUtama == "" {
		errs["diagnosa_utama"] = "Diagnosa utama wajib diisi"
	}
	if d.KeadaanPulang == "" {
		errs["keadaan_pulang"] = "Keadaan pulang wajib diisi"
	} else if !d.KeadaanPulang.IsValidRalan() {
		errs["keadaan_pulang"] = "Keadaan pulang tidak valid"
	}

	if len(d.KodeDiagnosaUtama) > 10 {
		errs["kode_diagnosa_utama"] = "Kode diagnosa utama maksimal 10 karakter"
	}
	if len(d.KodeDiagnosaSekunder) > 10 {
		errs["kode_diagnosa_sekunder"] = "Kode diagnosa sekunder maksimal 10 karakter"
	}
	if len(d.KodeDiagnosaSekunder2) > 10 {
		errs["kode_diagnosa_sekunder2"] = "Kode diagnosa sekunder 2 maksimal 10 karakter"
	}
	if len(d.DiagnosaSekunder3) > 80 {
		errs["diagnosa_sekunder3"] = "Diagnosa sekunder 3 maksimal 80 karakter"
	}
	if len(d.KodeDiagnosaSekunder3) > 10 {
		errs["kode_diagnosa_sekunder3"] = "Kode diagnosa sekunder 3 maksimal 10 karakter"
	}
	if len(d.DiagnosaSekunder4) > 80 {
		errs["diagnosa_sekunder4"] = "Diagnosa sekunder 4 maksimal 80 karakter"
	}
	if len(d.KodeDiagnosaSekunder4) > 10 {
		errs["kode_diagnosa_sekunder4"] = "Kode diagnosa sekunder 4 maksimal 10 karakter"
	}
	if len(d.ProsedurUtama) > 80 {
		errs["prosedur_utama"] = "Prosedur utama maksimal 80 karakter"
	}
	if len(d.KodeProsedurUtama) > 8 {
		errs["kode_prosedur_utama"] = "Kode prosedur utama maksimal 8 karakter"
	}
	if len(d.ProsedurSekunder) > 80 {
		errs["prosedur_sekunder"] = "Prosedur sekunder 1 maksimal 80 karakter"
	}
	if len(d.KodeProsedurSekunder) > 8 {
		errs["kode_prosedur_sekunder"] = "Kode prosedur sekunder 1 maksimal 8 karakter"
	}
	if len(d.ProsedurSekunder2) > 80 {
		errs["prosedur_sekunder2"] = "Prosedur sekunder 2 maksimal 80 karakter"
	}
	if len(d.KodeProsedurSekunder2) > 8 {
		errs["kode_prosedur_sekunder2"] = "Kode prosedur sekunder 2 maksimal 8 karakter"
	}
	if len(d.ProsedurSekunder3) > 80 {
		errs["prosedur_sekunder3"] = "Prosedur sekunder 3 maksimal 80 karakter"
	}
	if len(d.KodeProsedurSekunder3) > 8 {
		errs["kode_prosedur_sekunder3"] = "Kode prosedur sekunder 3 maksimal 8 karakter"
	}
}

type ResumePasienRalan struct {
	IdKunjungan string `json:"id_kunjungan"`
	NoRawat     string `json:"no_rawat"`
	KodeDokter  string `json:"kode_dokter"`
	NamaDokter  string `json:"nama_dokter"`

	DataResumePasienRalan
}

type SimpanResumePasienRalanRequest struct {
	NoRawat string `json:"no_rawat"`
	DataResumePasienRalan
}

func (req *SimpanResumePasienRalanRequest) Sanitize() {
	req.NoRawat = strings.TrimSpace(req.NoRawat)
	req.DataResumePasienRalan.Sanitize()
}

func (req *SimpanResumePasienRalanRequest) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	if req.NoRawat == "" {
		errs["no_rawat"] = "Nomor rawat wajib diisi"
	}

	req.DataResumePasienRalan.Validate(errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateResumePasienRalanRequest struct {
	DataResumePasienRalan
}

func (req *UpdateResumePasienRalanRequest) Sanitize() {
	req.DataResumePasienRalan.Sanitize()
}

func (req *UpdateResumePasienRalanRequest) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)

	req.DataResumePasienRalan.Validate(errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}
