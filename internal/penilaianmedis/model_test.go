package penilaianmedis

import (
	"testing"
)

func TestEnums_IsValid(t *testing.T) {
	// Anamnesis
	if !Autoanamnesis.IsValid() || !Alloanamnesis.IsValid() {
		t.Error("expected valid Anamnesis enums to return true")
	}
	if Anamnesis("Invalid").IsValid() {
		t.Error("expected invalid Anamnesis to return false")
	}

	// Keadaan
	if !KeadaanSehat.IsValid() || !KeadaanSakitRingan.IsValid() || !KeadaanSakitSedang.IsValid() || !KeadaanSakitBerat.IsValid() {
		t.Error("expected valid Keadaan enums to return true")
	}
	if Keadaan("Kritis").IsValid() {
		t.Error("expected invalid Keadaan to return false")
	}

	// KesadaranAwal
	if !KesadaranComposMentis.IsValid() || !KesadaranApatis.IsValid() || !KesadaranSomnolen.IsValid() || !KesadaranSopor.IsValid() || !KesadaranKoma.IsValid() {
		t.Error("expected valid Kesadaran enums to return true")
	}
	if KesadaranAwal("Delirium").IsValid() {
		t.Error("expected invalid Kesadaran to return false")
	}

	// StatusFisik
	if !StatusFisikNormal.IsValid() || !StatusFisikAbnormal.IsValid() || !StatusFisikTidakDiperiksa.IsValid() {
		t.Error("expected valid StatusFisik enums to return true")
	}
	if StatusFisik("Baik").IsValid() {
		t.Error("expected invalid StatusFisik to return false")
	}
}

func TestOpsiReferensi(t *testing.T) {
	if len(DaftarOpsiAnamnesis()) != 2 {
		t.Errorf("expected 2 anamnesis options, got %d", len(DaftarOpsiAnamnesis()))
	}
	if len(DaftarOpsiKeadaan()) != 4 {
		t.Errorf("expected 4 keadaan options, got %d", len(DaftarOpsiKeadaan()))
	}
	if len(DaftarOpsiKesadaran()) != 5 {
		t.Errorf("expected 5 kesadaran options, got %d", len(DaftarOpsiKesadaran()))
	}
	if len(DaftarOpsiStatusFisik()) != 3 {
		t.Errorf("expected 3 status fisik options, got %d", len(DaftarOpsiStatusFisik()))
	}

	ref := GetReferensiPenilaianMedis()
	if len(ref.Anamnesis) != 2 || len(ref.Keadaan) != 4 || len(ref.Kesadaran) != 5 || len(ref.StatusFisik) != 3 {
		t.Errorf("unexpected global referensi struct counts: %+v", ref)
	}
}

func TestSimpanPenilaianMedisRalanRequest_Sanitize_Defaults(t *testing.T) {
	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: " 2026/04/22/000001 ",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			KeluhanUtama: " Demam tinggi ",
			SuhuTubuh:    " 37,5 ",
			BeratBadan:   " 65,5 ",
			TinggiBadan:  " 170,0 ",
			Diagnosis:    " Febris ",
			TataLaksana:  " Paracetamol ",
		},
	}

	req.Sanitize()

	if req.NoRawat != "2026/04/22/000001" {
		t.Errorf("expected trimmed no_rawat, got %s", req.NoRawat)
	}
	if req.TanggalPenilaian == "" {
		t.Error("expected default TanggalPenilaian to be populated")
	}
	if req.KeluhanUtama != "Demam tinggi" {
		t.Errorf("expected trimmed keluhan utama, got %s", req.KeluhanUtama)
	}
	if req.SuhuTubuh != "37.5" || req.BeratBadan != "65.5" || req.TinggiBadan != "170.0" {
		t.Errorf("expected comma replaced with dot in TTV: suhu=%s, bb=%s, tb=%s", req.SuhuTubuh, req.BeratBadan, req.TinggiBadan)
	}

	// Cek auto defaults
	if req.Anamnesis != Autoanamnesis {
		t.Errorf("expected default Anamnesis to be Autoanamnesis, got %s", req.Anamnesis)
	}
	if req.Keadaan != KeadaanSehat {
		t.Errorf("expected default Keadaan to be Sehat, got %s", req.Keadaan)
	}
	if req.Kesadaran != KesadaranComposMentis {
		t.Errorf("expected default Kesadaran to be Compos Mentis, got %s", req.Kesadaran)
	}
	if req.Kepala != StatusFisikNormal || req.Gigi != StatusFisikNormal || req.Thoraks != StatusFisikNormal {
		t.Errorf("expected default Pemeriksaan Fisik to be Normal, got kepala=%s", req.Kepala)
	}
}

