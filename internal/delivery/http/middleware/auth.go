package middleware

import (
	"context"
	"net/http"
	"strings"

	"tracking-backend/internal/delivery/http/response"
	"tracking-backend/internal/domain"
	"tracking-backend/internal/usecase"
)

type contextKey string

const userContextKey contextKey = "user"

// Auth creates an authentication middleware that validates the Bearer token
// and stores the JWT claims in the request context.
func Auth(authUC usecase.AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			claims, err := authUC.VerifyToken(parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext extracts JWT claims from the request context.
func UserFromContext(ctx context.Context) (*domain.JWTClaims, bool) {
	claims, ok := ctx.Value(userContextKey).(*domain.JWTClaims)
	return claims, ok
}

// UserIDFromContext extracts the user ID from the request context.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	claims, ok := UserFromContext(ctx)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}
