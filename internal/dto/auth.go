package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	KodeDokter string `json:"kode_dokter"`
	NamaDokter string `json:"nama_dokter"`
}
