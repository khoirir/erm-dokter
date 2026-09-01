package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	KodeDokter string `json:"kode_dokter"`
	NamaUser   string `json:"nama_user"`
	jwt.RegisteredClaims
}

func GenerateToken(kodeDokter, namaUser, secretKey string, duration time.Duration) (string, error) {
	claims := Claims{
		KodeDokter: kodeDokter,
		NamaUser:   namaUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(secretKey))
}

func ValidateToken(tokenStr, secretKey string) (*Claims, error) {
	tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode enkripsi token tidak valid")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := tkn.Claims.(*Claims)
	if !ok || !tkn.Valid {
		return nil, errors.New("token tidak valid atau sudah kedaluwarsa")
	}

	return claims, nil
}

func ExtractClaimsUnverified(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, _, err := jwt.NewParser().ParseUnverified(tokenStr, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

