package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerCustomerRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {

	customers := router.PathPrefix("/customers").Subrouter()

	customers.HandleFunc("", handlers.Customer.GetAllCustomers).
		Methods(http.MethodGet)

	customers.HandleFunc("", handlers.Customer.CreateCustomer).
		Methods(http.MethodPost)

	customers.HandleFunc("/{id}", handlers.Customer.GetCustomerByID).
		Methods(http.MethodGet)

	customers.HandleFunc("/{id}", handlers.Customer.UpdateCustomer).
		Methods(http.MethodPatch)

	customers.HandleFunc("/{id}", handlers.Customer.DeleteCustomer).
		Methods(http.MethodDelete)
}
