package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerCustomerRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	customers := router.PathPrefix("/customers").Subrouter()

	HandleAdminOrGuide(
		customers,
		"",
		handlers.Customer.GetAllCustomers,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		customers,
		"/{id}",
		handlers.Customer.GetCustomerByID,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		customers,
		"",
		handlers.Customer.CreateCustomer,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		customers,
		"/{id}",
		handlers.Customer.UpdateCustomer,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		customers,
		"/{id}",
		handlers.ReservationCustomerApplication.DeleteCustomer,
		http.MethodDelete,
		cfg,
	)
}
