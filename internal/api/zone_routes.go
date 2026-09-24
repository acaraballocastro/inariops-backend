package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerZoneRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	zones := router.PathPrefix("/zones").Subrouter()

	HandleAdminOrGuide(
		zones,
		"",
		handlers.Zone.ListZones,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		zones,
		"/{name}",
		handlers.Zone.GetZoneByName,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		zones,
		"",
		handlers.Zone.CreateZone,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		zones,
		"/{id}",
		handlers.Zone.UpdateZone,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		zones,
		"/{name}",
		handlers.Zone.DeleteZone,
		http.MethodDelete,
		cfg,
	)
}
