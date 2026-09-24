package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"inariops/internal/modules/auth"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/middleware"
	"inariops/internal/shared/response"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(
	handlers *bootstrap.Handlers,
	cfg config.Config,
	jwtManager *auth.JWTManager,
) *mux.Router {

	router := mux.NewRouter()

	router.Use(mux.CORSMethodMiddleware(router))

	router.PathPrefix("/api/v1").
		Methods(http.MethodOptions).
		HandlerFunc(apiPreflightHandler)

	// =====================
	// Health
	// =====================

	router.HandleFunc("/health", healthCheck).
		Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()

	// =====================
	// Public
	// =====================

	registerAuthRoutes(api, handlers)

	// =====================
	// Protected
	// =====================

	protected := api.PathPrefix("").Subrouter()

	protected.Use(
		middleware.AuthMiddleware(cfg, jwtManager),
	)

	registerUserRoutes(
		protected,
		handlers,
		cfg,
	)

	registerReservationRoutes(
		protected,
		handlers,
		cfg,
	)

	registerTourDayRoutes(
		protected,
		handlers,
		cfg,
	)

	registerActivityRoutes(
		protected,
		handlers,
		cfg,
	)

	registerItineraryRoutes(
		protected,
		handlers,
		cfg,
	)

	registerPlacesRoutes(
		protected,
		handlers,
		cfg,
	)

	registerGuideRoutes(
		protected,
		handlers,
		cfg,
	)

	registerCustomerRoutes(
		protected,
		handlers,
		cfg,
	)

	registerAgencyRoutes(
		protected,
		handlers,
		cfg,
	)

	registerLanguageRoutes(
		protected,
		handlers,
		cfg,
	)

	registerZoneRoutes(
		protected,
		handlers,
		cfg,
	)

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
	response.Error(
		w,
		http.StatusNotFound,
		"endpoint not found",
	)

	logger.Error(
		"404 - Not Found: %s %s",
		r.Method,
		r.RequestURI,
	)
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(
		w,
		http.StatusMethodNotAllowed,
		"method not allowed",
	)

	logger.Error(
		"405 - Method Not Allowed: %s %s",
		r.Method,
		r.RequestURI,
	)
}

func apiPreflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Access-Control-Allow-Origin",
		"*",
	)

	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Content-Type, Authorization",
	)

	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, POST, PATCH, DELETE, OPTIONS",
	)

	w.WriteHeader(http.StatusNoContent)
}
