package middleware

import (
	"context"
	"net/http"
	"strings"

	"inariops/internal/config"
	"inariops/internal/modules/auth"
)

type contextKey string

const (
	UserClaimsKey contextKey = "user_claims"
)

func AuthMiddleware(
	cfg config.Config,
	jwtManager *auth.JWTManager,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Development mode:
			// Authentication is disabled so the API can be used
			// without JWT tokens.
			if cfg.IsDevelopment() {
				next.ServeHTTP(w, r)
				return
			}

			// Production mode:
			// JWT is required.
			token, err := extractBearerToken(r)
			if err != nil {
				http.Error(
					w,
					"authorization token required",
					http.StatusUnauthorized,
				)
				return
			}

			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				http.Error(
					w,
					"invalid or expired authorization token",
					http.StatusUnauthorized,
				)
				return
			}

			// Store claims in the request context so handlers/services
			// can access the authenticated user later.
			ctx := context.WithValue(
				r.Context(),
				UserClaimsKey,
				claims,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

func extractBearerToken(r *http.Request) (string, error) {

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", auth.ErrInvalidToken
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 {
		return "", auth.ErrInvalidToken
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", auth.ErrInvalidToken
	}

	token := strings.TrimSpace(parts[1])

	if token == "" {
		return "", auth.ErrInvalidToken
	}

	return token, nil
}

func GetClaims(r *http.Request) (*auth.Claims, bool) {

	claims, ok := r.Context().Value(UserClaimsKey).(*auth.Claims)

	return claims, ok
}
