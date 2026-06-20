package api

import (
	"database/sql"
	"net/http"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/reservations"
	"inariops/internal/modules/users"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/response"

	"github.com/gorilla/mux"
)

// NewRouter initializes and configures all API routes with proper organization,
// HTTP methods, and API versioning.
func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize repositories
	authRepo := auth.NewRepository(db)
	userRepo := users.NewRepository(db)
	reservationRepo := reservations.NewRepository(db)

	// Initialize services
	authService := auth.NewService(authRepo)
	userService := users.NewService(userRepo, authService)
	reservationService := reservations.NewService(reservationRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)
	reservationHandler := reservations.NewHandler(reservationService)

	router.HandleFunc("/health", healthCheck).Methods(http.MethodGet)

	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	authRoutes := apiV1.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("", authHandler.HandleAuth).Methods(http.MethodPost, http.MethodPatch)

	// Users routes
	usersRoutes := apiV1.PathPrefix("/users").Subrouter()
	usersRoutes.HandleFunc("", userHandler.HandleUsers).Methods(http.MethodGet, http.MethodPost)
	usersRoutes.HandleFunc("/{id}", userHandler.HandleUsers).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)

	// Reservations routes
	reservationsRoutes := apiV1.PathPrefix("/reservations").Subrouter()
	reservationsRoutes.HandleFunc("", reservationHandler.HandleReservations).Methods(http.MethodGet, http.MethodPost)
	reservationsRoutes.HandleFunc("/{id}", reservationHandler.HandleReservations).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)

	// 404 handler for undefined routes
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// Method not allowed handler
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	return router
}

// healthCheck returns a simple health status response
func healthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	logger.Info("Health check requested")
}

// notFoundHandler handles requests to undefined routes
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotFound, "endpoint not found")
	logger.Error("404 - Not Found: %s %s", r.Method, r.RequestURI)
}

// methodNotAllowedHandler handles requests with unsupported HTTP methods
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	logger.Error("405 - Method Not Allowed: %s %s", r.Method, r.RequestURI)
}
