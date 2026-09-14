package resumepasien_test

import (
	"testing"

	"erm-dokter/internal/resumepasien"
)

func TestKeadaanPulang_IsValid(t *testing.T) {
	tests := []struct {
		val        resumepasien.KeadaanPulang
		validRalan bool
		validRanap bool
	}{
		{resumepasien.KeadaanPulangHidup, true, false},
		{resumepasien.KeadaanPulangMembaik, false, true},
		{resumepasien.KeadaanPulangSembuh, false, true},
		{resumepasien.KeadaanPulangRujuk, false, true},
		{resumepasien.KeadaanPulangKeadaanKhusus, false, true},
		{resumepasien.KeadaanPulangMeninggal, true, true},
		{"Lainnya", false, false},
		{"", false, false},
	}

	for _, tt := range tests {
		if got := tt.val.IsValidRalan(); got != tt.validRalan {
			t.Errorf("KeadaanPulang(%s).IsValidRalan() = %v, want %v", tt.val, got, tt.validRalan)
		}
		if got := tt.val.IsValidRanap(); got != tt.validRanap {
			t.Errorf("KeadaanPulang(%s).IsValidRanap() = %v, want %v", tt.val, got, tt.validRanap)
		}
	}
}

func TestSimpanResumePasienRalanRequest_Sanitize_Defaults(t *testing.T) {
	req := resumepasien.SimpanResumePasienRalanRequest{
		NoRawat: "  2026/09/07/000001  ",
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:  "  Demam tinggi sejak 3 hari  ",
			DiagnosaUtama: "  Demam Tifoid  ",
			// Kode ICD dibiarkan kosong untuk menguji default string kosong ""
			KodeDiagnosaUtama: "",
			KeadaanPulang:     "", // dibiarkan kosong untuk menguji default "Hidup"
		},
	}

	req.Sanitize()

	if req.NoRawat != "2026/09/07/000001" {
		t.Errorf("expected NoRawat '2026/09/07/000001', got '%s'", req.NoRawat)
	}
	if req.KeluhanUtama != "Demam tinggi sejak 3 hari" {
		t.Errorf("expected KeluhanUtama 'Demam tinggi sejak 3 hari', got '%s'", req.KeluhanUtama)
	}
	if req.DiagnosaUtama != "Demam Tifoid" {
		t.Errorf("expected DiagnosaUtama 'Demam Tifoid', got '%s'", req.DiagnosaUtama)
	}
	if req.KodeDiagnosaUtama != "" {
		t.Errorf("expected KodeDiagnosaUtama '', got '%s'", req.KodeDiagnosaUtama)
	}
	if req.KeadaanPulang != resumepasien.KeadaanPulangHidup {
		t.Errorf("expected default KeadaanPulang 'Hidup', got '%s'", req.KeadaanPulang)
	}
}

func TestSimpanResumePasienRalanRequest_Validation(t *testing.T) {
	// Kasus: Field wajib kosong
	reqEmpty := resumepasien.SimpanResumePasienRalanRequest{}
	reqEmpty.Sanitize()
	errs := reqEmpty.Validate()
	if errs == nil {
		t.Fatal("expected validation errors, got nil")
	}
	if _, ok := errs["no_rawat"]; !ok {
		t.Error("expected error on no_rawat")
	}
	if _, ok := errs["keluhan_utama"]; !ok {
		t.Error("expected error on keluhan_utama")
	}
	if _, ok := errs["diagnosa_utama"]; !ok {
		t.Error("expected error on diagnosa_utama")
	}

	// Kasus: Batas panjang karakter terlampaui
	longStr81 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	reqLong := resumepasien.SimpanResumePasienRalanRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:      "Keluhan valid",
			DiagnosaUtama:     "Diagnosa valid",
			KodeDiagnosaUtama: "TOOLONGICDCODE", // > 10 chars
			DiagnosaSekunder3: longStr81,          // > 80 chars
			ProsedurUtama:     longStr81,          // > 80 chars
			KodeProsedurUtama: "TOOLONGPROCEDURE", // > 8 chars
			KeadaanPulang:     "InvalidKeadaan",
		},
	}
	errsLong := reqLong.Validate()
	if errsLong == nil {
		t.Fatal("expected validation errors for long strings, got nil")
	}
	if _, ok := errsLong["kode_diagnosa_utama"]; !ok {
		t.Error("expected error on kode_diagnosa_utama")
	}
	if _, ok := errsLong["diagnosa_sekunder3"]; !ok {
		t.Error("expected error on diagnosa_sekunder3")
	}
	if _, ok := errsLong["prosedur_utama"]; !ok {
		t.Error("expected error on prosedur_utama")
	}
	if _, ok := errsLong["kode_prosedur_utama"]; !ok {
		t.Error("expected error on kode_prosedur_utama")
	}
	if _, ok := errsLong["keadaan_pulang"]; !ok {
		t.Error("expected error on keadaan_pulang")
	}

	// Kasus: Valid penuh
	reqValid := resumepasien.SimpanResumePasienRalanRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:                 "Demam",
			JalannyaPenyakit:             "Demam naik turun",
			HasilPemeriksaanRadiologi:    "Widal positif",
			HasilPemeriksaanLaboratorium: "Leukosit normal",
			DiagnosaUtama:                "Thypoid Fever",
			KodeDiagnosaUtama:            "A01.0",
			DiagnosaSekunder:             "Dispepsia",
			ProsedurUtama:                "Konseling & Terapi Oral",
			KeadaanPulang:                resumepasien.KeadaanPulangHidup,
			ObatAtauInstruksi:            "Ciprofloxacin 2x500mg, Paracetamol 3x500mg",
		},
	}
	reqValid.Sanitize()
	if errsValid := reqValid.Validate(); errsValid != nil {
		t.Errorf("expected no validation errors, got: %v", errsValid)
	}
}

