package usecase

import (
	"context"
	"errors"
	"time"

	"erm-dokter/internal/domain"
	"erm-dokter/pkg/token"

	"github.com/go-playground/validator/v10"
)

type authUsecase struct {
	authRepo  domain.AuthRepository
	jwtSecret string
	validate  *validator.Validate
}

func NewAuthUsecase(repo domain.AuthRepository, jwtSecret string) domain.AuthUsecase {
	return &authUsecase{
		authRepo:  repo,
		jwtSecret: jwtSecret,
		validate:  validator.New(),
	}
}

func (u *authUsecase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, errors.New("username dan password wajib diisi")
	}
	user, err := u.authRepo.CariByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("username atau password salah")
	}

	if req.Password != user.Password {
		return nil, errors.New("username atau password salah")
	}

	tkn, err := token.GenerateToken(user.KodeDokter, user.NamaUser, u.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.New("gagal membuat token autentikasi")
	}

	return &domain.LoginResponse{
		Token:      tkn,
		KodeDokter: user.KodeDokter,
		NamaDokter: user.NamaUser,
	}, nil
}
