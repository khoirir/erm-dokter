package master_test

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/master"
	"erm-dokter/internal/pkg/logger"
)

type mockRepository struct {
	penjaminData   []master.Penjamin
	depoData       []master.Depo
	poliklinikData []master.Poliklinik
	bangsalData    []master.Bangsal
	kelasData      []master.KelasKamar
	icd10Data      []master.ICD10
	icd10Total     int
	icd9Data       []master.ICD9
	icd9Total      int
	checkICD10Map  map[string]bool
	checkICD9Map   map[string]bool
	err            error
}

func (m *mockRepository) DaftarPenjamin(ctx context.Context) ([]master.Penjamin, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.penjaminData, nil
}

func (m *mockRepository) DaftarDepo(ctx context.Context) ([]master.Depo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.depoData, nil
}

func (m *mockRepository) DaftarPoliklinik(ctx context.Context) ([]master.Poliklinik, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.poliklinikData, nil
}

func (m *mockRepository) DaftarBangsal(ctx context.Context) ([]master.Bangsal, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.bangsalData, nil
}

func (m *mockRepository) DaftarKelas(ctx context.Context) ([]master.KelasKamar, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.kelasData, nil
}

func (m *mockRepository) DaftarICD10(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD10, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.icd10Data, m.icd10Total, nil
}

func (m *mockRepository) DaftarICD9(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD9, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.icd9Data, m.icd9Total, nil
}

func (m *mockRepository) FetchAllICD10(ctx context.Context) ([]master.ICD10, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.icd10Data, nil
}

func (m *mockRepository) FetchAllICD9(ctx context.Context) ([]master.ICD9, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.icd9Data, nil
}

func (m *mockRepository) CekKeberadaanICD10(ctx context.Context, listKode []string) (map[string]bool, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.checkICD10Map, nil
}

func (m *mockRepository) CekKeberadaanICD9(ctx context.Context, listKode []string) (map[string]bool, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.checkICD9Map, nil
}

