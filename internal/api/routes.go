package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/response"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(
	handlers *bootstrap.Handlers,
) *mux.Router {

	router := mux.NewRouter()

	router.Use(mux.CORSMethodMiddleware(router))

	// Health
	router.HandleFunc("/health", healthCheck).
		Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()

	registerAuthRoutes(api, handlers)
	registerUserRoutes(api, handlers)
	registerReservationRoutes(api, handlers)
	registerTourDayRoutes(api, handlers)
	registerGuideRoutes(api, handlers)
	registerCustomerRoutes(api, handlers)
	registerAgencyRoutes(api, handlers)
	registerLanguageRoutes(api, handlers)

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
