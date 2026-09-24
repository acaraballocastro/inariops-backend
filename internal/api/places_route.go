package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerPlacesRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	places := router.PathPrefix(
		"/tour-days/{tour_day_id}/places",
	).Subrouter()

	HandleAdminOrGuide(
		places,
		"",
		handlers.PlaceApplication.GetPlacesByTourDayID,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		places,
		"",
		handlers.PlaceApplication.CreatePlace,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		places,
		"/{place_id}",
		handlers.PlaceApplication.UpdatePlace,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		places,
		"/{place_id}",
		handlers.PlaceApplication.DeletePlace,
		http.MethodDelete,
		cfg,
	)
}
