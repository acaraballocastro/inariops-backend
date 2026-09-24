package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerItineraryRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	itinerary := router.PathPrefix(
		"/tour-days/{tour_day_id}/itinerary",
	).Subrouter()

	HandleAdminOrGuide(
		itinerary,
		"",
		handlers.ItineraryApplication.GetItineraryByTourDayID,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		itinerary,
		"",
		handlers.ItineraryApplication.CreateItineraryItem,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		itinerary,
		"/{item_id}",
		handlers.ItineraryApplication.UpdateItineraryItem,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		itinerary,
		"/{item_id}",
		handlers.ItineraryApplication.DeleteItineraryItem,
		http.MethodDelete,
		cfg,
	)
}
