package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerAgencyRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	agencies := router.PathPrefix("/agencies").Subrouter()

	HandleAdminOrGuide(
		agencies,
		"",
		handlers.Agency.ListAgencies,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		agencies,
		"/{id}",
		handlers.Agency.GetAgencyByID,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		agencies,
		"",
		handlers.Agency.CreateAgency,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		agencies,
		"/{id}",
		handlers.Agency.UpdateAgency,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		agencies,
		"/{id}",
		handlers.Agency.DeleteAgency,
		http.MethodDelete,
		cfg,
	)
}
