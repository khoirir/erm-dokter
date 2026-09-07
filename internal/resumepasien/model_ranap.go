package resumepasien

import (
	"strings"
	"time"

	"erm-dokter/internal/shared/apperror"
)

type InfoKamarInap struct {
	TanggalMasuk  string `json:"tanggal_masuk,omitempty"`
	JamMasuk      string `json:"jam_masuk,omitempty"`
	TanggalKeluar string `json:"tanggal_keluar,omitempty"`
	JamKeluar     string `json:"jam_keluar,omitempty"`
	Kamar         string `json:"kamar,omitempty"`
}

type DataResumePasienRanap struct {
	// Riwayat Masuk
	DiagnosaAwal string `json:"diagnosa_awal"`
	Alasan       string `json:"alasan"`

	// Riwayat & Pemeriksaan Klinis Selama Rawat Inap
	KeluhanUtama         string `json:"keluhan_utama"`
	PemeriksaanFisik     string `json:"pemeriksaan_fisik"`
	JalannyaPenyakit     string `json:"jalannya_penyakit"`
	PemeriksaanPenunjang string `json:"pemeriksaan_penunjang"`
	HasilLaborat         string `json:"hasil_laborat"`
	TindakanDanOperasi   string `json:"tindakan_dan_operasi"`
	ObatDiRS             string `json:"obat_di_rs"`

	// Diagnosa Akhir & ICD 10
	DiagnosaUtama       string `json:"diagnosa_utama"`
	KdDiagnosaUtama     string `json:"kd_diagnosa_utama,omitempty"`
	DiagnosaSekunder    string `json:"diagnosa_sekunder"`
	KdDiagnosaSekunder  string `json:"kd_diagnosa_sekunder,omitempty"`
	DiagnosaSekunder2   string `json:"diagnosa_sekunder2"`
	KdDiagnosaSekunder2 string `json:"kd_diagnosa_sekunder2,omitempty"`
	DiagnosaSekunder3   string `json:"diagnosa_sekunder3"`
	KdDiagnosaSekunder3 string `json:"kd_diagnosa_sekunder3,omitempty"`
	DiagnosaSekunder4   string `json:"diagnosa_sekunder4"`
	KdDiagnosaSekunder4 string `json:"kd_diagnosa_sekunder4,omitempty"`

	// Prosedur & ICD 9
	ProsedurUtama       string `json:"prosedur_utama"`
	KdProsedurUtama     string `json:"kd_prosedur_utama,omitempty"`
	ProsedurSekunder    string `json:"prosedur_sekunder"`
	KdProsedurSekunder  string `json:"kd_prosedur_sekunder,omitempty"`
	ProsedurSekunder2   string `json:"prosedur_sekunder2"`
	KdProsedurSekunder2 string `json:"kd_prosedur_sekunder2,omitempty"`
	ProsedurSekunder3   string `json:"prosedur_sekunder3"`
	KdProsedurSekunder3 string `json:"kd_prosedur_sekunder3,omitempty"`

	// Kondisi Khusus & Edukasi
	Alergi   string `json:"alergi"`
	Diet     string `json:"diet"`
	LabBelum string `json:"lab_belum"`
	Edukasi  string `json:"edukasi"`

	// Pemulangan Pasien
	CaraKeluar     CaraKeluar    `json:"cara_keluar"`
	KetKeluar      string        `json:"ket_keluar"`
	Keadaan        KeadaanPulang `json:"keadaan"`
	KetKeadaan     string        `json:"ket_keadaan"`
	Dilanjutkan    Dilanjutkan   `json:"dilanjutkan"`
	KetDilanjutkan string        `json:"ket_dilanjutkan"`
	Kontrol        string        `json:"kontrol"`
	ObatPulang     string        `json:"obat_pulang"`
}

