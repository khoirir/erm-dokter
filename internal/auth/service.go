package auth

import (
	"context"
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
		s.log.Error("Gagal memverifikasi login: %v", err)
		return nil, err
	}

	if user == nil {
		s.log.Warn("Login gagal untuk username: %s", req.Username)
		return nil, apperror.NewBusinessError("username atau password salah")
	}

	tkn, err := token.GenerateToken(user.IDUser, user.NamaUser, s.jwtSecret, 24*time.Hour)
	if err != nil {
		s.log.Error("Gagal membuat token: %v", err)
		return nil, err
	}

	s.log.Info("Login berhasil untuk dokter: %s", user.IDUser)
	return &LoginResponse{
		Token:      tkn,
		KodeDokter: user.IDUser,
		NamaDokter: user.NamaUser,
	}, nil
}
