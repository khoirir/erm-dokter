package docs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"erm-dokter/internal/docs"
	"gopkg.in/yaml.v3"
)

func TestMergeOpenAPISpecs(t *testing.T) {
	baseBytes, err := os.ReadFile("base.yaml")
	if err != nil {
		t.Fatalf("failed to read base.yaml: %v", err)
	}

	entries, err := os.ReadDir("modules")
	if err != nil {
		t.Fatalf("failed to read modules dir: %v", err)
	}

	var moduleFiles [][]byte
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".yaml") {
			content, err := os.ReadFile(filepath.Join("modules", entry.Name()))
			if err != nil {
				t.Fatalf("failed to read %s: %v", entry.Name(), err)
			}
			moduleFiles = append(moduleFiles, content)
		}
	}

	merged, err := docs.MergeOpenAPISpecs(baseBytes, moduleFiles)
	if err != nil {
		t.Fatalf("failed to merge openapi specs: %v", err)
	}

	// Verifikasi hasil valid YAML
	var root yaml.Node
	if err := yaml.Unmarshal(merged, &root); err != nil {
		t.Fatalf("merged result is not valid yaml: %v", err)
	}

	mergedStr := string(merged)

	// Pastikan metadata dasar ada
	if !strings.Contains(mergedStr, "openapi: 3.1.0") {
		t.Error("expected openapi: 3.1.0")
	}
	if !strings.Contains(mergedStr, "title: ERM Dokter API") {
		t.Error("expected title: ERM Dokter API")
	}

	// Pastikan endpoint dari berbagai modul ada
	expectedPaths := []string{
		"/health:",
		"/api/v1/auth/login:",
		"/api/v1/master/poliklinik:",
		"/api/v1/rawat-jalan/antrean:",
		"/api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}:",
		"/api/v1/resep/{id_kunjungan}/{status_lanjut}:",
		"/api/v1/penilaian-medis/ranap-kandungan/{id_kunjungan}:",
		"/api/v1/laboratorium/pk/permintaan/{id_kunjungan}/{status_lanjut}:",
		"/api/v1/radiologi/{id_kunjungan}/{status_lanjut}:",
		"/api/v1/berkas-digital/{id_berkas}:",
	}
	for _, p := range expectedPaths {
		if !strings.Contains(mergedStr, p) {
			t.Errorf("expected merged spec to contain path '%s'", p)
		}
	}

	// Pastikan skema dari berbagai modul ada
	expectedSchemas := []string{
		"BaseResponse:",
		"LoginRequest:",
		"Poliklinik:",
		"KunjunganRawatJalan:",
		"Pemeriksaan:",
		"Resep:",
		"PenilaianMedisRanapKandungan:",
		"ItemPemeriksaanLabPKRequest:",
		"ItemPemeriksaanRadiologiRequest:",
		"BerkasDigital:",
	}
	for _, s := range expectedSchemas {
		if !strings.Contains(mergedStr, s) {
			t.Errorf("expected merged spec to contain schema '%s'", s)
		}
	}
}
