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

func (r *authRepository) CariByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := fmt.Sprintf(`
		SELECT 
			AES_DECRYPT(u.id_user, '%[1]s') AS kode_dokter,
			COALESCE(d.nm_dokter, AES_DECRYPT(u.id_user, '%[1]s')) AS nama_user,
			AES_DECRYPT(u.password, '%[2]s') AS password
		FROM user u
		LEFT JOIN dokter d ON AES_DECRYPT(u.id_user, '%[1]s') = d.kd_dokter
		WHERE u.id_user = AES_ENCRYPT(?, '%[1]s')
		LIMIT 1
	`, r.userKey, r.passwordKey)

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.KodeDokter,
		&user.NamaUser,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("gagal mencari user di database: %w", err)
	}

	return &user, nil
}
