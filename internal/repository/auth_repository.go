package repository

import (
	"context"
	"database/sql"
	"fmt"

	"erm-dokter/internal/domain"
)

type authRepository struct {
	db          *sql.DB
	userKey     string
	passwordKey string
}

func NewAuthRepository(db *sql.DB, userKey, passwordKey string) domain.AuthRepository {
	return &authRepository{
		db:          db,
		userKey:     userKey,
		passwordKey: passwordKey,
	}
}

func (r *authRepository) VerifikasiLogin(ctx context.Context, username string, password string) (*domain.User, error) {
	query := `
		SELECT 
			AES_DECRYPT(u.id_user, ?) AS kode_dokter,
			d.nm_dokter AS nama_user
		FROM user u
		INNER JOIN dokter d ON AES_DECRYPT(u.id_user, ?) = d.kd_dokter
		WHERE u.id_user = AES_ENCRYPT(?, ?)
		  AND u.password = AES_ENCRYPT(?, ?)
		LIMIT 1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query,
		r.userKey,
		r.userKey,
		username, r.userKey,
		password, r.passwordKey,
	).Scan(
		&user.IDUser,
		&user.NamaUser,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("gagal memverifikasi login: %w", err)
	}

	return &user, nil
}