func TestUpdateResumePasienRalanRequest_Validation(t *testing.T) {
	reqValid := resumepasien.UpdateResumePasienRalanRequest{
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:  "Batuk pilek",
			DiagnosaUtama: "ISPA",
			KeadaanPulang: resumepasien.KeadaanPulangHidup,
		},
	}
	reqValid.Sanitize()
	if errs := reqValid.Validate(); errs != nil {
		t.Errorf("expected valid update request, got: %v", errs)
	}

	reqInvalid := resumepasien.UpdateResumePasienRalanRequest{
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:  "",
			DiagnosaUtama: "",
			KeadaanPulang: "Asal",
		},
	}
	reqInvalid.Sanitize()
	errs := reqInvalid.Validate()
	if errs == nil {
		t.Fatal("expected validation errors, got nil")
	}
	if _, ok := errs["keluhan_utama"]; !ok {
		t.Error("expected error on keluhan_utama")
	}
	if _, ok := errs["diagnosa_utama"]; !ok {
		t.Error("expected error on diagnosa_utama")
	}
	if _, ok := errs["keadaan_pulang"]; !ok {
		t.Error("expected error on keadaan_pulang")
	}
}

func TestRanapEnums_IsValid(t *testing.T) {
	// Cara Keluar
	if !resumepasien.CaraKeluarAtasIzinDokter.IsValid() {
		t.Error("expected CaraKeluarAtasIzinDokter to be valid")
	}
	if resumepasien.CaraKeluar("Sembarangan").IsValid() {
		t.Error("expected Sembarangan to be invalid")
	}

	// Keadaan Pulang
	if !resumepasien.KeadaanPulangMembaik.IsValid() {
		t.Error("expected KeadaanPulangMembaik to be valid")
	}
	if resumepasien.KeadaanPulang("Sembarangan").IsValid() {
		t.Error("expected Sembarangan to be invalid")
	}

	// Dilanjutkan
	if !resumepasien.DilanjutkanKembaliKeRS.IsValid() {
		t.Error("expected DilanjutkanKembaliKeRS to be valid")
	}
	if resumepasien.Dilanjutkan("Sembarangan").IsValid() {
		t.Error("expected Sembarangan to be invalid")
	}
}

func TestSimpanResumePasienRanapRequest_Sanitize_Defaults(t *testing.T) {
	req := resumepasien.SimpanResumePasienRanapRequest{
		NoRawat: "  2026/09/07/000001  ",
		DataResumePasienRanap: resumepasien.DataResumePasienRanap{
			DiagnosaAwal: "  Demam Thypoid  ",
			AlasanRawat:  "  Panas tinggi 5 hari  ",
			KeluhanUtama: "  Demam menggigil  ",
			WaktuKontrol: "  2026-09-15 10:00  ", // YYYY-MM-DD HH:mm -> harus auto append :00
		},
	}
	req.Sanitize()

	if req.NoRawat != "2026/09/07/000001" {
		t.Errorf("expected NoRawat '2026/09/07/000001', got '%s'", req.NoRawat)
	}
	if req.DiagnosaAwal != "Demam Thypoid" {
		t.Errorf("expected DiagnosaAwal 'Demam Thypoid', got '%s'", req.DiagnosaAwal)
	}
	if req.AlasanRawat != "Panas tinggi 5 hari" {
		t.Errorf("expected AlasanRawat 'Panas tinggi 5 hari', got '%s'", req.AlasanRawat)
	}
	if req.WaktuKontrol != "2026-09-15 10:00:00" {
		t.Errorf("expected WaktuKontrol '2026-09-15 10:00:00', got '%s'", req.WaktuKontrol)
	}
}