func TestMasterService_DaftarPenjamin(t *testing.T) {
	repo := &mockRepository{
		penjaminData: []master.Penjamin{
			{ItemMaster: master.ItemMaster{Kode: "BPJ", Nama: "BPJS KESEHATAN"}},
		},
	}
	svc := master.NewService(repo, logger.New())

	data, err := svc.DaftarPenjamin(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "BPJ" {
		t.Errorf("unexpected penjamin data: %+v", data)
	}

	repo.err = errors.New("db error")
	_, err = svc.DaftarPenjamin(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_DaftarDepo(t *testing.T) {
	repo := &mockRepository{
		depoData: []master.Depo{
			{ItemMaster: master.ItemMaster{Kode: "DPRJ", Nama: "Depo Rawat Jalan"}},
		},
	}
	svc := master.NewService(repo, logger.New())

	data, err := svc.DaftarDepo(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "DPRJ" {
		t.Errorf("unexpected depo data: %+v", data)
	}

	repo.err = errors.New("db error")
	_, err = svc.DaftarDepo(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_DaftarPoliklinik(t *testing.T) {
	repo := &mockRepository{
		poliklinikData: []master.Poliklinik{
			{ItemMaster: master.ItemMaster{Kode: "INT", Nama: "Poli Penyakit Dalam"}},
		},
	}
	svc := master.NewService(repo, logger.New())

	data, err := svc.DaftarPoliklinik(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "INT" {
		t.Errorf("unexpected poliklinik data: %+v", data)
	}

	repo.err = errors.New("db error")
	_, err = svc.DaftarPoliklinik(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_DaftarBangsal(t *testing.T) {
	repo := &mockRepository{
		bangsalData: []master.Bangsal{
			{ItemMaster: master.ItemMaster{Kode: "B01", Nama: "Melati"}},
		},
	}
	svc := master.NewService(repo, logger.New())

	data, err := svc.DaftarBangsal(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "B01" || data[0].Nama != "Melati" {
		t.Errorf("unexpected bangsal data: %+v", data)
	}

	repo.err = errors.New("db error")
	_, err = svc.DaftarBangsal(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_DaftarKelas(t *testing.T) {
	repo := &mockRepository{
		kelasData: []master.KelasKamar{
			{ItemMaster: master.ItemMaster{Kode: "Kelas 1", Nama: "Kelas 1"}},
		},
	}
	svc := master.NewService(repo, logger.New())

	data, err := svc.DaftarKelas(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "Kelas 1" || data[0].Nama != "Kelas 1" {
		t.Errorf("unexpected kelas data: %+v", data)
	}

	repo.err = errors.New("db error")
	_, err = svc.DaftarKelas(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_DaftarICD10(t *testing.T) {
	repo := &mockRepository{
		icd10Data: []master.ICD10{
			{ItemMaster: master.ItemMaster{Kode: "I63.9", Nama: "Cerebral infarction, unspecified"}},
		},
		icd10Total: 1,
	}
	svc := master.NewService(repo, logger.New())

	filter := master.FilterMasterICD{Keyword: "cerebral", Page: 1, Limit: 20}
	data, meta, err := svc.DaftarICD10(context.Background(), filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "I63.9" {
		t.Errorf("unexpected ICD-10 data: %+v", data)
	}
	if meta.TotalRecords != 1 || meta.CurrentPage != 1 {
		t.Errorf("unexpected meta: %+v", meta)
	}

	// Test fallback ke DB jika cache gagal dimuat
	repoErr := &mockRepository{err: errors.New("db error")}
	svcErr := master.NewService(repoErr, logger.New())
	_, _, err = svcErr.DaftarICD10(context.Background(), filter)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_DaftarICD9(t *testing.T) {
	repo := &mockRepository{
		icd9Data: []master.ICD9{
			{ItemMaster: master.ItemMaster{Kode: "87.03", Nama: "Computerized axial tomography of head"}},
		},
		icd9Total: 1,
	}
	svc := master.NewService(repo, logger.New())

	filter := master.FilterMasterICD{Keyword: "tomography", Page: 1, Limit: 20}
	data, meta, err := svc.DaftarICD9(context.Background(), filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) != 1 || data[0].Kode != "87.03" {
		t.Errorf("unexpected ICD-9 data: %+v", data)
	}
	if meta.TotalRecords != 1 || meta.CurrentPage != 1 {
		t.Errorf("unexpected meta: %+v", meta)
	}

	// Test fallback ke DB jika cache gagal dimuat
	repoErr := &mockRepository{err: errors.New("db error")}
	svcErr := master.NewService(repoErr, logger.New())
	_, _, err = svcErr.DaftarICD9(context.Background(), filter)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_CekKeberadaanICD(t *testing.T) {
	repo := &mockRepository{
		icd10Data: []master.ICD10{
			{ItemMaster: master.ItemMaster{Kode: "I63.9", Nama: "Cerebral infarction"}},
		},
		icd9Data: []master.ICD9{
			{ItemMaster: master.ItemMaster{Kode: "87.03", Nama: "CT Head"}},
		},
		checkICD10Map: map[string]bool{"I63.9": true, "A00": false},
		checkICD9Map:  map[string]bool{"87.03": true, "99.99": false},
	}
	svc := master.NewService(repo, logger.New())

	res10, err := svc.CekKeberadaanICD10(context.Background(), []string{"I63.9", "A00"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res10["I63.9"] || res10["A00"] {
		t.Errorf("unexpected res10: %+v", res10)
	}

	res9, err := svc.CekKeberadaanICD9(context.Background(), []string{"87.03", "99.99"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res9["87.03"] || res9["99.99"] {
		t.Errorf("unexpected res9: %+v", res9)
	}

	// Test fallback jika cache belum loaded dan DB error
	repoErr := &mockRepository{err: errors.New("db error")}
	svcErr := master.NewService(repoErr, logger.New())
	_, err = svcErr.CekKeberadaanICD10(context.Background(), []string{"I63.9"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	_, err = svcErr.CekKeberadaanICD9(context.Background(), []string{"87.03"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMasterService_SyncICD(t *testing.T) {
	repo := &mockRepository{
		icd10Data: []master.ICD10{
			{ItemMaster: master.ItemMaster{Kode: "I63.9", Nama: "Cerebral infarction"}},
		},
		icd9Data: []master.ICD9{
			{ItemMaster: master.ItemMaster{Kode: "87.03", Nama: "CT Head"}},
		},
	}
	svc := master.NewService(repo, logger.New())

	res, err := svc.SyncICD(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.TotalICD10 != 1 || res.TotalICD9 != 1 {
		t.Errorf("unexpected sync result: %+v", res)
	}

	repo.err = errors.New("db error")
	_, err = svc.SyncICD(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}


