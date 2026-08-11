package usecase

import (
	"context"
	"fmt"
	"time"

	"erm-dokter/internal/domain"
	"erm-dokter/internal/dto"
	"erm-dokter/pkg/logger"
	"erm-dokter/pkg/token"

	"github.com/go-playground/validator/v10"
)

type authUsecase struct {
	authRepo  domain.AuthRepository
	Log       *logger.Logger
	Validate  *validator.Validate
	jwtSecret string
}

func NewAuthUsecase(repo domain.AuthRepository, jwtSecret string, validate *validator.Validate, log *logger.Logger) domain.AuthUsecase {
	return &authUsecase{
		authRepo:  repo,
		jwtSecret: jwtSecret,
		Validate:  validate,
		Log:       log,
	}
}

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		u.Log.Warn("Invalid request body: %+v", err)
		return nil, domain.NewBusinessError("username dan password wajib diisi")
	}

	user, err := u.authRepo.VerifikasiLogin(ctx, req.Username, req.Password)
	if err != nil {
		u.Log.Error("Gagal memverifikasi login: %v", err)
		return nil, fmt.Errorf("gagal memverifikasi login: %w", err)
	}

	if user == nil {
		u.Log.Warn("Login gagal untuk username: %s", req.Username)
		return nil, domain.NewBusinessError("username atau password salah")
	}

	tkn, err := token.GenerateToken(user.IDUser, user.NamaUser, u.jwtSecret, 24*time.Hour)
	if err != nil {
		u.Log.Error("Gagal membuat token: %v", err)
		return nil, fmt.Errorf("gagal membuat token: %w", err)
	}

	u.Log.Info("Login berhasil untuk dokter: %s", user.IDUser)
	return &dto.LoginResponse{
		Token:      tkn,
		KodeDokter: user.IDUser,
		NamaDokter: user.NamaUser,
	}, nil
}