func TestSimpanPenilaianMedisRalanRequest_Validation_RequiredFields(t *testing.T) {
	// 1. Missing no_rawat, keluhan_utama, diagnosis, tata_laksana
	req := SimpanPenilaianMedisRalanRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors, got nil")
	}
	if _, exists := errs["no_rawat"]; !exists {
		t.Error("expected error on no_rawat")
	}
	if _, exists := errs["keluhan_utama"]; !exists {
		t.Error("expected error on keluhan_utama")
	}
	if _, exists := errs["diagnosis"]; !exists {
		t.Error("expected error on diagnosis")
	}
	if _, exists := errs["tata_laksana"]; !exists {
		t.Error("expected error on tata_laksana")
	}

	// 2. Alloanamnesis missing hubungan
	req = SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			Anamnesis:    Alloanamnesis,
			Hubungan:     "", // kosong
			KeluhanUtama: "Nyeri dada",
			Diagnosis:    "Angina",
			TataLaksana:  "ISDN 5mg",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["hubungan"] == "" {
		t.Errorf("expected error on hubungan for Alloanamnesis, got: %+v", errs)
	}

	// 3. Alloanamnesis with valid hubungan
	req.Hubungan = "Anak Kandung"
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid request to have no errors, got: %+v", errs)
	}
}

func TestSimpanPenilaianMedisRalanRequest_Validation_TTV(t *testing.T) {
	// 1. Invalid Tensi format (missing slash)
	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			KeluhanUtama: "Demam",
			Diagnosis:    "Febris",
			TataLaksana:  "Paracetamol",
			Tensi:        "12080",
		},
	}
	req.Sanitize()
	errs := req.Validate()
	if errs == nil || errs["tensi"] == "" {
		t.Errorf("expected error on invalid tensi format, got: %+v", errs)
	}

	// 2. Sistolik <= Diastolik
	req.Tensi = "80/120"
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["tensi"] == "" {
		t.Errorf("expected error when sistolik <= diastolik, got: %+v", errs)
	}

	// 3. Suhu Tubuh out of range
	req.Tensi = "120/80"
	req.SuhuTubuh = "10.0"
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["suhu_tubuh"] == "" {
		t.Errorf("expected error on out of range suhu_tubuh, got: %+v", errs)
	}

	// 4. Valid TTV
	req.SuhuTubuh = "36.8"
	req.Nadi = "80"
	req.Respirasi = "20"
	req.SpO2 = "99"
	req.BeratBadan = "60"
	req.TinggiBadan = "165"
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid TTV to have no errors, got: %+v", errs)
	}
}

func TestSimpanPenilaianMedisRalanRequest_Validation_TanggalPenilaian(t *testing.T) {
	// 1. Invalid date format
	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: "31-08-2026", // format salah
			KeluhanUtama:     "Demam",
			Diagnosis:        "Febris",
			TataLaksana:      "Paracetamol",
		},
	}
	req.Sanitize()
	errs := req.Validate()
	if errs == nil || errs["tanggal_penilaian"] == "" {
		t.Errorf("expected error on invalid tanggal_penilaian format, got: %+v", errs)
	}

	// 2. Future date
	future := "2099-01-01 10:00:00"
	req.TanggalPenilaian = future
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["tanggal_penilaian"] == "" {
		t.Errorf("expected error on future tanggal_penilaian, got: %+v", errs)
	}

	// 3. Valid past date
	req.TanggalPenilaian = "2026-04-22 09:30:00"
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid tanggal_penilaian to have no errors, got: %+v", errs)
	}
}

func TestUpdatePenilaianMedisRalanRequest_Validation(t *testing.T) {
	req := UpdatePenilaianMedisRalanRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty update request")
	}
	if _, exists := errs["keluhan_utama"]; !exists {
		t.Error("expected error on keluhan_utama")
	}
	if _, exists := errs["no_rawat"]; exists {
		t.Error("did not expect error on no_rawat for UpdatePenilaianMedisRalanRequest")
	}

	req = UpdatePenilaianMedisRalanRequest{
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Batuk pilek",
			Diagnosis:        "ISPA",
			TataLaksana:      "Amoxicillin 3x500mg",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid update request to have no errors, got: %+v", errs)
	}
}

