package middleware

import (
	"net/http"

	"inariops/internal/config"
	"inariops/internal/domain"

	"github.com/gorilla/mux"
)

type UserOwnershipChecker interface {
	UserExists(id string) (bool, error)
}

type GuideOwnershipChecker interface {
	GetGuideByID(id string) (domain.Guide, error)
}

func RequireAdminOrOwnUser(
	cfg config.Config,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if cfg.IsDevelopment() {
				next.ServeHTTP(w, r)
				return
			}

			claims, ok := GetClaims(r)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if domain.UserRole(claims.Role) == domain.RoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			if domain.UserRole(claims.Role) != domain.RoleGuide {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			userID := mux.Vars(r)["id"]

			if claims.UserID != userID {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func routeParam(r *http.Request, name string) string {
	// Implementado mediante mux.Vars en el siguiente paso.
	return ""
}
