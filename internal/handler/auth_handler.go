package handler

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/domain"
	"erm-dokter/pkg/response"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(usecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: usecase,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method tidak diizinkan", nil)
		return
	}

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Format request JSON tidak valid", nil)
		return
	}

	resp, err := h.authUsecase.Login(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(w, "Login berhasil", resp)
}
