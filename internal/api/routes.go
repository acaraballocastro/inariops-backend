package api

import (
	"database/sql"
	"net/http"

	"inariops/internal/modules/users"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	mux.HandleFunc("/users", userHandler.HandleUsers)

	return mux
}
