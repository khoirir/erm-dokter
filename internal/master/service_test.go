package master_test

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/master"
	"erm-dokter/internal/pkg/logger"
)

type mockRepository struct {
	penjaminData []master.Penjamin
	depoData     []master.Depo
	err          error
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

func TestMasterService_DaftarPenjamin(t *testing.T) {
	repo := &mockRepository{
		penjaminData: []master.Penjamin{
			{Kode: "BPJ", Nama: "BPJS KESEHATAN"},
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
			{Kode: "DPRJ", Nama: "Depo Rawat Jalan"},
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
