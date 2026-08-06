package domain

import (
	"context"
)

type Penjamin struct {
    KodePenjamin string `json:"kode_penjamin"`
    Nama         string `json:"nama"`
}

type PenjaminRepository interface {
	DaftarPenjamin(ctx context.Context) ([]Penjamin, error)
}
type PenjaminUsecase interface {
	DaftarPenjamin(ctx context.Context) ([]Penjamin, error)
}