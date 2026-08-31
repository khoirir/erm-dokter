package rujukaninternal

import (
	"testing"
)

func TestIdOpsiPoliDokter_CompositeKey(t *testing.T) {
	id := IdOpsiPoliDokter{
		KodePoli:   "INT",
		KodeDokter: "DR001",
	}

	key := id.CompositeKey()
	expected := "INT~DR001"
	if key != expected {
		t.Fatalf("expected composite key %s, got %s", expected, key)
	}

	parsed, err := ParseIdOpsiPoliDokter(key)
	if err != nil {
		t.Fatalf("unexpected error parsing key: %v", err)
	}
	if parsed.KodePoli != "INT" || parsed.KodeDokter != "DR001" {
		t.Fatalf("unexpected parsed values: %+v", parsed)
	}

	_, err = ParseIdOpsiPoliDokter("invalid_key")
	if err == nil {
		t.Fatal("expected error on invalid key, got nil")
	}
}

func TestIdRujukanInternal_CompositeKey(t *testing.T) {
	id := IdRujukanInternal{
		NoRawat:    "2026/04/22/000001",
		KodePoli:   "INT",
		KodeDokter: "DR001",
	}

	key := id.CompositeKey()
	expected := "2026/04/22/000001~INT~DR001"
	if key != expected {
		t.Fatalf("expected composite key %s, got %s", expected, key)
	}

	parsed, err := ParseIdRujukanInternal(key)
	if err != nil {
		t.Fatalf("unexpected error parsing key: %v", err)
	}
	if parsed.NoRawat != "2026/04/22/000001" || parsed.KodePoli != "INT" || parsed.KodeDokter != "DR001" {
		t.Fatalf("unexpected parsed values: %+v", parsed)
	}

	_, err = ParseIdRujukanInternal("invalid~key")
	if err == nil {
		t.Fatal("expected error on invalid key, got nil")
	}
}

func TestSimpanRujukanRequest_Validation(t *testing.T) {
	req := SimpanRujukanRequest{
		IdTujuan: "",
	}
	errs := req.Validate()
	if errs == nil || errs["id_tujuan"] == "" {
		t.Fatal("expected validation error on empty id_tujuan")
	}

	req.IdTujuan = " valid_encrypted_id "
	errs = req.Validate()
	if errs != nil {
		t.Fatalf("expected no validation error, got %+v", errs)
	}
	if req.IdTujuan != "valid_encrypted_id" {
		t.Fatalf("expected sanitized id_tujuan, got %s", req.IdTujuan)
	}
}