// ==========================================
// IGD MODEL & VALIDATION TESTS
// ==========================================

func TestSimpanPenilaianMedisIGDRequest_Sanitize_Defaults(t *testing.T) {
	req := SimpanPenilaianMedisIGDRequest{
		NoRawat: " 2026/04/22/000002 ",
		DataPenilaianMedisIGD: DataPenilaianMedisIGD{
			KeluhanUtama: " Nyeri dada tembus ke belakang ",
			SuhuTubuh:    " 36,8 ",
			BeratBadan:   " 70,0 ",
			TinggiBadan:  " 172,0 ",
			Diagnosis:    " STEMI Anterior ",
			TataLaksana:  " Loading Aspilet & Clopidogrel ",
			EKG:          " ST Elevasi di V1-V4 ",
			Radiologi:    " Cardiomegaly ",
			Laboratorium: " Troponin I Positif ",
		},
	}

	req.Sanitize()

	if req.NoRawat != "2026/04/22/000002" {
		t.Errorf("expected trimmed no_rawat, got %s", req.NoRawat)
	}
	if req.TanggalPenilaian == "" {
		t.Error("expected default TanggalPenilaian to be populated")
	}
	if req.KeluhanUtama != "Nyeri dada tembus ke belakang" {
		t.Errorf("expected trimmed keluhan utama, got %s", req.KeluhanUtama)
	}
	if req.EKG != "ST Elevasi di V1-V4" || req.Radiologi != "Cardiomegaly" || req.Laboratorium != "Troponin I Positif" {
		t.Errorf("expected trimmed EKG/Rad/Lab, got ekg=%s, rad=%s, lab=%s", req.EKG, req.Radiologi, req.Laboratorium)
	}

	// Cek default 8 organ fisik IGD (termasuk Mata dan Leher)
	if req.Kepala != StatusFisikNormal || req.Mata != StatusFisikNormal || req.Gigi != StatusFisikNormal ||
		req.Leher != StatusFisikNormal || req.Thoraks != StatusFisikNormal || req.Abdomen != StatusFisikNormal ||
		req.Genital != StatusFisikNormal || req.Ekstremitas != StatusFisikNormal {
		t.Errorf("expected default Pemeriksaan Fisik IGD to be Normal, got mata=%s, leher=%s", req.Mata, req.Leher)
	}
}

func TestSimpanPenilaianMedisIGDRequest_Validation(t *testing.T) {
	// 1. Missing required fields
	req := SimpanPenilaianMedisIGDRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors on empty IGD request, got nil")
	}
	if _, exists := errs["no_rawat"]; !exists {
		t.Error("expected error on no_rawat")
	}
	if _, exists := errs["keluhan_utama"]; !exists {
		t.Error("expected error on keluhan_utama")
	}
	if _, exists := errs["diagnosis"]; !exists {
		t.Error("expected error on diagnosis")
	}
	if _, exists := errs["tata_laksana"]; !exists {
		t.Error("expected error on tata_laksana")
	}

	// 2. Valid IGD Request
	req = SimpanPenilaianMedisIGDRequest{
		NoRawat: "2026/04/22/000002",
		DataPenilaianMedisIGD: DataPenilaianMedisIGD{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Sesak napas akut",
			Diagnosis:        "Asma Akut Berat",
			TataLaksana:      "Nebulisasi Ventolin + Pulmicort",
			Tensi:            "130/80",
			SuhuTubuh:        "36.7",
			Nadi:             "105",
			Respirasi:        "28",
			SpO2:             "94",
			EKG:              "Sinus Takikardia",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid IGD request to have no errors, got: %+v", errs)
	}
}

func TestUpdatePenilaianMedisIGDRequest_Validation(t *testing.T) {
	req := UpdatePenilaianMedisIGDRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty update IGD request")
	}

	req = UpdatePenilaianMedisIGDRequest{
		DataPenilaianMedisIGD: DataPenilaianMedisIGD{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Lemas & pusing",
			Diagnosis:        "Hipoglikemia",
			TataLaksana:      "Bolus Dextrose 40%",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid update IGD request to have no errors, got: %+v", errs)
	}
}
