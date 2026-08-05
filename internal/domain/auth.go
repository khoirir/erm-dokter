package domain

import (
	"context"
)

type User struct {
	KodeDokter string `json:"kode_dokter"`
	NamaUser   string `json:"nama_user"`
	Password   string `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}

type AuthRepository interface {
	CariByUsername(ctx context.Context, username string) (*User, error)
}
type AuthUsecase interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}
