package auth

import (
	"context"
	"fmt"
	"time"

	"erm-dokter/internal/pkg/logger"
	"erm-dokter/internal/pkg/token"
	"erm-dokter/internal/shared/apperror"
)

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

type service struct {
	repo      Repository
	log       *logger.Logger
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string, log *logger.Logger) Service {
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.VerifikasiLogin(ctx, req.Username, req.Password)
	if err != nil {
		s.log.Error("Gagal login: %v", err)
		return nil, err
	}

	if user == nil {
		return nil, apperror.NewUnauthorizedError(
			"Username atau password salah",
			fmt.Sprintf("Kredensial login tidak cocok untuk username '%s'", req.Username),
		)
	}

	tkn, err := token.GenerateToken(user.IDUser, user.NamaUser, s.jwtSecret, 24*time.Hour)
	if err != nil {
		s.log.Error("Gagal membuat token: %v", err)
		return nil, err
	}

	return &LoginResponse{
		Token:      tkn,
		KodeDokter: user.IDUser,
		NamaDokter: user.NamaUser,
	}, nil
}
