package api

import (
	"database/sql"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/users"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	userRepo := users.NewRepository(db)

	userService := users.NewService(userRepo)

	userHandler := users.NewHandler(userService)

	authRepo := auth.NewRepository(db)

	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/auth/change-password", authHandler.ChangePassword).Methods("POST")

	router.HandleFunc("/users", userHandler.HandleUsers)

	return router
}
