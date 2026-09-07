package resumepasien_test

import (
	"testing"

	"erm-dokter/internal/resumepasien"
)

func TestKondisiPulang_IsValid(t *testing.T) {
	tests := []struct {
		val   resumepasien.KondisiPulang
		valid bool
	}{
		{resumepasien.KondisiPulangHidup, true},
		{resumepasien.KondisiPulangMeninggal, true},
		{"Lainnya", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := tt.val.IsValid(); got != tt.valid {
			t.Errorf("KondisiPulang(%s).IsValid() = %v, want %v", tt.val, got, tt.valid)
		}
	}
}

func TestSimpanResumePasienRequest_Sanitize_Defaults(t *testing.T) {
	req := resumepasien.SimpanResumePasienRequest{
		NoRawat: "  2026/09/07/000001  ",
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:  "  Demam tinggi sejak 3 hari  ",
			DiagnosaUtama: "  Demam Tifoid  ",
			// Kode ICD dibiarkan kosong untuk menguji default string kosong ""
			KodeDiagnosaUtama: "",
			KondisiPulang:     "", // dibiarkan kosong untuk menguji default "Hidup"
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
	if req.KondisiPulang != resumepasien.KondisiPulangHidup {
		t.Errorf("expected default KondisiPulang 'Hidup', got '%s'", req.KondisiPulang)
	}
}

func TestSimpanResumePasienRequest_Validation(t *testing.T) {
	// Kasus: Field wajib kosong
	reqEmpty := resumepasien.SimpanResumePasienRequest{}
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
	reqLong := resumepasien.SimpanResumePasienRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:      "Keluhan valid",
			DiagnosaUtama:     "Diagnosa valid",
			KodeDiagnosaUtama: "TOOLONGICDCODE", // > 10 chars
			DiagnosaSekunder3: longStr81,          // > 80 chars
			ProsedurUtama:     longStr81,          // > 80 chars
			KodeProsedurUtama: "TOOLONGPROCEDURE", // > 8 chars
			KondisiPulang:     "InvalidKondisi",
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
	if _, ok := errsLong["kondisi_pulang"]; !ok {
		t.Error("expected error on kondisi_pulang")
	}

	// Kasus: Valid penuh
	reqValid := resumepasien.SimpanResumePasienRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:         "Demam",
			JalannyaPenyakit:     "Demam naik turun",
			PemeriksaanPenunjang: "Widal positif",
			HasilLaborat:         "Leukosit normal",
			DiagnosaUtama:        "Thypoid Fever",
			KodeDiagnosaUtama:    "A01.0",
			DiagnosaSekunder:     "Dispepsia",
			ProsedurUtama:        "Konseling & Terapi Oral",
			KondisiPulang:        resumepasien.KondisiPulangHidup,
			ObatPulang:           "Ciprofloxacin 2x500mg, Paracetamol 3x500mg",
		},
	}
	reqValid.Sanitize()
	if errsValid := reqValid.Validate(); errsValid != nil {
		t.Errorf("expected no validation errors, got: %v", errsValid)
	}
}

func TestUpdateResumePasienRequest_Validation(t *testing.T) {
	reqValid := resumepasien.UpdateResumePasienRequest{
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:  "Batuk pilek",
			DiagnosaUtama: "ISPA",
			KondisiPulang: resumepasien.KondisiPulangHidup,
		},
	}
	reqValid.Sanitize()
	if errs := reqValid.Validate(); errs != nil {
		t.Errorf("expected valid update request, got: %v", errs)
	}

	reqInvalid := resumepasien.UpdateResumePasienRequest{
		DataResumePasienRalan: resumepasien.DataResumePasienRalan{
			KeluhanUtama:  "",
			DiagnosaUtama: "",
			KondisiPulang: "Asal",
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
	if _, ok := errs["kondisi_pulang"]; !ok {
		t.Error("expected error on kondisi_pulang")
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
			Alasan:       "  Panas tinggi 5 hari  ",
			KeluhanUtama: "  Demam menggigil  ",
			Kontrol:      "  2026-09-15 10:00  ", // YYYY-MM-DD HH:mm -> harus auto append :00
		},
	}
	req.Sanitize()

	if req.NoRawat != "2026/09/07/000001" {
		t.Errorf("expected NoRawat '2026/09/07/000001', got '%s'", req.NoRawat)
	}
	if req.DiagnosaAwal != "Demam Thypoid" {
		t.Errorf("expected DiagnosaAwal 'Demam Thypoid', got '%s'", req.DiagnosaAwal)
	}
	if req.Kontrol != "2026-09-15 10:00:00" {
		t.Errorf("expected Kontrol '2026-09-15 10:00:00', got '%s'", req.Kontrol)
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
	if _, ok := errs["alasan"]; !ok {
		t.Error("expected error on alasan")
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
	if _, ok := errs["keadaan"]; !ok {
		t.Error("expected error on keadaan")
	}
	if _, ok := errs["dilanjutkan"]; !ok {
		t.Error("expected error on dilanjutkan")
	}

	// Valid
	reqValid := resumepasien.SimpanResumePasienRanapRequest{
		NoRawat: "2026/09/07/000001",
		DataResumePasienRanap: resumepasien.DataResumePasienRanap{
			DiagnosaAwal:  "Febris H-3",
			Alasan:        "Demam tinggi dan dehidrasi",
			KeluhanUtama:  "Demam tinggi dan lemas",
			DiagnosaUtama: "DHF Grade 1",
			CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
			Keadaan:       resumepasien.KeadaanPulangMembaik,
			Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
			Kontrol:       "2026-09-15 10:00:00",
		},
	}
	reqValid.Sanitize()
	if errsValid := reqValid.Validate(); errsValid != nil {
		t.Errorf("expected no validation errors, got: %v", errsValid)
	}
}

func TestUpdateResumePasienRanapRequest_Validation(t *testing.T) {
	reqValid := resumepasien.UpdateResumePasienRanapRequest{
		DataResumePasienRanap: resumepasien.DataResumePasienRanap{
			DiagnosaAwal:  "Febris",
			Alasan:        "Demam",
			KeluhanUtama:  "Keluhan",
			DiagnosaUtama: "DHF",
			CaraKeluar:    resumepasien.CaraKeluarAtasIzinDokter,
			Keadaan:       resumepasien.KeadaanPulangMembaik,
			Dilanjutkan:   resumepasien.DilanjutkanKembaliKeRS,
		},
	}
	reqValid.Sanitize()
	if errs := reqValid.Validate(); errs != nil {
		t.Errorf("expected valid update request, got: %v", errs)
	}
}
