package api

import (
	"database/sql"

	"inariops/internal/modules/users"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	userRepo := users.NewRepository(db)

	userService := users.NewService(userRepo)

	userHandler := users.NewHandler(userService)

	router.HandleFunc("/users", userHandler.HandleUsers)

	return router
}
