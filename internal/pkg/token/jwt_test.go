package token

import (
	"testing"
	"time"
)

func TestToken_GenerateValidateExtract(t *testing.T) {
	secret := "mysecretkey12345678901234567890"
	dokter := "DR001"
	nama := "dr. Handi"

	tkn, err := GenerateToken(dokter, nama, secret, 1*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := ValidateToken(tkn, secret)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}
	if claims.KodeDokter != dokter || claims.NamaUser != nama {
		t.Errorf("expected claims %s/%s, got %s/%s", dokter, nama, claims.KodeDokter, claims.NamaUser)
	}

	unverified, err := ExtractClaimsUnverified(tkn)
	if err != nil {
		t.Fatalf("unexpected error extracting unverified claims: %v", err)
	}
	if unverified.KodeDokter != dokter || unverified.NamaUser != nama {
		t.Errorf("expected unverified claims %s/%s, got %s/%s", dokter, nama, unverified.KodeDokter, unverified.NamaUser)
	}
}