func TestSimpanResumePasienRanapRequest_Sanitize_WaktuKontrol_Meninggal_Dan_Rujuk(t *testing.T) {
	// Kasus 1: KeadaanPulang Meninggal -> waktu_kontrol harus otomatis di-set ke 0000-00-00 00:00:00
	reqMeninggal := resumepasien.SimpanResumePasienRanapRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRanap: resumepasien.DataResumePasienRanap{
			KeadaanPulang: resumepasien.KeadaanPulangMeninggal,
			WaktuKontrol:  "2026-09-15 10:00:00", // meskipun diisi, harus ter-override ke default
		},
	}
	reqMeninggal.Sanitize()
	if reqMeninggal.WaktuKontrol != "0000-00-00 00:00:00" {
		t.Errorf("expected WaktuKontrol '0000-00-00 00:00:00' for Meninggal, got '%s'", reqMeninggal.WaktuKontrol)
	}

	// Kasus 2: KeadaanPulang Rujuk -> waktu_kontrol harus otomatis di-set ke 0000-00-00 00:00:00
	reqRujuk := resumepasien.SimpanResumePasienRanapRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRanap: resumepasien.DataResumePasienRanap{
			KeadaanPulang: resumepasien.KeadaanPulangRujuk,
			WaktuKontrol:  "",
		},
	}
	reqRujuk.Sanitize()
	if reqRujuk.WaktuKontrol != "0000-00-00 00:00:00" {
		t.Errorf("expected WaktuKontrol '0000-00-00 00:00:00' for Rujuk, got '%s'", reqRujuk.WaktuKontrol)
	}
}

