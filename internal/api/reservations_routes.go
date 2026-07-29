package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerReservationRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {

	reservations := router.PathPrefix("/reservations").Subrouter()

	reservations.HandleFunc("", handlers.Reservation.GetAllReservations).
		Methods(http.MethodGet)

	reservations.HandleFunc("", handlers.ReservationApplication.CreateReservation).
		Methods(http.MethodPost)

	reservations.HandleFunc("/{code}", handlers.Reservation.GetReservationByCode).
		Methods(http.MethodGet)

	reservations.HandleFunc("/{code}", handlers.ReservationApplication.UpdateReservation).
		Methods(http.MethodPatch)

	reservations.HandleFunc("/{code}", handlers.ReservationApplication.DeleteReservation).
		Methods(http.MethodDelete)

	reservations.HandleFunc("/{code}/customers", handlers.ReservationCustomerApplication.GetCustomersByReservationCode).
		Methods(http.MethodGet)
}