func (d *DataResumePasienRanap) Sanitize() {
	d.DiagnosaAwal = strings.TrimSpace(d.DiagnosaAwal)
	d.Alasan = strings.TrimSpace(d.Alasan)
	d.KeluhanUtama = strings.TrimSpace(d.KeluhanUtama)
	d.PemeriksaanFisik = strings.TrimSpace(d.PemeriksaanFisik)
	d.JalannyaPenyakit = strings.TrimSpace(d.JalannyaPenyakit)
	d.PemeriksaanPenunjang = strings.TrimSpace(d.PemeriksaanPenunjang)
	d.HasilLaborat = strings.TrimSpace(d.HasilLaborat)
	d.TindakanDanOperasi = strings.TrimSpace(d.TindakanDanOperasi)
	d.ObatDiRS = strings.TrimSpace(d.ObatDiRS)

	d.DiagnosaUtama = strings.TrimSpace(d.DiagnosaUtama)
	d.KdDiagnosaUtama = strings.TrimSpace(d.KdDiagnosaUtama)
	d.DiagnosaSekunder = strings.TrimSpace(d.DiagnosaSekunder)
	d.KdDiagnosaSekunder = strings.TrimSpace(d.KdDiagnosaSekunder)
	d.DiagnosaSekunder2 = strings.TrimSpace(d.DiagnosaSekunder2)
	d.KdDiagnosaSekunder2 = strings.TrimSpace(d.KdDiagnosaSekunder2)
	d.DiagnosaSekunder3 = strings.TrimSpace(d.DiagnosaSekunder3)
	d.KdDiagnosaSekunder3 = strings.TrimSpace(d.KdDiagnosaSekunder3)
	d.DiagnosaSekunder4 = strings.TrimSpace(d.DiagnosaSekunder4)
	d.KdDiagnosaSekunder4 = strings.TrimSpace(d.KdDiagnosaSekunder4)

	d.ProsedurUtama = strings.TrimSpace(d.ProsedurUtama)
	d.KdProsedurUtama = strings.TrimSpace(d.KdProsedurUtama)
	d.ProsedurSekunder = strings.TrimSpace(d.ProsedurSekunder)
	d.KdProsedurSekunder = strings.TrimSpace(d.KdProsedurSekunder)
	d.ProsedurSekunder2 = strings.TrimSpace(d.ProsedurSekunder2)
	d.KdProsedurSekunder2 = strings.TrimSpace(d.KdProsedurSekunder2)
	d.ProsedurSekunder3 = strings.TrimSpace(d.ProsedurSekunder3)
	d.KdProsedurSekunder3 = strings.TrimSpace(d.KdProsedurSekunder3)

	d.Alergi = strings.TrimSpace(d.Alergi)
	d.Diet = strings.TrimSpace(d.Diet)
	d.LabBelum = strings.TrimSpace(d.LabBelum)
	d.Edukasi = strings.TrimSpace(d.Edukasi)

	d.CaraKeluar = CaraKeluar(strings.TrimSpace(string(d.CaraKeluar)))
	d.KetKeluar = strings.TrimSpace(d.KetKeluar)
	d.Keadaan = KeadaanPulang(strings.TrimSpace(string(d.Keadaan)))
	d.KetKeadaan = strings.TrimSpace(d.KetKeadaan)
	d.Dilanjutkan = Dilanjutkan(strings.TrimSpace(string(d.Dilanjutkan)))
	d.KetDilanjutkan = strings.TrimSpace(d.KetDilanjutkan)

	d.Kontrol = strings.TrimSpace(d.Kontrol)
	if d.Kontrol == "-" {
		d.Kontrol = ""
	}
	if len(d.Kontrol) == 16 { // YYYY-MM-DD HH:mm
		d.Kontrol += ":00"
	}
	d.ObatPulang = strings.TrimSpace(d.ObatPulang)
}

