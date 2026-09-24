package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerItineraryRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {

	itinerary := router.PathPrefix("/tour-days/{tour_day_id}/itinerary").Subrouter()

	itinerary.HandleFunc(
		"",
		handlers.ItineraryApplication.GetItineraryByTourDayID,
	).Methods(http.MethodGet)

	itinerary.HandleFunc(
		"",
		handlers.ItineraryApplication.CreateItineraryItem,
	).Methods(http.MethodPost)

	itinerary.HandleFunc(
		"/{item_id}",
		handlers.ItineraryApplication.UpdateItineraryItem,
	).Methods(http.MethodPatch)

	itinerary.HandleFunc(
		"/{item_id}",
		handlers.ItineraryApplication.DeleteItineraryItem,
	).Methods(http.MethodDelete)
}
