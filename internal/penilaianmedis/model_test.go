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
	// 1. Empty request
	req := SimpanPenilaianMedisRalanRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty request, got nil")
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

	// 2. Alloanamnesis without hubungan
	req = SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			Anamnesis:    Alloanamnesis,
			Hubungan:     "",
			KeluhanUtama: "Nyeri dada",
			Diagnosis:    "Angina",
			TataLaksana:  "ISDN 5mg",
		},
	}
	errs = req.Validate()
	if errs == nil || errs["hubungan"] == "" {
		t.Errorf("expected error on hubungan for Alloanamnesis, got: %+v", errs)
	}

	// 3. Alloanamnesis with hubungan (Valid)
	req.Hubungan = "Anak Kandung"
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid request to have no errors, got: %+v", errs)
	}
}

func TestSimpanPenilaianMedisRalanRequest_Validation_InvalidEnums(t *testing.T) {
	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			KeluhanUtama: "Pusing",
			Diagnosis:    "Vertigo",
			TataLaksana:  "Betahistine",
			Keadaan:      Keadaan("InvalidKeadaan"),
			Kesadaran:    KesadaranAwal("InvalidKesadaran"),
			Kepala:       StatusFisik("InvalidKepala"),
		},
	}

	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for invalid enums")
	}
	if _, exists := errs["keadaan"]; !exists {
		t.Error("expected error for invalid keadaan")
	}
	if _, exists := errs["kesadaran"]; !exists {
		t.Error("expected error for invalid kesadaran")
	}
	if _, exists := errs["kepala"]; !exists {
		t.Error("expected error for invalid kepala")
	}
}

func TestSimpanPenilaianMedisRalanRequest_Validation_TTV(t *testing.T) {
	// 1. Tensi validation
	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			KeluhanUtama: "Sakit kepala",
			Diagnosis:    "Hipertensi",
			TataLaksana:  "Amlodipine 5mg",
			Tensi:        "120", // format salah
		},
	}
	errs := req.Validate()
	if errs == nil || errs["tensi"] == "" {
		t.Errorf("expected error on invalid tensi format, got: %+v", errs)
	}

	req.Tensi = "80/120" // sistolik <= diastolik
	errs = req.Validate()
	if errs == nil || errs["tensi"] == "" {
		t.Errorf("expected error when sistolik <= diastolik, got: %+v", errs)
	}

	req.Tensi = "120/80" // valid tensi

	// 2. Suhu validation
	req.SuhuTubuh = "20.0" // out of range (< 25)
	errs = req.Validate()
	if errs == nil || errs["suhu_tubuh"] == "" {
		t.Errorf("expected error on suhu out of range, got: %+v", errs)
	}
	req.SuhuTubuh = "36.5" // valid suhu

	// 3. Nadi validation
	req.Nadi = "10" // out of range (< 20)
	errs = req.Validate()
	if errs == nil || errs["nadi"] == "" {
		t.Errorf("expected error on nadi out of range, got: %+v", errs)
	}
	req.Nadi = "80" // valid nadi

	// 4. Respirasi validation
	req.Respirasi = "150" // out of range (> 100)
	errs = req.Validate()
	if errs == nil || errs["respirasi"] == "" {
		t.Errorf("expected error on respirasi out of range, got: %+v", errs)
	}
	req.Respirasi = "20" // valid respirasi

	// 5. Tinggi & Berat Badan
	req.TinggiBadan = "300" // out of range (> 250)
	errs = req.Validate()
	if errs == nil || errs["tinggi_badan"] == "" {
		t.Errorf("expected error on tinggi badan out of range, got: %+v", errs)
	}
	req.TinggiBadan = "170"

	req.BeratBadan = "600" // out of range (> 500)
	errs = req.Validate()
	if errs == nil || errs["berat_badan"] == "" {
		t.Errorf("expected error on berat badan out of range, got: %+v", errs)
	}
	req.BeratBadan = "65"

	// 6. SpO2
	req.SpO2 = "120" // out of range (> 100)
	errs = req.Validate()
	if errs == nil || errs["spo2"] == "" {
		t.Errorf("expected error on SpO2 out of range, got: %+v", errs)
	}
	req.SpO2 = "98"

	// Valid all
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
	errs := req.Validate()
	if errs == nil || errs["tanggal_penilaian"] == "" {
		t.Errorf("expected error on invalid tanggal_penilaian format, got: %+v", errs)
	}

	// 2. Future date
	future := "2099-01-01 10:00:00"
	req.TanggalPenilaian = future
	errs = req.Validate()
	if errs == nil || errs["tanggal_penilaian"] == "" {
		t.Errorf("expected error on future tanggal_penilaian, got: %+v", errs)
	}

	// 3. Valid past date
	req.TanggalPenilaian = "2026-04-22 09:30:00"
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid tanggal_penilaian to have no errors, got: %+v", errs)
	}
}

func TestUpdatePenilaianMedisRalanRequest_Validation(t *testing.T) {
	// 1. Empty update request (missing required clinical fields)
	req := UpdatePenilaianMedisRalanRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty update request")
	}
	if _, exists := errs["keluhan_utama"]; !exists {
		t.Error("expected error on keluhan_utama")
	}
	// Note: no_rawat is NOT required in UpdatePenilaianMedisRalanRequest
	if _, exists := errs["no_rawat"]; exists {
		t.Error("did not expect error on no_rawat for UpdatePenilaianMedisRalanRequest")
	}

	// 2. Valid update request
	req = UpdatePenilaianMedisRalanRequest{
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Batuk pilek",
			Diagnosis:        "ISPA",
			TataLaksana:      "Amoxicillin 3x500mg",
		},
	}
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid update request to have no errors, got: %+v", errs)
	}
}

