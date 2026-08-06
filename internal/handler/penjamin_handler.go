package handler

import (
	"erm-dokter/internal/domain"
	"erm-dokter/pkg/response"
	"net/http"
)

type PenjaminHandler struct {
	penjaminUsecase domain.PenjaminUsecase
}

func NewPenjaminHandler(usecase domain.PenjaminUsecase) *PenjaminHandler {
	return &PenjaminHandler{
		penjaminUsecase: usecase,
	}
}

func (h *PenjaminHandler) DaftarPenjamin(w http.ResponseWriter, r *http.Request) {
	daftarPenjamin, err := h.penjaminUsecase.DaftarPenjamin(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, "Berhasil mengambil daftar penjamin", daftarPenjamin)
}