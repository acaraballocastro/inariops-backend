package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerReservationRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	reservations := router.PathPrefix("/reservations").Subrouter()

	// =====================
	// ADMIN + GUIDE
	// =====================

	HandleAdminOrGuide(
		reservations,
		"",
		handlers.Reservation.GetAllReservations,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		reservations,
		"/{code}",
		handlers.Reservation.GetReservationByCode,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		reservations,
		"/details/{code}",
		handlers.ReservationApplication.GetReservationDetailByCode,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		reservations,
		"/{code}/customers",
		handlers.ReservationCustomerApplication.GetCustomersByReservationCode,
		http.MethodGet,
		cfg,
	)

	// =====================
	// ADMIN ONLY
	// =====================

	HandleAdmin(
		reservations,
		"",
		handlers.ReservationApplication.CreateReservation,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		reservations,
		"/{code}",
		handlers.ReservationApplication.UpdateReservation,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		reservations,
		"/{code}/signature",
		handlers.ReservationApplication.UpdateSignatureStatus,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		reservations,
		"/{code}/voucher",
		handlers.ReservationApplication.UpdateVoucherStatus,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		reservations,
		"/{code}",
		handlers.ReservationApplication.DeleteReservation,
		http.MethodDelete,
		cfg,
	)

	HandleAdmin(
		reservations,
		"/erase/{code}",
		handlers.ReservationApplication.EraseReservation,
		http.MethodDelete,
		cfg,
	)
}
