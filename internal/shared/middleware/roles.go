package middleware

import (
	"net/http"

	"inariops/internal/config"
	"inariops/internal/domain"
)

func RequireRole(
	cfg config.Config,
	roles ...domain.UserRole,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Development:
			// Authentication and authorization are disabled.
			if cfg.IsDevelopment() {
				next.ServeHTTP(w, r)
				return
			}

			claims, ok := GetClaims(r)

			if !ok {
				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			userRole := domain.UserRole(claims.Role)

			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(
				w,
				"forbidden",
				http.StatusForbidden,
			)
		})
	}
}

func RequireAdmin(
	cfg config.Config,
) func(http.Handler) http.Handler {

	return RequireRole(
		cfg,
		domain.RoleAdmin,
	)
}

func RequireGuide(
	cfg config.Config,
) func(http.Handler) http.Handler {

	return RequireRole(
		cfg,
		domain.RoleGuide,
	)
}

func RequireAdminOrGuide(
	cfg config.Config,
) func(http.Handler) http.Handler {

	return RequireRole(
		cfg,
		domain.RoleAdmin,
		domain.RoleGuide,
	)
}
