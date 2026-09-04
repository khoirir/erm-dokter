package rawatinap_test

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
)

type mockRepository struct {
	cekStatusKamarInapFn func(ctx context.Context, noRawat string) (bool, bool, error)
}

func (m *mockRepository) CekStatusKamarInap(ctx context.Context, noRawat string) (bool, bool, error) {
	if m.cekStatusKamarInapFn != nil {
		return m.cekStatusKamarInapFn(ctx, noRawat)
	}
	return false, false, nil
}

func TestService_CekStatusKamarInap(t *testing.T) {
	log := logger.New()

	t.Run("Empty no_rawat returns false, false, nil without calling repo", func(t *testing.T) {
		mockRepo := &mockRepository{
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				t.Fatal("Repo should not be called for empty no_rawat")
				return false, false, nil
			},
		}
		svc := rawatinap.NewService(mockRepo, log)

		isAktif, hasRecord, err := svc.CekStatusKamarInap(context.Background(), "   ")
		if err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}
		if isAktif || hasRecord {
			t.Fatalf("Expected false, false, got %v, %v", isAktif, hasRecord)
		}
	})

	t.Run("Active inpatient returns true, true", func(t *testing.T) {
		mockRepo := &mockRepository{
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				if noRawat != "2026/09/04/000001" {
					t.Fatalf("Unexpected noRawat %s", noRawat)
				}
				return true, true, nil
			},
		}
		svc := rawatinap.NewService(mockRepo, log)

		isAktif, hasRecord, err := svc.CekStatusKamarInap(context.Background(), "2026/09/04/000001")
		if err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}
		if !isAktif || !hasRecord {
			t.Fatalf("Expected true, true, got %v, %v", isAktif, hasRecord)
		}
	})

	t.Run("Checked-out inpatient returns false, true", func(t *testing.T) {
		mockRepo := &mockRepository{
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, true, nil
			},
		}
		svc := rawatinap.NewService(mockRepo, log)

		isAktif, hasRecord, err := svc.CekStatusKamarInap(context.Background(), "2026/09/04/000001")
		if err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}
		if isAktif || !hasRecord {
			t.Fatalf("Expected false, true, got %v, %v", isAktif, hasRecord)
		}
	})

	t.Run("Non-inpatient returns false, false", func(t *testing.T) {
		mockRepo := &mockRepository{
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, nil
			},
		}
		svc := rawatinap.NewService(mockRepo, log)

		isAktif, hasRecord, err := svc.CekStatusKamarInap(context.Background(), "2026/09/04/000001")
		if err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}
		if isAktif || hasRecord {
			t.Fatalf("Expected false, false, got %v, %v", isAktif, hasRecord)
		}
	})

	t.Run("Database error propagated", func(t *testing.T) {
		dbErr := errors.New("db connection failure")
		mockRepo := &mockRepository{
			cekStatusKamarInapFn: func(ctx context.Context, noRawat string) (bool, bool, error) {
				return false, false, dbErr
			},
		}
		svc := rawatinap.NewService(mockRepo, log)

		_, _, err := svc.CekStatusKamarInap(context.Background(), "2026/09/04/000001")
		if !errors.Is(err, dbErr) {
			t.Fatalf("Expected %v, got %v", dbErr, err)
		}
	})
}