func TestSimpanResumePasienRanapRequest_Validation(t *testing.T) {
	// Kosong
	reqEmpty := resumepasien.SimpanResumePasienRanapRequest{}
	reqEmpty.Sanitize()
	errs := reqEmpty.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty ranap request")
	}
	if _, ok := errs["no_rawat"]; !ok {
		t.Error("expected error on no_rawat")
	}
	if _, ok := errs["diagnosa_awal"]; !ok {
		t.Error("expected error on diagnosa_awal")
	}
	if _, ok := errs["alasan_rawat"]; !ok {
		t.Error("expected error on alasan_rawat")
	}
	if _, ok := errs["keluhan_utama"]; !ok {
		t.Error("expected error on keluhan_utama")
	}
	if _, ok := errs["diagnosa_utama"]; !ok {
		t.Error("expected error on diagnosa_utama")
	}
	if _, ok := errs["cara_keluar"]; !ok {
		t.Error("expected error on cara_keluar")
	}
	if _, ok := errs["keadaan_pulang"]; !ok {
		t.Error("expected error on keadaan_pulang")
	}
	if _, ok := errs["dilanjutkan"]; !ok {
		t.Error("expected error on dilanjutkan")
	}

	// Valid
	reqValid := resumepasien.SimpanResumePasienRanapRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRanap: resumepasien.DataResumePasienRanap{
			DiagnosaAwal:  "Febris H-3",
			AlasanRawat:   "Demam tinggi dan dehidrasi",
			KeluhanUtama:  "Demam tinggi dan lemas",
			DiagnosaUtama: "DHF Grade 1",
			CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
			KeadaanPulang: resumepasien.KeadaanPulangMembaik,
			Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
			WaktuKontrol:  "2026-09-15 10:00:00",
		},
	}
	reqValid.Sanitize()
	if errsValid := reqValid.Validate(); errsValid != nil {
		t.Errorf("expected no validation errors, got: %v", errsValid)
	}

	// Invalid format: "-" tidak boleh dan harus ditolak
	reqInvalidDash := reqValid
	reqInvalidDash.WaktuKontrol = "-"
	reqInvalidDash.Sanitize()
	errsDash := reqInvalidDash.Validate()
	if errsDash == nil || errsDash["waktu_kontrol"] == "" {
		t.Errorf("expected validation error for WaktuKontrol '-', got: %v", errsDash)
	}

	// Format 0000-00-00 00:00:00 tidak boleh untuk pasien selain meninggal/rujuk (harus ditolak)
	reqZero := reqValid
	reqZero.WaktuKontrol = "0000-00-00 00:00:00"
	reqZero.Sanitize()
	errsZero := reqZero.Validate()
	if errsZero == nil || errsZero["waktu_kontrol"] == "" {
		t.Errorf("expected validation error for 0000-00-00 00:00:00 on Membaik, got: %v", errsZero)
	}

	// String kosong tidak boleh untuk pasien selain meninggal/rujuk (harus ditolak)
	reqEmptyWaktu := reqValid
	reqEmptyWaktu.WaktuKontrol = ""
	reqEmptyWaktu.Sanitize()
	errsEmpty := reqEmptyWaktu.Validate()
	if errsEmpty == nil || errsEmpty["waktu_kontrol"] == "" {
		t.Errorf("expected validation error for empty WaktuKontrol on Membaik, got: %v", errsEmpty)
	}

	// Jam kontrol >= 14:00 harus ditolak (jam 14:00:00)
	reqHour14 := reqValid
	reqHour14.WaktuKontrol = "2026-09-15 14:00:00"
	reqHour14.Sanitize()
	errsHour14 := reqHour14.Validate()
	if errsHour14 == nil || errsHour14["waktu_kontrol"] != "Jam kontrol harus kurang dari jam 14:00" {
		t.Errorf("expected 'Jam kontrol harus kurang dari jam 14:00' for 14:00:00, got: %v", errsHour14)
	}

	// Jam kontrol >= 14:00 harus ditolak (jam 15:30:00)
	reqHour15 := reqValid
	reqHour15.WaktuKontrol = "2026-09-15 15:30:00"
	reqHour15.Sanitize()
	errsHour15 := reqHour15.Validate()
	if errsHour15 == nil || errsHour15["waktu_kontrol"] != "Jam kontrol harus kurang dari jam 14:00" {
		t.Errorf("expected 'Jam kontrol harus kurang dari jam 14:00' for 15:30:00, got: %v", errsHour15)
	}

	// Jam kontrol sebelum 14:00 valid (misal jam 13:59:00)
	reqHour13 := reqValid
	reqHour13.WaktuKontrol = "2026-09-15 13:59:00"
	reqHour13.Sanitize()
	if errsHour13 := reqHour13.Validate(); errsHour13 != nil {
		t.Errorf("expected no validation error for 13:59:00, got: %v", errsHour13)
	}

	// Pasien Meninggal mengirim string kosong -> disanitize ke 0000-00-00 00:00:00 dan valid
	reqMeninggal := reqValid
	reqMeninggal.KeadaanPulang = resumepasien.KeadaanPulangMeninggal
	reqMeninggal.WaktuKontrol = ""
	reqMeninggal.Sanitize()
	if reqMeninggal.WaktuKontrol != "0000-00-00 00:00:00" {
		t.Errorf("expected WaktuKontrol '0000-00-00 00:00:00' for Meninggal, got '%s'", reqMeninggal.WaktuKontrol)
	}
	if errsMeninggal := reqMeninggal.Validate(); errsMeninggal != nil {
		t.Errorf("expected no validation error for Meninggal with empty WaktuKontrol, got: %v", errsMeninggal)
	}

	// Pasien Rujuk mengirim string kosong -> disanitize ke 0000-00-00 00:00:00 dan valid
	reqRujuk := reqValid
	reqRujuk.KeadaanPulang = resumepasien.KeadaanPulangRujuk
	reqRujuk.WaktuKontrol = ""
	reqRujuk.Sanitize()
	if reqRujuk.WaktuKontrol != "0000-00-00 00:00:00" {
		t.Errorf("expected WaktuKontrol '0000-00-00 00:00:00' for Rujuk, got '%s'", reqRujuk.WaktuKontrol)
	}
	if errsRujuk := reqRujuk.Validate(); errsRujuk != nil {
		t.Errorf("expected no validation error for Rujuk with empty WaktuKontrol, got: %v", errsRujuk)
	}

	// Pasien Rujuk mengirim tanggal kontrol valid -> dibiarkan dan valid
	reqRujukDate := reqValid
	reqRujukDate.KeadaanPulang = resumepasien.KeadaanPulangRujuk
	reqRujukDate.WaktuKontrol = "2026-09-18 10:00:00"
	reqRujukDate.Sanitize()
	if reqRujukDate.WaktuKontrol != "2026-09-18 10:00:00" {
		t.Errorf("expected WaktuKontrol '2026-09-18 10:00:00' for Rujuk with valid date, got '%s'", reqRujukDate.WaktuKontrol)
	}
	if errsRujukDate := reqRujukDate.Validate(); errsRujukDate != nil {
		t.Errorf("expected no validation error for Rujuk with valid date, got: %v", errsRujukDate)
	}
}

func TestUpdateResumePasienRanapRequest_Validation(t *testing.T) {
	reqValid := resumepasien.UpdateResumePasienRanapRequest{
		DataResumePasienRanap: resumepasien.DataResumePasienRanap{
			DiagnosaAwal:  "Febris",
			AlasanRawat:   "Demam",
			KeluhanUtama:  "Keluhan",
			DiagnosaUtama: "DHF",
			CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
			KeadaanPulang: resumepasien.KeadaanPulangMembaik,
			Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
			WaktuKontrol:  "2026-09-15 10:00:00",
		},
	}
	reqValid.Sanitize()
	if errs := reqValid.Validate(); errs != nil {
		t.Errorf("expected valid update request, got: %v", errs)
	}
}
