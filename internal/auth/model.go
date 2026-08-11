package auth

type User struct {
	IDUser   string `json:"id_user"`
	NamaUser string `json:"nama_user"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}
