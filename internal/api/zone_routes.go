package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerZoneRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {

	zones := router.PathPrefix("/zones").Subrouter()

	zones.HandleFunc("", handlers.Zone.ListZones).
		Methods(http.MethodGet)

	zones.HandleFunc("", handlers.Zone.CreateZone).
		Methods(http.MethodPost)

	zones.HandleFunc("/{name}", handlers.Zone.GetZoneByName).
		Methods(http.MethodGet)

	zones.HandleFunc("/{id}", handlers.Zone.UpdateZone).
		Methods(http.MethodPatch)

	zones.HandleFunc("/{name}", handlers.Zone.DeleteZone).
		Methods(http.MethodDelete)
}