func (d *DataResumePasienRanap) Validate(errs apperror.ValidationError) {
	if d.DiagnosaAwal == "" {
		errs["diagnosa_awal"] = "Diagnosa awal masuk wajib diisi"
	} else if len(d.DiagnosaAwal) > 500 {
		errs["diagnosa_awal"] = "Diagnosa awal masuk maksimal 500 karakter"
	}

	if d.Alasan == "" {
		errs["alasan"] = "Alasan masuk dirawat wajib diisi"
	} else if len(d.Alasan) > 100 {
		errs["alasan"] = "Alasan masuk dirawat maksimal 100 karakter"
	}

	if d.KeluhanUtama == "" {
		errs["keluhan_utama"] = "Keluhan utama wajib diisi"
	}

	if d.DiagnosaUtama == "" {
		errs["diagnosa_utama"] = "Diagnosa utama wajib diisi"
	} else if len(d.DiagnosaUtama) > 150 {
		errs["diagnosa_utama"] = "Diagnosa utama maksimal 150 karakter"
	}

	if len(d.KdDiagnosaUtama) > 10 {
		errs["kd_diagnosa_utama"] = "Kode diagnosa utama maksimal 10 karakter"
	}
	if len(d.DiagnosaSekunder) > 150 {
		errs["diagnosa_sekunder"] = "Diagnosa sekunder 1 maksimal 150 karakter"
	}
	if len(d.KdDiagnosaSekunder) > 10 {
		errs["kd_diagnosa_sekunder"] = "Kode diagnosa sekunder 1 maksimal 10 karakter"
	}
	if len(d.DiagnosaSekunder2) > 150 {
		errs["diagnosa_sekunder2"] = "Diagnosa sekunder 2 maksimal 150 karakter"
	}
	if len(d.KdDiagnosaSekunder2) > 10 {
		errs["kd_diagnosa_sekunder2"] = "Kode diagnosa sekunder 2 maksimal 10 karakter"
	}
	if len(d.DiagnosaSekunder3) > 150 {
		errs["diagnosa_sekunder3"] = "Diagnosa sekunder 3 maksimal 150 karakter"
	}
	if len(d.KdDiagnosaSekunder3) > 10 {
		errs["kd_diagnosa_sekunder3"] = "Kode diagnosa sekunder 3 maksimal 10 karakter"
	}
	if len(d.DiagnosaSekunder4) > 150 {
		errs["diagnosa_sekunder4"] = "Diagnosa sekunder 4 maksimal 150 karakter"
	}
	if len(d.KdDiagnosaSekunder4) > 10 {
		errs["kd_diagnosa_sekunder4"] = "Kode diagnosa sekunder 4 maksimal 10 karakter"
	}

	if len(d.ProsedurUtama) > 150 {
		errs["prosedur_utama"] = "Prosedur utama maksimal 150 karakter"
	}
	if len(d.KdProsedurUtama) > 8 {
		errs["kd_prosedur_utama"] = "Kode prosedur utama maksimal 8 karakter"
	}
	if len(d.ProsedurSekunder) > 150 {
		errs["prosedur_sekunder"] = "Prosedur sekunder 1 maksimal 150 karakter"
	}
	if len(d.KdProsedurSekunder) > 8 {
		errs["kd_prosedur_sekunder"] = "Kode prosedur sekunder 1 maksimal 8 karakter"
	}
	if len(d.ProsedurSekunder2) > 150 {
		errs["prosedur_sekunder2"] = "Prosedur sekunder 2 maksimal 150 karakter"
	}
	if len(d.KdProsedurSekunder2) > 8 {
		errs["kd_prosedur_sekunder2"] = "Kode prosedur sekunder 2 maksimal 8 karakter"
	}
	if len(d.ProsedurSekunder3) > 150 {
		errs["prosedur_sekunder3"] = "Prosedur sekunder 3 maksimal 150 karakter"
	}
	if len(d.KdProsedurSekunder3) > 8 {
		errs["kd_prosedur_sekunder3"] = "Kode prosedur sekunder 3 maksimal 8 karakter"
	}

	if len(d.Alergi) > 100 {
		errs["alergi"] = "Alergi maksimal 100 karakter"
	}

	if d.CaraKeluar == "" {
		errs["cara_keluar"] = "Cara keluar wajib diisi"
	} else if !d.CaraKeluar.IsValid() {
		errs["cara_keluar"] = "Cara keluar tidak valid (pilihan: Atas Izin Dokter, Pindah RS, Pulang Atas Permintaan Sendiri, Lainnya)"
	}
	if len(d.KetKeluar) > 50 {
		errs["ket_keluar"] = "Keterangan cara keluar maksimal 50 karakter"
	}

	if d.Keadaan == "" {
		errs["keadaan"] = "Keadaan pulang wajib diisi"
	} else if !d.Keadaan.IsValid() {
		errs["keadaan"] = "Keadaan pulang tidak valid (pilihan: Membaik, Sembuh, Rujuk, Keadaan Khusus, Meninggal)"
	}
	if len(d.KetKeadaan) > 50 {
		errs["ket_keadaan"] = "Keterangan keadaan pulang maksimal 50 karakter"
	}

	if d.Dilanjutkan == "" {
		errs["dilanjutkan"] = "Status dilanjutkan wajib diisi"
	} else if !d.Dilanjutkan.IsValid() {
		errs["dilanjutkan"] = "Status dilanjutkan tidak valid (pilihan: Kembali Ke RS, RS Lain, Dokter Luar, Puskesmes, Lainnya)"
	}
	if len(d.KetDilanjutkan) > 50 {
		errs["ket_dilanjutkan"] = "Keterangan status dilanjutkan maksimal 50 karakter"
	}

	if d.Kontrol != "" && d.Kontrol != "0000-00-00 00:00:00" {
		if _, err := time.ParseInLocation("2006-01-02 15:04:05", d.Kontrol, time.Local); err != nil {
			errs["kontrol"] = "Format tanggal & jam kontrol harus YYYY-MM-DD HH:mm:ss"
		}
	}
}

type ResumePasienRanap struct {
	IdKunjungan string `json:"id_kunjungan"`
	NoRawat     string `json:"-"`
	KodeDokter  string `json:"kode_dokter"`
	NamaDokter  string `json:"nama_dokter"`
	InfoKamarInap
	DataResumePasienRanap
}

type SimpanResumePasienRanapRequest struct {
	NoRawat string `json:"no_rawat"`
	DataResumePasienRanap
}

func (req *SimpanResumePasienRanapRequest) Sanitize() {
	req.NoRawat = strings.TrimSpace(req.NoRawat)
	req.DataResumePasienRanap.Sanitize()
}

func (req *SimpanResumePasienRanapRequest) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)
	if req.NoRawat == "" {
		errs["no_rawat"] = "Nomor rawat wajib diisi"
	}
	req.DataResumePasienRanap.Validate(errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateResumePasienRanapRequest struct {
	DataResumePasienRanap
}

func (req *UpdateResumePasienRanapRequest) Sanitize() {
	req.DataResumePasienRanap.Sanitize()
}

func (req *UpdateResumePasienRanapRequest) Validate() apperror.ValidationError {
	errs := make(apperror.ValidationError)
	req.DataResumePasienRanap.Validate(errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}
