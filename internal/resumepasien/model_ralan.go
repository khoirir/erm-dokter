package resumepasien

import (
	"strings"

	"erm-dokter/internal/shared/apperror"
)

// DataResumePasienRalan memuat seluruh data isian resume pasien rawat jalan (clinical content)
type DataResumePasienRalan struct {
	KeluhanUtama          string        `json:"keluhan_utama"`
	JalannyaPenyakit      string        `json:"jalannya_penyakit"`
	PemeriksaanPenunjang  string        `json:"pemeriksaan_penunjang"`
	HasilLaborat          string        `json:"hasil_laborat"`
	DiagnosaUtama         string        `json:"diagnosa_utama"`
	KodeDiagnosaUtama     string        `json:"kode_diagnosa_utama,omitempty"`
	DiagnosaSekunder      string        `json:"diagnosa_sekunder"`
	KodeDiagnosaSekunder  string        `json:"kode_diagnosa_sekunder,omitempty"`
	DiagnosaSekunder2     string        `json:"diagnosa_sekunder2"`
	KodeDiagnosaSekunder2 string        `json:"kode_diagnosa_sekunder2,omitempty"`
	DiagnosaSekunder3     string        `json:"diagnosa_sekunder3"`
	KodeDiagnosaSekunder3 string        `json:"kode_diagnosa_sekunder3,omitempty"`
	DiagnosaSekunder4     string        `json:"diagnosa_sekunder4"`
	KodeDiagnosaSekunder4 string        `json:"kode_diagnosa_sekunder4,omitempty"`
	ProsedurUtama         string        `json:"prosedur_utama"`
	KodeProsedurUtama     string        `json:"kode_prosedur_utama,omitempty"`
	ProsedurSekunder      string        `json:"prosedur_sekunder"`
	KodeProsedurSekunder  string        `json:"kode_prosedur_sekunder,omitempty"`
	ProsedurSekunder2     string        `json:"prosedur_sekunder2"`
	KodeProsedurSekunder2 string        `json:"kode_prosedur_sekunder2,omitempty"`
	ProsedurSekunder3     string        `json:"prosedur_sekunder3"`
	KodeProsedurSekunder3 string        `json:"kode_prosedur_sekunder3,omitempty"`
	KondisiPulang         KondisiPulang `json:"kondisi_pulang"`
	ObatPulang            string        `json:"obat_pulang"`
}

func (d *DataResumePasienRalan) Sanitize() {
	d.KeluhanUtama = strings.TrimSpace(d.KeluhanUtama)
	d.JalannyaPenyakit = strings.TrimSpace(d.JalannyaPenyakit)
	d.PemeriksaanPenunjang = strings.TrimSpace(d.PemeriksaanPenunjang)
	d.HasilLaborat = strings.TrimSpace(d.HasilLaborat)
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
	d.KondisiPulang = KondisiPulang(strings.TrimSpace(string(d.KondisiPulang)))
	if d.KondisiPulang == "" {
		d.KondisiPulang = KondisiPulangHidup
	}
	d.ObatPulang = strings.TrimSpace(d.ObatPulang)
}

func (d *DataResumePasienRalan) Validate(errs apperror.ValidationError) {
	if d.KeluhanUtama == "" {
		errs["keluhan_utama"] = "Keluhan utama wajib diisi"
	}
	if d.DiagnosaUtama == "" {
		errs["diagnosa_utama"] = "Diagnosa utama wajib diisi"
	}
	if d.KondisiPulang == "" {
		errs["kondisi_pulang"] = "Kondisi pulang wajib diisi"
	} else if !d.KondisiPulang.IsValid() {
		errs["kondisi_pulang"] = "Kondisi pulang harus 'Hidup' atau 'Meninggal'"
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

type DataResumePasien = DataResumePasienRalan

// ResumePasienRalan merepresentasikan data lengkap domain resume pasien rawat jalan untuk response API
type ResumePasienRalan struct {
	IdKunjungan string `json:"id_kunjungan"`
	NoRawat     string `json:"no_rawat"`
	KodeDokter  string `json:"kode_dokter"`
	NamaDokter  string `json:"nama_dokter"`

	DataResumePasienRalan
}

// SimpanResumePasienRalanRequest request payload untuk menyimpan resume pasien rawat jalan baru
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

// UpdateResumePasienRalanRequest request payload untuk memperbarui resume pasien rawat jalan
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

// Type aliases untuk backwards-compatibility
type ResumePasien = ResumePasienRalan
type SimpanResumePasienRequest = SimpanResumePasienRalanRequest
type UpdateResumePasienRequest = UpdateResumePasienRalanRequest

