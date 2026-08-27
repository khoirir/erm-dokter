package auth

import (
	"encoding/json"
	"net/http"

	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared/apperror"
)

type Handler struct {
	authService Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		authService: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, loginRateLimit func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("POST /api/v1/auth/login", loginRateLimit(h.Login))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.NewBusinessError("Format request JSON tidak valid"))
		return
	}

	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	resp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Login berhasil", resp)
}
