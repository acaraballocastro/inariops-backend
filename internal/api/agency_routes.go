package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerAgencyRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {

	agencies := router.PathPrefix("/agencies").Subrouter()

	agencies.HandleFunc("", handlers.Agency.ListAgencies).
		Methods(http.MethodGet)

	agencies.HandleFunc("", handlers.Agency.CreateAgency).
		Methods(http.MethodPost)

	agencies.HandleFunc("/{id}", handlers.Agency.GetAgencyByID).
		Methods(http.MethodGet)

	agencies.HandleFunc("/{id}", handlers.Agency.UpdateAgency).
		Methods(http.MethodPatch)

	agencies.HandleFunc("/{id}", handlers.Agency.DeleteAgency).
		Methods(http.MethodDelete)
}
