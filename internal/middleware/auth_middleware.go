package middleware

import (
	"context"
	"net/http"
	"strings"

	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/pkg/token"
	"erm-dokter/internal/shared/apperror"
)

type contextKey string

const UserClaimKey contextKey = "userClaim"

const (
	RoleDokter  = "dokter"
	RoleService = "service"
)

func JWTMiddleware(jwtSecret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Kredensial login tidak valid", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Error(w, http.StatusUnauthorized, "Kredensial login tidak valid", nil)
				return
			}

			claims, err := token.ValidateToken(parts[1], jwtSecret)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Kredensial login tidak valid atau telah kedaluwarsa", nil)
				return
			}

			if claims.Role == "" {
				claims.Role = RoleDokter
			}

			if claims.KodeDokter != "" {
				w.Header().Set("X-User-ID", claims.KodeDokter)
			}

			ctx := context.WithValue(r.Context(), UserClaimKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func ServiceOrJWTMiddleware(jwtSecret, serviceAPIKey string) func(http.HandlerFunc) http.HandlerFunc {
	cleanServiceKey := strings.TrimSpace(serviceAPIKey)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))
			if apiKey != "" {
				if cleanServiceKey == "" || apiKey != cleanServiceKey {
					response.Error(w, http.StatusUnauthorized, "API Key tidak valid", nil)
					return
				}

				serviceClaims := &token.Claims{
					KodeDokter: "",
					NamaUser:   "External Service",
					Role:       RoleService,
				}
				w.Header().Set("X-User-ID", "service")
				ctx := context.WithValue(r.Context(), UserClaimKey, serviceClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			JWTMiddleware(jwtSecret)(next).ServeHTTP(w, r)
		}
	}
}

func AuthMiddleware(jwtSecret, serviceAPIKey string) func(http.HandlerFunc) http.HandlerFunc {
	return ServiceOrJWTMiddleware(jwtSecret, serviceAPIKey)
}

func GetKodeDokter(ctx context.Context, allowService ...bool) (string, error) {
	claims, ok := ctx.Value(UserClaimKey).(*token.Claims)
	if !ok || claims == nil {
		return "", apperror.NewUnauthorizedError("Kredensial login tidak ditemukan")
	}

	canService := len(allowService) > 0 && allowService[0]
	if claims.Role == RoleService {
		if canService {
			return "", nil
		}
		return "", apperror.NewUnauthorizedError("Identitas dokter pada akun ini tidak valid")
	}

	if strings.TrimSpace(claims.KodeDokter) == "" {
		return "", apperror.NewUnauthorizedError("Identitas dokter pada akun ini tidak valid")
	}

	return claims.KodeDokter, nil
}

func IsService(ctx context.Context) bool {
	claims, ok := ctx.Value(UserClaimKey).(*token.Claims)
	return ok && claims != nil && claims.Role == RoleService
}
