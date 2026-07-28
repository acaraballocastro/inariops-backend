package api

import (
	"database/sql"
	"net/http"

	reservationcustomersapp "inariops/internal/application/reservation_customers"
	reservationsapp "inariops/internal/application/reservations"
	tourdaysapp "inariops/internal/application/tour_days"
	"inariops/internal/modules/auth"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/guides"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/modules/users"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/response"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Mux middleware
	router.Use(mux.CORSMethodMiddleware(router))

	// Repositories
	authRepo := auth.NewRepository(db)
	guidesRepo := guides.NewRepository(db)
	userRepo := users.NewRepository(db)
	reservationRepo := reservations.NewRepository(db)
	tourDayRepo := tourdays.NewRepository(db)
	customerRepo := customers.NewRepository(db)
	reservatationscustomersRepo := reservationscustomers.NewRepository(db)
	assigmentRepo := tourdaysapp.NewRepository(db)

	// Services
	authService := auth.NewService(authRepo)
	guidesService := guides.NewService(guidesRepo)

	//TODO: Check if we need to pass guidesService to userService, MAYBE IMPLEMENTATION OF A APPLICATION SERVICE FOR USERS
	userService := users.NewService(userRepo, authService, guidesService)

	reservationService := reservations.NewService(reservationRepo)
	tourDayService := tourdays.NewService(tourDayRepo)
	customerService := customers.NewService(customerRepo)

	// Application Services
	reservationUseCase := reservationsapp.NewService(reservationRepo, customerRepo, reservatationscustomersRepo, tourDayRepo)
	tourdayUseCase := tourdaysapp.NewService(tourDayRepo, guidesRepo, assigmentRepo)
	reservationCustomerUseCase := reservationcustomersapp.NewService(reservationRepo, customerRepo, reservatationscustomersRepo)

	// Handlers
	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)
	reservationHandler := reservations.NewHandler(reservationService)
	tourDayHandler := tourdays.NewHandler(tourDayService)
	customerHandler := customers.NewHandler(customerService)
	guideHandler := guides.NewHandler(guidesService)

	// Application Handlers
	reservationAppHandler := reservationsapp.NewHandler(reservationUseCase)
	tourDayAppHandler := tourdaysapp.NewHandler(tourdayUseCase)
	reservationCustomerAppHandler := reservationcustomersapp.NewHandler(reservationCustomerUseCase)

	// Health
	router.HandleFunc("/health", healthCheck).Methods(http.MethodGet)

	// Global OPTIONS handler (CORS preflight)
	router.PathPrefix("/").Methods(http.MethodOptions).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

	// API v1
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// =====================
	// AUTH
	// =====================

	apiV1.HandleFunc("/auth", authHandler.Login).
		Methods(http.MethodPost)

	apiV1.HandleFunc("/auth", authHandler.ChangePassword).
		Methods(http.MethodPatch)

	// =====================
	// USERS
	// =====================

	apiV1.HandleFunc("/users", userHandler.GetAllUsers).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/users/guides", userHandler.GetAllGuides).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/users", userHandler.CreateUser).
		Methods(http.MethodPost)

	apiV1.HandleFunc("/users/{id}", userHandler.GetUserByID).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/users/{id}", userHandler.UpdateUser).
		Methods(http.MethodPatch)

	apiV1.HandleFunc("/users/{id}", userHandler.DeactivateUser).
		Methods(http.MethodDelete)

	// =====================
	// RESERVATIONS
	// =====================

	apiV1.HandleFunc("/reservations", reservationHandler.GetAllReservations).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/reservations/{code}", reservationHandler.GetReservationByCode).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/reservations/{code}/customers", reservationCustomerAppHandler.GetCustomersByReservationCode).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/reservations", reservationAppHandler.CreateReservation).
		Methods(http.MethodPost)

	apiV1.HandleFunc("/reservations/{code}", reservationAppHandler.UpdateReservation).
		Methods(http.MethodPatch)

	apiV1.HandleFunc("/reservations/{code}", reservationAppHandler.DeleteReservation).
		Methods(http.MethodDelete)

		// =====================
		// TOUR DAYS
		// =====================

	// =====================
	// TOUR DAYS
	// =====================

	// Specific routes first
	apiV1.HandleFunc("/tour-days/assign-guide", tourDayAppHandler.AssignGuide).
		Methods(http.MethodPatch)

	apiV1.HandleFunc("/tour-days/unassign-guide", tourDayAppHandler.UnassignGuide).
		Methods(http.MethodPatch)

	apiV1.HandleFunc(
		"/tour-days/by-reservation/{reservation_id}",
		tourDayHandler.GetTourDaysByReservationID,
	).Methods(http.MethodGet)

	// Generic routes afterwards
	apiV1.HandleFunc("/tour-days", tourDayHandler.CreateTourDay).
		Methods(http.MethodPost)

	apiV1.HandleFunc("/tour-days/{id}", tourDayHandler.GetTourDayByID).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/tour-days/{id}", tourDayHandler.UpdateTourDay).
		Methods(http.MethodPatch)

	apiV1.HandleFunc("/tour-days/{id}", tourDayHandler.CancelTourDay).
		Methods(http.MethodDelete)

	// =====================
	// GUIDES
	// =====================
	apiV1.HandleFunc("/guides", userHandler.GetAllGuides).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/guides/{id}", guideHandler.GetGuideByID).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/guides/user/{user_id}", guideHandler.GetGuideByUserID).
		Methods(http.MethodGet)

	// ======================
	// Customers
	// ======================
	apiV1.HandleFunc("/customers", customerHandler.GetAllCustomers).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/customers/{id}", customerHandler.GetCustomerByID).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/customers", customerHandler.CreateCustomer).
		Methods(http.MethodPost)

	apiV1.HandleFunc("/customers/{id}", customerHandler.UpdateCustomer).
		Methods(http.MethodPatch)

	apiV1.HandleFunc("/customers/{id}", customerHandler.DeleteCustomer).
		Methods(http.MethodDelete)

	// Error handlers
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	return router
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})

	logger.Info("Health check requested")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotFound, "endpoint not found")
	logger.Error("404 - Not Found: %s %s", r.Method, r.RequestURI)
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	logger.Error("405 - Method Not Allowed: %s %s", r.Method, r.RequestURI)
}
