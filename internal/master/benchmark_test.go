package master_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"erm-dokter/internal/master"
	"erm-dokter/internal/shared"
)

func BenchmarkMasterHandler_DaftarICD10(b *testing.B) {
	// Siapkan 20 data dummy ICD10
	dummyList := make([]master.ICD10, 20)
	for i := 0; i < 20; i++ {
		dummyList[i] = master.ICD10{
			ItemMaster: master.ItemMaster{
				Kode: fmt.Sprintf("I%02d.9", i),
				Nama: fmt.Sprintf("Diagnosis Penyakit ICD-10 Nomor %d", i),
			},
		}
	}

	mockSvc := &mockMasterService{
		daftarICD10Fn: func(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD10, shared.PaginationMeta, error) {
			return dummyList, shared.NewPaginationMeta(15000, filter.Page, filter.Limit), nil
		},
	}

	handler := master.NewHandler(mockSvc)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, noOpMw, noOpMw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd10?keyword=stroke&page=1&limit=20", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				b.Fatalf("expected 200, got %d", rr.Code)
			}
		}
	})
}

func BenchmarkMasterHandler_DaftarICD9(b *testing.B) {
	dummyList := make([]master.ICD9, 20)
	for i := 0; i < 20; i++ {
		dummyList[i] = master.ICD9{
			ItemMaster: master.ItemMaster{
				Kode: fmt.Sprintf("87.%02d", i),
				Nama: fmt.Sprintf("Prosedur Tindakan Medis ICD-9 Nomor %d", i),
			},
		}
	}

	mockSvc := &mockMasterService{
		daftarICD9Fn: func(ctx context.Context, filter master.FilterMasterICD) ([]master.ICD9, shared.PaginationMeta, error) {
			return dummyList, shared.NewPaginationMeta(4000, filter.Page, filter.Limit), nil
		},
	}

	handler := master.NewHandler(mockSvc)
	mux := http.NewServeMux()
	noOpMw := func(next http.HandlerFunc) http.HandlerFunc { return next }
	handler.RegisterRoutes(mux, noOpMw, noOpMw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd9?keyword=tomo&page=1&limit=20", nil)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				b.Fatalf("expected 200, got %d", rr.Code)
			}
		}
	})
}
