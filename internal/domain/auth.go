package domain

import (
	"context"

	"erm-dokter/internal/dto"
)

type User struct {
	IDUser   string `json:"id_user"`
	NamaUser string `json:"nama_user"`
}

type AuthRepository interface {
	VerifikasiLogin(ctx context.Context, username string, password string) (*User, error)
}

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}
