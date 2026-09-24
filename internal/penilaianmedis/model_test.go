package penilaianmedis

import (
	"testing"
)

func TestEnums_IsValid(t *testing.T) {

	if !Autoanamnesis.IsValid() || !Alloanamnesis.IsValid() {
		t.Error("expected valid Anamnesis enums to return true")
	}
	if Anamnesis("Invalid").IsValid() {
		t.Error("expected invalid Anamnesis to return false")
	}

	if !KeadaanSehat.IsValid() || !KeadaanSakitRingan.IsValid() || !KeadaanSakitSedang.IsValid() || !KeadaanSakitBerat.IsValid() {
		t.Error("expected valid Keadaan enums to return true")
	}
	if Keadaan("Kritis").IsValid() {
		t.Error("expected invalid Keadaan to return false")
	}

	if !KesadaranComposMentis.IsValid() || !KesadaranApatis.IsValid() || !KesadaranSomnolen.IsValid() || !KesadaranSopor.IsValid() || !KesadaranKoma.IsValid() {
		t.Error("expected valid Kesadaran enums to return true")
	}
	if KesadaranAwal("Delirium").IsValid() {
		t.Error("expected invalid Kesadaran to return false")
	}

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
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["hubungan"] == "" {
		t.Errorf("expected error on hubungan for Alloanamnesis, got: %+v", errs)
	}

	req.Hubungan = "Anak Kandung"
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid request to have no errors, got: %+v", errs)
	}
}

