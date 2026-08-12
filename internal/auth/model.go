package auth

import (
	"erm-dokter/internal/shared/apperror"
	"strings"
)

type User struct {
	IDUser   string `json:"id_user"`
	NamaUser string `json:"nama_user"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}

func (r *LoginRequest) Sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)
}

func (r *LoginRequest) Validate() apperror.ValidationError {
	r.Sanitize()
	errs := make(apperror.ValidationError)
	if r.Username == "" {
		errs["username"] = "username wajib diisi"
	}
	if r.Password == "" {
		errs["password"] = "password wajib diisi"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

