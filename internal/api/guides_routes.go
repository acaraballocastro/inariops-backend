package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerGuideRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	guides := router.PathPrefix("/guides").Subrouter()

	// =====================
	// ADMIN + GUIDE
	// =====================

	HandleAdminOrGuide(
		guides,
		"/{id}",
		handlers.Guide.GetGuideByID,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		guides,
		"/user/{user_id}",
		handlers.Guide.GetGuideByUserID,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		guides,
		"/{id}",
		handlers.GuideApplication.UpdateGuide,
		http.MethodPatch,
		cfg,
	)

	// Languages

	HandleAdminOrGuide(
		guides,
		"/{id}/languages",
		handlers.GuideApplication.GetLanguagesByGuideID,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		guides,
		"/{id}/languages",
		handlers.GuideApplication.AddLanguageToGuide,
		http.MethodPost,
		cfg,
	)

	HandleAdminOrGuide(
		guides,
		"/{id}/languages/{language_id}",
		handlers.GuideApplication.RemoveLanguageFromGuide,
		http.MethodDelete,
		cfg,
	)

	// Availabilities

	HandleAdminOrGuide(
		guides,
		"/{id}/availabilities",
		handlers.GuideApplication.GetAvailabilitiesByGuideID,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		guides,
		"/{id}/availabilities",
		handlers.GuideApplication.CreateAvailability,
		http.MethodPost,
		cfg,
	)

	HandleAdminOrGuide(
		guides,
		"/{id}/availabilities/{availability_id}",
		handlers.GuideApplication.UpdateAvailability,
		http.MethodPatch,
		cfg,
	)

	HandleAdminOrGuide(
		guides,
		"/{id}/availabilities/{availability_id}",
		handlers.GuideApplication.DeleteAvailability,
		http.MethodDelete,
		cfg,
	)

	// =====================
	// ADMIN ONLY
	// =====================

	HandleAdmin(
		guides,
		"",
		handlers.GuideApplication.GetAllGuidesDetail,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		guides,
		"",
		handlers.GuideApplication.CreateGuide,
		http.MethodPost,
		cfg,
	)
}
