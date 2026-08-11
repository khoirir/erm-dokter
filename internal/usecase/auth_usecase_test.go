package usecase_test

import (
	"context"
	"errors"
	"testing"

	"erm-dokter/internal/domain"
	"erm-dokter/internal/dto"
	"erm-dokter/internal/usecase"
	"erm-dokter/pkg/logger"

	"github.com/go-playground/validator/v10"
)

// --- Mock Repository ---

type mockAuthRepository struct {
	verifikasiLoginFunc func(ctx context.Context, username, password string) (*domain.User, error)
}

func (m *mockAuthRepository) VerifikasiLogin(ctx context.Context, username, password string) (*domain.User, error) {
	if m.verifikasiLoginFunc != nil {
		return m.verifikasiLoginFunc(ctx, username, password)
	}
	return nil, nil
}

// --- Tests ---

func TestLogin_EmptyUsername(t *testing.T) {
	repo := &mockAuthRepository{}
	validate := validator.New()
	log := logger.New()
	uc := usecase.NewAuthUsecase(repo, "test-secret", validate, log)

	_, err := uc.Login(context.Background(), dto.LoginRequest{
		Username: "",
		Password: "pass123",
	})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLogin_EmptyPassword(t *testing.T) {
	repo := &mockAuthRepository{}
	validate := validator.New()
	log := logger.New()
	uc := usecase.NewAuthUsecase(repo, "test-secret", validate, log)

	_, err := uc.Login(context.Background(), dto.LoginRequest{
		Username: "dokter1",
		Password: "",
	})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLogin_WrongCredentials(t *testing.T) {
	repo := &mockAuthRepository{
		verifikasiLoginFunc: func(ctx context.Context, username, password string) (*domain.User, error) {
			return nil, nil
		},
	}
	validate := validator.New()
	log := logger.New()
	uc := usecase.NewAuthUsecase(repo, "test-secret", validate, log)

	_, err := uc.Login(context.Background(), dto.LoginRequest{
		Username: "dokter1",
		Password: "wrongpass",
	})
	if err == nil {
		t.Fatal("expected error for wrong credentials, got nil")
	}

	var businessErr *domain.BusinessError
	if !errors.As(err, &businessErr) {
		t.Fatalf("expected *domain.BusinessError, got %T", err)
	}
}

func TestLogin_Success(t *testing.T) {
	repo := &mockAuthRepository{
		verifikasiLoginFunc: func(ctx context.Context, username, password string) (*domain.User, error) {
			return &domain.User{IDUser: "DK001", NamaUser: "Dr. Budi"}, nil
		},
	}
	validate := validator.New()
	log := logger.New()
	uc := usecase.NewAuthUsecase(repo, "test-secret-key-12345", validate, log)

	resp, err := uc.Login(context.Background(), dto.LoginRequest{
		Username: "dokter1",
		Password: "pass123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Token == "" {
		t.Fatal("expected response with token")
	}
	if resp.KodeDokter != "DK001" {
		t.Errorf("expected KodeDokter=DK001, got %s", resp.KodeDokter)
	}
}

func TestLogin_DatabaseError(t *testing.T) {
	repo := &mockAuthRepository{
		verifikasiLoginFunc: func(ctx context.Context, username, password string) (*domain.User, error) {
			return nil, errors.New("connection refused")
		},
	}
	validate := validator.New()
	log := logger.New()
	uc := usecase.NewAuthUsecase(repo, "test-secret", validate, log)

	_, err := uc.Login(context.Background(), dto.LoginRequest{
		Username: "dokter1",
		Password: "pass123",
	})
	if err == nil {
		t.Fatal("expected error for database failure, got nil")
	}

	var businessErr *domain.BusinessError
	if errors.As(err, &businessErr) {
		t.Error("database error should NOT be a BusinessError")
	}
}
