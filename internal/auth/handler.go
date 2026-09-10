package auth

import (
	"net/http"

	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/shared"
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

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	loginRateLimit func(http.HandlerFunc) http.HandlerFunc,
	authMiddleware func(http.HandlerFunc) http.HandlerFunc,
	timeoutMiddleware func(http.HandlerFunc) http.HandlerFunc,
) {
	mux.HandleFunc("POST /api/v1/auth/login", loginRateLimit(timeoutMiddleware(h.Login)))
	mux.HandleFunc("POST /api/v1/auth/logout", authMiddleware(timeoutMiddleware(h.Logout)))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !shared.DecodeJSON(w, r, &req, "Data login tidak valid") {
		return
	}

	req.Sanitize()
	if errs := req.Validate(); errs != nil {
		apperror.HandleError(w, errs)
		return
	}

	w.Header().Set("X-User-ID", req.Username)

	resp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response.Success(w, "Login berhasil", resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	response.Success(w, "Logout berhasil", nil)
}

