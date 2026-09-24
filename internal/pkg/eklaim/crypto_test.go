package eklaim_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"erm-dokter/internal/pkg/eklaim"
)

const testKeyHex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	originalData := []byte(`{"metadata":{"method":"grouper","stage":"1"},"data":{"nomor_sep":"SIM-12345"}}`)

	encrypted, err := eklaim.Encrypt(originalData, testKeyHex)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == "" {
		t.Fatal("Encrypted string should not be empty")
	}

	decrypted, err := eklaim.Decrypt(encrypted, testKeyHex)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(originalData, decrypted) {
		t.Errorf("Decrypted data mismatch.\nExpected: %s\nGot: %s", string(originalData), string(decrypted))
	}
}

func TestDecrypt_ChunkSplitFormat(t *testing.T) {
	originalData := []byte(`{"test":"chunk_split_with_newlines"}`)

	encrypted, err := eklaim.Encrypt(originalData, testKeyHex)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	chunked := encrypted[:10] + "\r\n" + encrypted[10:30] + "\n" + encrypted[30:] + "\r\n"

	decrypted, err := eklaim.Decrypt(chunked, testKeyHex)
	if err != nil {
		t.Fatalf("Decrypt with chunked format failed: %v", err)
	}

	if !bytes.Equal(originalData, decrypted) {
		t.Errorf("Expected decrypted data match, got: %s", string(decrypted))
	}
}

func TestDecrypt_WithEncryptedDataBanner(t *testing.T) {
	originalData := []byte(`{"metadata":{"code":200,"message":"Ok"}}`)

	encrypted, err := eklaim.Encrypt(originalData, testKeyHex)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	bannered := "----BEGIN ENCRYPTED DATA----\r\n" + encrypted + "\r\n----END ENCRYPTED DATA----"

	decrypted, err := eklaim.Decrypt(bannered, testKeyHex)
	if err != nil {
		t.Fatalf("Decrypt with banner failed: %v", err)
	}

	if !bytes.Equal(originalData, decrypted) {
		t.Errorf("Expected decrypted data match, got: %s", string(decrypted))
	}
}

func TestDecrypt_ServerPayloadFormat(t *testing.T) {
	mockResponse := []byte(`{"metadata":{"code":400,"message":"Parameter tidak berlaku: diagnosa, procedure","error_no":"E2044"}}`)
	encrypted, err := eklaim.Encrypt(mockResponse, testKeyHex)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	payloadWithBanner := "----BEGIN ENCRYPTED DATA----\r\n" + encrypted + "\r\n----END ENCRYPTED DATA----"
	decrypted, err := eklaim.Decrypt(payloadWithBanner, testKeyHex)
	if err != nil {
		t.Fatalf("Decrypt server payload failed: %v", err)
	}

	if !bytes.Equal(mockResponse, decrypted) {
		t.Errorf("Expected %s, got %s", string(mockResponse), string(decrypted))
	}
}

func TestEncrypt_InvalidKeyLength(t *testing.T) {
	shortKey := "126658800f3b5001"
	_, err := eklaim.Encrypt([]byte("test"), shortKey)
	if err == nil {
		t.Fatal("Expected error for short key, got nil")
	}

	nonHexKey := "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"
	_, err = eklaim.Encrypt([]byte("test"), nonHexKey)
	if err == nil {
		t.Fatal("Expected error for non-hex key, got nil")
	}
}

func TestDecrypt_SignatureMismatch(t *testing.T) {
	originalData := []byte(`{"status":"secret"}`)

	encrypted, err := eklaim.Encrypt(originalData, testKeyHex)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	wrongKey := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	_, err = eklaim.Decrypt(encrypted, wrongKey)
	if err == nil {
		t.Fatal("Expected error for wrong key, got nil")
	}
	if err.Error() != "SIGNATURE_NOT_MATCH" {
		t.Errorf("Expected 'SIGNATURE_NOT_MATCH', got: %v", err)
	}
}

func TestDecrypt_CorruptedPayload(t *testing.T) {
	_, err := eklaim.Decrypt("invalid-base64!!", testKeyHex)
	if err == nil {
		t.Fatal("Expected error for invalid base64, got nil")
	}

	shortPayload := base64.StdEncoding.EncodeToString([]byte("short"))
	_, err = eklaim.Decrypt(shortPayload, testKeyHex)
	if err == nil {
		t.Fatal("Expected error for payload too short, got nil")
	}
}
