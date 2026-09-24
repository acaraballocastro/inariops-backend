package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerPlacesRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {
	places := router.PathPrefix("/tour-days/{tour_day_id}/places").Subrouter()

	places.HandleFunc(
		"",
		handlers.PlaceApplication.GetPlacesByTourDayID,
	).Methods(http.MethodGet)

	places.HandleFunc(
		"",
		handlers.PlaceApplication.CreatePlace,
	).Methods(http.MethodPost)

	places.HandleFunc(
		"/{place_id}",
		handlers.PlaceApplication.UpdatePlace,
	).Methods(http.MethodPatch)

	places.HandleFunc(
		"/{place_id}",
		handlers.PlaceApplication.DeletePlace,
	).Methods(http.MethodDelete)
}
