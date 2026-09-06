package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerActivityRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {
	activities := router.PathPrefix("/activities").Subrouter()

	activities.HandleFunc("", handlers.Activity.GetAllActivities).
		Methods(http.MethodGet)

	activities.HandleFunc("/{id}", handlers.Activity.GetActivityByID).
		Methods(http.MethodGet)

	activities.HandleFunc("", handlers.Activity.CreateActivity).
		Methods(http.MethodPost)

	activities.HandleFunc("/{id}", handlers.Activity.UpdateActivity).
		Methods(http.MethodPatch)

	activities.HandleFunc("/{id}", handlers.Activity.DeleteActivity).
		Methods(http.MethodDelete)
}
