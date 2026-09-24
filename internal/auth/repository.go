package auth

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	VerifikasiLogin(ctx context.Context, username string, password string) (*User, error)
}

type repository struct {
	db          *sql.DB
	userKey     string
	passwordKey string
}

func NewRepository(db *sql.DB, userKey string, passwordKey string) Repository {
	return &repository{
		db:          db,
		userKey:     userKey,
		passwordKey: passwordKey,
	}
}

func (r *repository) VerifikasiLogin(ctx context.Context, username string, password string) (*User, error) {
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

	var user User
	err := r.db.QueryRowContext(ctx, query,
		r.userKey,
		r.userKey,
		username, r.userKey,
		password, r.passwordKey,
	).Scan(
		&user.IDUser,
		&user.NamaUser,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
