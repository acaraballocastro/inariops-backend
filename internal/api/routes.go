package api

import (
	"database/sql"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/users"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	authRepo := auth.NewRepository(db)
	userRepo := users.NewRepository(db)

	authService := auth.NewService(authRepo)
	userService := users.NewService(userRepo, authService)

	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)

	router.HandleFunc("/auth", authHandler.HandleAuth)
	router.HandleFunc("/users", userHandler.HandleUsers)

	return router
}
