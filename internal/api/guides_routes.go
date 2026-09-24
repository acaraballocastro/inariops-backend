package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerGuideRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {
	guides := router.PathPrefix("/guides").Subrouter()

	guides.HandleFunc("/{id}", handlers.Guide.GetGuideByID).
		Methods(http.MethodGet)

	guides.HandleFunc("/user/{user_id}", handlers.Guide.GetGuideByUserID).
		Methods(http.MethodGet)

	guides.HandleFunc("", handlers.GuideApplication.GetAllGuidesDetail).
		Methods(http.MethodGet)

	guides.HandleFunc("", handlers.GuideApplication.CreateGuide).
		Methods(http.MethodPost)

	guides.HandleFunc("/{id}", handlers.GuideApplication.UpdateGuide).
		Methods(http.MethodPatch)

	guides.HandleFunc("/{id}/languages", handlers.GuideApplication.GetLanguagesByGuideID).
		Methods(http.MethodGet)

	guides.HandleFunc("/{id}/languages", handlers.GuideApplication.AddLanguageToGuide).
		Methods(http.MethodPost)

	guides.HandleFunc("/{id}/languages/{language_id}", handlers.GuideApplication.RemoveLanguageFromGuide).
		Methods(http.MethodDelete)

	// Availabilities routes
	guides.HandleFunc("/{id}/availabilities", handlers.GuideApplication.GetAvailabilitiesByGuideID).
		Methods(http.MethodGet)

	guides.HandleFunc("/{id}/availabilities", handlers.GuideApplication.CreateAvailability).
		Methods(http.MethodPost)

	guides.HandleFunc("/{id}/availabilities/{availability_id}", handlers.GuideApplication.UpdateAvailability).
		Methods(http.MethodPatch)

	guides.HandleFunc("/{id}/availabilities/{availability_id}", handlers.GuideApplication.DeleteAvailability).
		Methods(http.MethodDelete)
}