func TestSimpanPenilaianMedisRalanRequest_Validation_TTV(t *testing.T) {

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

	req.Tensi = "80/120"
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["tensi"] == "" {
		t.Errorf("expected error when sistolik <= diastolik, got: %+v", errs)
	}

	req.Tensi = "120/80"
	req.SuhuTubuh = "10.0"
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["suhu_tubuh"] == "" {
		t.Errorf("expected error on out of range suhu_tubuh, got: %+v", errs)
	}

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

	req := SimpanPenilaianMedisRalanRequest{
		NoRawat: "2026/04/22/000001",
		DataPenilaianMedisRalan: DataPenilaianMedisRalan{
			TanggalPenilaian: "31-08-2026",
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

	future := "2099-01-01 10:00:00"
	req.TanggalPenilaian = future
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["tanggal_penilaian"] == "" {
		t.Errorf("expected error on future tanggal_penilaian, got: %+v", errs)
	}

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

	if req.Kepala != StatusFisikNormal || req.Mata != StatusFisikNormal || req.Gigi != StatusFisikNormal ||
		req.Leher != StatusFisikNormal || req.Thoraks != StatusFisikNormal || req.Abdomen != StatusFisikNormal ||
		req.Genital != StatusFisikNormal || req.Ekstremitas != StatusFisikNormal {
		t.Errorf("expected default Pemeriksaan Fisik IGD to be Normal, got mata=%s, leher=%s", req.Mata, req.Leher)
	}
}

func TestSimpanPenilaianMedisIGDRequest_Validation(t *testing.T) {

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

func TestSimpanPenilaianMedisRanapRequest_Sanitize_Defaults(t *testing.T) {
	req := SimpanPenilaianMedisRanapRequest{
		NoRawat: " 2026/04/22/000003 ",
		DataPenilaianMedisRanap: DataPenilaianMedisRanap{
			KeluhanUtama: " Sesak nafas berat ",
			SuhuTubuh:    " 38,2 ",
			BeratBadan:   " 70,5 ",
			TinggiBadan:  " 165,0 ",
			Diagnosis:    " Pneumonia Bilateral ",
			TataLaksana:  " Ceftriaxone 1g IV ",
		},
	}

	req.Sanitize()

	if req.NoRawat != "2026/04/22/000003" {
		t.Errorf("expected trimmed no_rawat, got %s", req.NoRawat)
	}
	if req.TanggalPenilaian == "" {
		t.Error("expected default TanggalPenilaian to be populated")
	}
	if req.KeluhanUtama != "Sesak nafas berat" {
		t.Errorf("expected trimmed keluhan utama, got %s", req.KeluhanUtama)
	}
	if req.SuhuTubuh != "38.2" || req.BeratBadan != "70.5" || req.TinggiBadan != "165.0" {
		t.Errorf("expected comma replaced with dot in TTV: suhu=%s, bb=%s, tb=%s", req.SuhuTubuh, req.BeratBadan, req.TinggiBadan)
	}

	if req.Anamnesis != Autoanamnesis {
		t.Errorf("expected default Anamnesis to be Autoanamnesis, got %s", req.Anamnesis)
	}
	if req.Keadaan != KeadaanSehat {
		t.Errorf("expected default Keadaan to be Sehat, got %s", req.Keadaan)
	}
	if req.Kesadaran != KesadaranComposMentis {
		t.Errorf("expected default Kesadaran to be Compos Mentis, got %s", req.Kesadaran)
	}
	if req.Kepala != StatusFisikNormal || req.Mata != StatusFisikNormal || req.Gigi != StatusFisikNormal ||
		req.TelingaHidungTenggorok != StatusFisikNormal || req.Thoraks != StatusFisikNormal ||
		req.Jantung != StatusFisikNormal || req.Paru != StatusFisikNormal ||
		req.Abdomen != StatusFisikNormal || req.Genital != StatusFisikNormal ||
		req.Ekstremitas != StatusFisikNormal || req.Kulit != StatusFisikNormal {
		t.Errorf("expected default Pemeriksaan Fisik to be Normal, got kepala=%s, jantung=%s, paru=%s", req.Kepala, req.Jantung, req.Paru)
	}
}

func TestSimpanPenilaianMedisRanapRequest_Validation(t *testing.T) {

	req := SimpanPenilaianMedisRanapRequest{}
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

	req = SimpanPenilaianMedisRanapRequest{
		NoRawat: "2026/04/22/000003",
		DataPenilaianMedisRanap: DataPenilaianMedisRanap{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Sesak nafas memberat sejak 2 hari",
			Diagnosis:        "Pneumonia Komuniti",
			TataLaksana:      "O2 nasal 3 lpm, IVFD RL 20 tpm, Ceftriaxone 1x2g IV",
			Tensi:            "120/80",
			SuhuTubuh:        "37.8",
			Nadi:             "92",
			Respirasi:        "24",
			SpO2:             "96",
			Laboratorium:     "Leukosit 14.500",
			Radiologi:        "Infiltrat pada lobus kanan bawah",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid Ranap request to have no errors, got: %+v", errs)
	}
}

func TestUpdatePenilaianMedisRanapRequest_Validation(t *testing.T) {
	req := UpdatePenilaianMedisRanapRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty update Ranap request")
	}

	req = UpdatePenilaianMedisRanapRequest{
		DataPenilaianMedisRanap: DataPenilaianMedisRanap{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Sesak nafas berkurang",
			Diagnosis:        "Pneumonia Komuniti Perbaikan",
			TataLaksana:      "Lanjut antibiotik IV",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid update Ranap request to have no errors, got: %+v", errs)
	}
}

func TestKontraksi_IsValid(t *testing.T) {
	if !KontraksiAda.IsValid() || !KontraksiTidak.IsValid() {
		t.Error("expected valid Kontraksi enums to return true")
	}
	if Kontraksi("Jarang").IsValid() {
		t.Error("expected invalid Kontraksi to return false")
	}
}

func TestSimpanPenilaianMedisRalanKandunganRequest_Sanitize_Defaults(t *testing.T) {
	req := SimpanPenilaianMedisRalanKandunganRequest{
		NoRawat: " 2026/04/22/000004 ",
		DataPenilaianMedisRalanKandungan: DataPenilaianMedisRalanKandungan{
			KeluhanUtama:       " Perut kencang-kencang ",
			SuhuTubuh:          " 36,8 ",
			BeratBadan:         " 62,5 ",
			TinggiBadan:        " 158,0 ",
			Diagnosis:          " G1P0A0 hamil 38 minggu ",
			TataLaksana:        " Observasi kemajuan persalinan ",
			TinggiFundusUteri:  " 30 cm ",
			TaksiranBeratJanin: " 3000 gr ",
			His:                " 3x/10m/40s ",
			DenyutJantungJanin: " 140 dpm ",
		},
	}

	req.Sanitize()

	if req.NoRawat != "2026/04/22/000004" {
		t.Errorf("expected trimmed no_rawat, got %s", req.NoRawat)
	}
	if req.TanggalPenilaian == "" {
		t.Error("expected default TanggalPenilaian to be populated")
	}
	if req.KeluhanUtama != "Perut kencang-kencang" {
		t.Errorf("expected trimmed keluhan utama, got %s", req.KeluhanUtama)
	}
	if req.SuhuTubuh != "36.8" || req.BeratBadan != "62.5" || req.TinggiBadan != "158.0" {
		t.Errorf("expected comma replaced with dot in TTV: suhu=%s, bb=%s, tb=%s", req.SuhuTubuh, req.BeratBadan, req.TinggiBadan)
	}
	if req.Kontraksi != KontraksiTidak {
		t.Errorf("expected default Kontraksi to be Tidak, got %s", req.Kontraksi)
	}
	if req.TinggiFundusUteri != "30 cm" || req.TaksiranBeratJanin != "3000 gr" || req.His != "3x/10m/40s" || req.DenyutJantungJanin != "140 dpm" {
		t.Errorf("expected trimmed obgyn fields, got tfu=%s, tbj=%s, his=%s, djj=%s", req.TinggiFundusUteri, req.TaksiranBeratJanin, req.His, req.DenyutJantungJanin)
	}

	if req.Kepala != StatusFisikNormal || req.Mata != StatusFisikNormal || req.Gigi != StatusFisikNormal ||
		req.TelingaHidungTenggorok != StatusFisikNormal || req.Thoraks != StatusFisikNormal ||
		req.Abdomen != StatusFisikNormal || req.Genital != StatusFisikNormal ||
		req.Ekstremitas != StatusFisikNormal || req.Kulit != StatusFisikNormal {
		t.Errorf("expected default Pemeriksaan Fisik to be Normal, got mata=%s, kulit=%s", req.Mata, req.Kulit)
	}
}

func TestSimpanPenilaianMedisRalanKandunganRequest_Validation(t *testing.T) {

	req := SimpanPenilaianMedisRalanKandunganRequest{}
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

	req = SimpanPenilaianMedisRalanKandunganRequest{
		NoRawat: "2026/04/22/000004",
		DataPenilaianMedisRalanKandungan: DataPenilaianMedisRalanKandungan{
			KeluhanUtama: "Mules-mules",
			Diagnosis:    "G1P0A0 inpartu",
			TataLaksana:  "Observasi his dan DJJ",
			Kontraksi:    "Sering",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["kontraksi"] == "" {
		t.Errorf("expected error on invalid kontraksi, got: %+v", errs)
	}

	req.Kontraksi = KontraksiAda
	req.TanggalPenilaian = "2026-04-22 09:30:00"
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid Ralan Kandungan request to have no errors, got: %+v", errs)
	}
}

func TestUpdatePenilaianMedisRalanKandunganRequest_Validation(t *testing.T) {
	req := UpdatePenilaianMedisRalanKandunganRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty update Ralan Kandungan request")
	}

	req = UpdatePenilaianMedisRalanKandunganRequest{
		DataPenilaianMedisRalanKandungan: DataPenilaianMedisRalanKandungan{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Mules-mules berkurang",
			Diagnosis:        "G1P0A0 belum inpartu",
			TataLaksana:      "Rawat jalan, kontrol ulang 3 hari",
			Kontraksi:        KontraksiTidak,
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid update Ralan Kandungan request to have no errors, got: %+v", errs)
	}
}

func TestSimpanPenilaianMedisRanapKandunganRequest_Sanitize_Defaults(t *testing.T) {
	req := SimpanPenilaianMedisRanapKandunganRequest{
		NoRawat: " 2026/04/22/000005 ",
		DataPenilaianMedisRanapKandungan: DataPenilaianMedisRanapKandungan{
			KeluhanUtama:       " Perut kencang-kencang teratur ",
			SuhuTubuh:          " 37,2 ",
			BeratBadan:         " 65,5 ",
			TinggiBadan:        " 160,0 ",
			Diagnosis:          " G2P1A0 hamil aterm inpartu kala I fase aktif ",
			TataLaksana:        " Siapkan partus set, pantau his dan djj ",
			Edukasi:            " Edukasi teknik relaksasi saat kontraksi ",
			TinggiFundusUteri:  " 32 cm ",
			TaksiranBeratJanin: " 3200 gr ",
			His:                " 4x/10m/45s ",
			DenyutJantungJanin: " 144 dpm ",
		},
	}

	req.Sanitize()

	if req.NoRawat != "2026/04/22/000005" {
		t.Errorf("expected trimmed no_rawat, got %s", req.NoRawat)
	}
	if req.TanggalPenilaian == "" {
		t.Error("expected default TanggalPenilaian to be populated")
	}
	if req.KeluhanUtama != "Perut kencang-kencang teratur" {
		t.Errorf("expected trimmed keluhan utama, got %s", req.KeluhanUtama)
	}
	if req.Edukasi != "Edukasi teknik relaksasi saat kontraksi" {
		t.Errorf("expected trimmed edukasi, got %s", req.Edukasi)
	}
	if req.SuhuTubuh != "37.2" || req.BeratBadan != "65.5" || req.TinggiBadan != "160.0" {
		t.Errorf("expected comma replaced with dot in TTV: suhu=%s, bb=%s, tb=%s", req.SuhuTubuh, req.BeratBadan, req.TinggiBadan)
	}
	if req.Kontraksi != KontraksiTidak {
		t.Errorf("expected default Kontraksi to be Tidak, got %s", req.Kontraksi)
	}
	if req.TinggiFundusUteri != "32 cm" || req.TaksiranBeratJanin != "3200 gr" || req.His != "4x/10m/45s" || req.DenyutJantungJanin != "144 dpm" {
		t.Errorf("expected trimmed obgyn fields, got tfu=%s, tbj=%s, his=%s, djj=%s", req.TinggiFundusUteri, req.TaksiranBeratJanin, req.His, req.DenyutJantungJanin)
	}

	if req.Kepala != StatusFisikNormal || req.Mata != StatusFisikNormal || req.Gigi != StatusFisikNormal ||
		req.TelingaHidungTenggorok != StatusFisikNormal || req.Thoraks != StatusFisikNormal ||
		req.Jantung != StatusFisikNormal || req.Paru != StatusFisikNormal ||
		req.Abdomen != StatusFisikNormal || req.Genital != StatusFisikNormal ||
		req.Ekstremitas != StatusFisikNormal || req.Kulit != StatusFisikNormal {
		t.Errorf("expected default Pemeriksaan Fisik to be Normal, got jantung=%s, paru=%s", req.Jantung, req.Paru)
	}
}

func TestSimpanPenilaianMedisRanapKandunganRequest_Validation(t *testing.T) {

	req := SimpanPenilaianMedisRanapKandunganRequest{}
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

	req = SimpanPenilaianMedisRanapKandunganRequest{
		NoRawat: "2026/04/22/000005",
		DataPenilaianMedisRanapKandungan: DataPenilaianMedisRanapKandungan{
			KeluhanUtama: "Mules-mules",
			Diagnosis:    "G2P1A0 inpartu",
			TataLaksana:  "Observasi his dan DJJ",
			Kontraksi:    "Sering",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs == nil || errs["kontraksi"] == "" {
		t.Errorf("expected error on invalid kontraksi, got: %+v", errs)
	}

	req.Kontraksi = KontraksiAda
	req.TanggalPenilaian = "2026-04-22 09:30:00"
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid Ranap Kandungan request to have no errors, got: %+v", errs)
	}
}

func TestUpdatePenilaianMedisRanapKandunganRequest_Validation(t *testing.T) {
	req := UpdatePenilaianMedisRanapKandunganRequest{}
	errs := req.Validate()
	if errs == nil {
		t.Fatal("expected validation errors for empty update Ranap Kandungan request")
	}

	req = UpdatePenilaianMedisRanapKandunganRequest{
		DataPenilaianMedisRanapKandungan: DataPenilaianMedisRanapKandungan{
			TanggalPenilaian: "2026-04-22 09:30:00",
			KeluhanUtama:     "Mules-mules berkurang",
			Diagnosis:        "G2P1A0",
			TataLaksana:      "Observasi lanjut",
			Kontraksi:        KontraksiTidak,
			Edukasi:          "Tetap tirah baring",
		},
	}
	req.Sanitize()
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected valid update Ranap Kandungan request to have no errors, got: %+v", errs)
	}
}
