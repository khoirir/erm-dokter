package resep

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"erm-dokter/internal/obat"
	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/rawatinap"
	"erm-dokter/internal/rawatjalan"
	"erm-dokter/internal/shared"
)

type mockConcurrentRepo struct {
	Repository
	simpanResepFunc func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanResepRequest) (*Resep, error)
}

func (m *mockConcurrentRepo) CekKeberadaanMetodeRacik(ctx context.Context, listKodeRacik []string) (map[string]bool, error) {
	return map[string]bool{"PULV": true}, nil
}

func (m *mockConcurrentRepo) SimpanResep(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanResepRequest) (*Resep, error) {
	if m.simpanResepFunc != nil {
		return m.simpanResepFunc(ctx, kodeDokter, statusLanjut, req)
	}
	return nil, nil
}

func (m *mockConcurrentRepo) HapusResep(ctx context.Context, noResep string) error {
	return nil
}

func (m *mockConcurrentRepo) UpdateResep(ctx context.Context, noResep string, req SimpanResepRequest) (*Resep, error) {
	return nil, nil
}

type mockConcurrentRJ struct {
	rawatjalan.Service
}

func (m *mockConcurrentRJ) GetWaktuRegistrasi(ctx context.Context, noRawat string) (string, string, bool, error) {
	return time.Now().Format("2006-01-02"), "07:00:00", true, nil
}

func (m *mockConcurrentRJ) GetInfoRegistrasi(ctx context.Context, noRawat string) (*rawatjalan.InfoRegistrasiPasien, error) {
	return &rawatjalan.InfoRegistrasiPasien{
		TanggalRegistrasi: time.Now().Format("2006-01-02"),
		JamRegistrasi:     "07:00:00",
		KodePenjamin:      "UMU",
		StatusBayar:       "Belum Bayar",
	}, nil
}

type mockConcurrentRI struct {
	rawatinap.Service
}

func (m *mockConcurrentRI) CekStatusKamarInap(ctx context.Context, noRawat string) (bool, bool, error) {
	return false, false, nil
}

type mockConcurrentObat struct {
	obat.Service
}

func (m *mockConcurrentObat) CekKeberadaanObat(ctx context.Context, listKodeObat []string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, k := range listKodeObat {
		res[k] = true
	}
	return res, nil
}

func TestSimpanResep_ConcurrentDoctors(t *testing.T) {
	const totalDokter = 20
	var sequenceCounter int64
	var mu sync.Mutex

	today := time.Now().Format("2006-01-02")

	mockRepo := &mockConcurrentRepo{
		simpanResepFunc: func(ctx context.Context, kodeDokter string, statusLanjut shared.StatusLanjut, req SimpanResepRequest) (*Resep, error) {
			mu.Lock()

			seq := atomic.AddInt64(&sequenceCounter, 1)
			noResep := fmt.Sprintf("%s%04d", time.Now().Format("20060102"), seq)
			mu.Unlock()

			return &Resep{
				NoResep:          noResep,
				NoRawat:          req.NoRawat,
				TanggalPeresepan: req.TanggalPeresepan,
				JamPeresepan:     req.JamPeresepan,
				KodeDokter:       kodeDokter,
				Status:           string(statusLanjut),
			}, nil
		},
	}

	mockRJ := &mockConcurrentRJ{}
	mockRI := &mockConcurrentRI{}
	mockObat := &mockConcurrentObat{}

	log := logger.New()
	svc := NewService(mockRepo, mockRJ, mockRI, mockObat, 48, log)

	var wg sync.WaitGroup
	startGate := make(chan struct{})
	generatedNumbers := sync.Map{}
	errorCount := int64(0)

	for i := 1; i <= totalDokter; i++ {
		wg.Add(1)
		go func(dokterIdx int) {
			defer wg.Done()
			<-startGate

			req := SimpanResepRequest{
				NoRawat:          fmt.Sprintf("2026/08/28/%06d", dokterIdx),
				TanggalPeresepan: today,
				JamPeresepan:     "08:00:00",
				ResepDokter: []ResepDokterInput{
					{
						ItemObatInput: ItemObatInput{
							KodeObat: "OBAT001",
							Jumlah:   10,
						},
						AturanPakai: "3x1",
					},
				},
			}

			kodeDokter := fmt.Sprintf("DK%03d", dokterIdx)
			result, err := svc.SimpanResep(context.Background(), kodeDokter, shared.StatusLanjutRawatJalan, req)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				t.Errorf("Dokter %s gagal simpan resep: %v", kodeDokter, err)
				return
			}

			if _, loaded := generatedNumbers.LoadOrStore(result.NoResep, true); loaded {
				t.Errorf("DUPLIKASI TERDETEKSI! NoResep kembar: %s", result.NoResep)
			}
		}(i)
	}

	close(startGate)
	wg.Wait()

	if errorCount > 0 {
		t.Fatalf("Terdapat %d transaksi gagal dari total %d dokter konkuren", errorCount, totalDokter)
	}

	totalUnik := 0
	generatedNumbers.Range(func(key, value any) bool {
		totalUnik++
		return true
	})

	if totalUnik != totalDokter {
		t.Fatalf("Diharapkan %d nomor resep unik, namun dihasilkan %d", totalDokter, totalUnik)
	}
}

func TestIsDuplicateKey(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "MySQL Error 1062 duplicate key",
			err:      errors.New("Error 1062 (23000): Duplicate entry '202608280001' for key 'resep_obat.PRIMARY'"),
			expected: true,
		},
		{
			name:     "Duplicate entry without error code",
			err:      errors.New("Duplicate entry '202608280001' for key 'PRIMARY'"),
			expected: true,
		},
		{
			name:     "Other database error (e.g. connection error)",
			err:      errors.New("connection refused"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicateKey(tt.err)
			if got != tt.expected {
				t.Errorf("isDuplicateKey(%v) = %v, expected %v", tt.err, got, tt.expected)
			}
		})
	}
}
