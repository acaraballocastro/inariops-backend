package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerUserRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {

	users := router.PathPrefix("/users").Subrouter()

	users.HandleFunc("", handlers.User.GetAllUsers).
		Methods(http.MethodGet)

	users.HandleFunc("", handlers.User.CreateUser).
		Methods(http.MethodPost)

	users.HandleFunc("/guides", handlers.User.GetAllGuides).
		Methods(http.MethodGet)

	users.HandleFunc("/{id}", handlers.User.GetUserByID).
		Methods(http.MethodGet)

	users.HandleFunc("/{id}", handlers.User.UpdateUser).
		Methods(http.MethodPatch)

	users.HandleFunc("/{id}", handlers.User.DeactivateUser).
		Methods(http.MethodDelete)

	users.HandleFunc("/guides", handlers.User.GetAllGuides).
		Methods(http.MethodGet)
}
