package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerTourDayRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	tourDays := router.PathPrefix("/tour-days").Subrouter()

	// =====================
	// ADMIN ONLY
	// =====================

	HandleAdmin(
		tourDays,
		"/assign-guide",
		handlers.TourDayApplication.AssignGuide,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		tourDays,
		"/unassign-guide",
		handlers.TourDayApplication.UnassignGuide,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		tourDays,
		"",
		handlers.TourDay.CreateTourDay,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		tourDays,
		"/{id}",
		handlers.TourDay.UpdateTourDay,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		tourDays,
		"/{id}/confirm",
		handlers.TourDayApplication.ConfirmTourDay,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		tourDays,
		"/{id}",
		handlers.TourDay.CancelTourDay,
		http.MethodDelete,
		cfg,
	)

	// =====================
	// ADMIN + GUIDE
	// =====================

	HandleAdminOrGuide(
		tourDays,
		"/available-for-guide/{guide_id}",
		handlers.TourDayApplication.GetTourDaysAvailableForGuide,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		tourDays,
		"/by-reservation/{reservation_id}",
		handlers.TourDay.GetTourDaysByReservationID,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		tourDays,
		"/{id}",
		handlers.TourDay.GetTourDayByID,
		http.MethodGet,
		cfg,
	)
}
