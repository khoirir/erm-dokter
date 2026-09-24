package eklaim

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	signatureSize = 10
	ivSize        = 16
	keySize       = 32
)

func Encrypt(data []byte, keyHex string) (string, error) {
	key, err := hex.DecodeString(strings.TrimSpace(keyHex))
	if err != nil {
		return "", fmt.Errorf("kunci enkripsi hex tidak valid: %w", err)
	}
	if len(key) != keySize {
		return "", errors.New("kunci enkripsi E-Klaim harus 256-bit (32 byte)")
	}

	iv := make([]byte, ivSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("gagal membuat IV acak: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("gagal menginisialisasi cipher AES: %w", err)
	}

	padded := pkcs7Pad(data, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	h := hmac.New(sha256.New, key)
	h.Write(ciphertext)
	mac := h.Sum(nil)
	signature := mac[:signatureSize]

	combined := make([]byte, 0, len(signature)+len(iv)+len(ciphertext))
	combined = append(combined, signature...)
	combined = append(combined, iv...)
	combined = append(combined, ciphertext...)

	return base64.StdEncoding.EncodeToString(combined), nil
}

func Decrypt(b64Str string, keyHex string) ([]byte, error) {
	key, err := hex.DecodeString(strings.TrimSpace(keyHex))
	if err != nil {
		return nil, fmt.Errorf("kunci enkripsi hex tidak valid: %w", err)
	}
	if len(key) != keySize {
		return nil, errors.New("kunci enkripsi E-Klaim harus 256-bit (32 byte)")
	}

	cleaned := cleanBase64Payload(b64Str)
	cleaned = strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' || r == ' ' || r == '\t' {
			return -1
		}
		return r
	}, cleaned)

	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return nil, fmt.Errorf("gagal decode base64 respons E-Klaim: %w", err)
	}

	minLen := signatureSize + ivSize
	if len(decoded) < minLen {
		return nil, errors.New("payload terenkripsi E-Klaim terlalu pendek")
	}

	signature := decoded[:signatureSize]
	iv := decoded[signatureSize : signatureSize+ivSize]
	ciphertext := decoded[signatureSize+ivSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("panjang ciphertext tidak kelipatan ukuran blok AES")
	}

	h := hmac.New(sha256.New, key)
	h.Write(ciphertext)
	expectedMac := h.Sum(nil)[:signatureSize]

	if !hmac.Equal(signature, expectedMac) {
		return nil, errors.New("SIGNATURE_NOT_MATCH")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi cipher AES: %w", err)
	}

	plainPadded := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainPadded, ciphertext)

	plain, err := pkcs7Unpad(plainPadded, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("gagal unpad data respons E-Klaim: %w", err)
	}

	return plain, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, errors.New("panjang data padding tidak valid")
	}
	padLen := int(data[length-1])
	if padLen == 0 || padLen > blockSize || padLen > length {
		return nil, errors.New("nilai byte padding tidak valid")
	}
	for i := length - padLen; i < length; i++ {
		if data[i] != byte(padLen) {
			return nil, errors.New("format byte padding tidak cocok")
		}
	}
	return data[:length-padLen], nil
}

func cleanBase64Payload(input string) string {
	lines := strings.Split(input, "\n")
	var sb strings.Builder
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "ENCRYPTED DATA") || strings.HasPrefix(trimmed, "---") || strings.HasSuffix(trimmed, "---") {
			continue
		}
		sb.WriteString(trimmed)
	}
	return sb.String()
}
