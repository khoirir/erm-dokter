package middleware

import (
	"context"
	"net/http"
	"strings"

	"erm-dokter/internal/pkg/response"
	"erm-dokter/internal/pkg/token"
)

type contextKey string

const UserClaimKey contextKey = "userClaim"

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
				response.Error(w, http.StatusUnauthorized, "Format kredensial login tidak valid", nil)
				return
			}

			claims, err := token.ValidateToken(parts[1], jwtSecret)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Kredensial login tidak valid atau telah kedaluwarsa", nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func GetUserClaim(ctx context.Context) *token.Claims {
	if claims, ok := ctx.Value(UserClaimKey).(*token.Claims); ok {
		return claims
	}
	return nil
}
