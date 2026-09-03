package obat_test

import (
	"testing"

	"erm-dokter/internal/obat"
)

func TestFilterDaftarObat_Validation(t *testing.T) {
	t.Run("Valid default filter", func(t *testing.T) {
		f := obat.FilterDaftarObat{
			Depo:  "DPRJ",
			Page:  1,
			Limit: 20,
		}
		if errs := f.Validate(); errs != nil {
			t.Errorf("Expected nil validation errors, got: %+v", errs)
		}
	})

	t.Run("Invalid OrderBy and SortOrder", func(t *testing.T) {
		f := obat.FilterDaftarObat{
			OrderBy:   "malicious_col",
			SortOrder: "DROP TABLE",
		}
		errs := f.Validate()
		if errs == nil {
			t.Fatal("Expected validation error for invalid order_by and sort_order")
		}
		if _, exists := errs["order_by"]; !exists {
			t.Error("Expected error on order_by")
		}
		if _, exists := errs["sort_order"]; !exists {
			t.Error("Expected error on sort_order")
		}
	})
}
