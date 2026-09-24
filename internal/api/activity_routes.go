package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerActivityRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	activities := router.PathPrefix("/activities").Subrouter()

	HandleAdminOrGuide(
		activities,
		"",
		handlers.Activity.GetAllActivities,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		activities,
		"/{id}",
		handlers.Activity.GetActivityByID,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		activities,
		"",
		handlers.Activity.CreateActivity,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		activities,
		"/{id}",
		handlers.Activity.UpdateActivity,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		activities,
		"/{id}",
		handlers.Activity.DeleteActivity,
		http.MethodDelete,
		cfg,
	)
}
