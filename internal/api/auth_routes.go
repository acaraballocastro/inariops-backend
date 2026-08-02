package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerAuthRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {
	auth := router.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/", handlers.Auth.Login).
		Methods(http.MethodPost)

	auth.HandleFunc("/", handlers.Auth.ChangePassword).
		Methods(http.MethodPatch)

	//TODO: Implement logout route
}
