package crypto_test

import (
	"testing"

	"erm-dokter/internal/pkg/crypto"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := "my-secret-key-for-testing"
	testCases := []struct {
		name      string
		plainText string
	}{
		{"simple text", "hello world"},
		{"no_rawat format", "2025/04/22/000001"},
		{"empty string", ""},
		{"special chars", "test/with+special=chars&more"},
		{"unicode", "Pasien: Budi Santoso (RM-001234)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := crypto.Encrypt(tc.plainText, key)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			if encrypted == tc.plainText && tc.plainText != "" {
				t.Error("Encrypted text should differ from plain text")
			}

			decrypted, err := crypto.Decrypt(encrypted, key)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if decrypted != tc.plainText {
				t.Errorf("Roundtrip failed: got %q, want %q", decrypted, tc.plainText)
			}
		})
	}
}

func TestEncryptWithEmptyKey(t *testing.T) {
	_, err := crypto.Encrypt("test", "")
	if err == nil {
		t.Error("Encrypt() with empty key should return error")
	}
}

func TestDecryptWithEmptyKey(t *testing.T) {
	_, err := crypto.Decrypt("test", "")
	if err == nil {
		t.Error("Decrypt() with empty key should return error")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	encrypted, err := crypto.Encrypt("hello", "correct-key")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	_, err = crypto.Decrypt(encrypted, "wrong-key")
	if err == nil {
		t.Error("Decrypt() with wrong key should return error")
	}
}

func TestDecryptInvalidCipherText(t *testing.T) {
	_, err := crypto.Decrypt("not-valid-base64!!!", "any-key")
	if err == nil {
		t.Error("Decrypt() with invalid cipher text should return error")
	}
}
