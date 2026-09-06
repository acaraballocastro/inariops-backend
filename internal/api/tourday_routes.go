package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerTourDayRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {
	tourDays := router.PathPrefix("/tour-days").Subrouter()

	// Specific tour day routes
	tourDays.HandleFunc("/assign-guide", handlers.TourDayApplication.AssignGuide).
		Methods(http.MethodPatch)

	tourDays.HandleFunc("/unassign-guide", handlers.TourDayApplication.UnassignGuide).
		Methods(http.MethodPatch)

	tourDays.HandleFunc("/available-for-guide/{guide_id}", handlers.TourDayApplication.GetTourDaysAvailableForGuide).
		Methods(http.MethodGet)

	tourDays.HandleFunc("/by-reservation/{reservation_id}", handlers.TourDay.GetTourDaysByReservationID).
		Methods(http.MethodGet)

	// CRUD operations for tour days
	tourDays.HandleFunc("", handlers.TourDay.CreateTourDay).
		Methods(http.MethodPost)

	tourDays.HandleFunc("/{id}", handlers.TourDay.GetTourDayByID).
		Methods(http.MethodGet)

	tourDays.HandleFunc("/{id}", handlers.TourDay.UpdateTourDay).
		Methods(http.MethodPatch)

	tourDays.HandleFunc("/{id}/confirm", handlers.TourDayApplication.ConfirmTourDay).
		Methods(http.MethodPatch)

	tourDays.HandleFunc("/{id}", handlers.TourDay.CancelTourDay).
		Methods(http.MethodDelete)
}
