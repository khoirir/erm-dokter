package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

func get32ByteKey(encryptionKey string) []byte {
	hash := sha256.Sum256([]byte(encryptionKey))
	return hash[:]
}

func Encrypt(plainText string, encryptionKey string) (string, error) {
	if encryptionKey == "" {
		return "", errors.New("encryption key tidak boleh kosong")
	}
	key := get32ByteKey(encryptionKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return base64.RawURLEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cryptoText string, encryptionKey string) (string, error) {
	if encryptionKey == "" {
		return "", errors.New("encryption key tidak boleh kosong")
	}
	key := get32ByteKey(encryptionKey)

	cipherText, err := base64.RawURLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", errors.New("format token url tidak valid")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ukuran cipher text terlalu pendek")
	}

	nonce, cipherTextData := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherTextData, nil)
	if err != nil {
		return "", errors.New("gagal dekripsi data: key salah atau data rusak")
	}

	return string(plainText), nil
}
