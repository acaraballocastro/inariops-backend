package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerLanguageRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	languages := router.PathPrefix("/languages").Subrouter()

	HandleAdminOrGuide(
		languages,
		"",
		handlers.Language.ListLanguages,
		http.MethodGet,
		cfg,
	)

	HandleAdminOrGuide(
		languages,
		"/{code}",
		handlers.Language.GetLanguageByCode,
		http.MethodGet,
		cfg,
	)

	HandleAdmin(
		languages,
		"",
		handlers.Language.CreateLanguage,
		http.MethodPost,
		cfg,
	)

	HandleAdmin(
		languages,
		"/{code}",
		handlers.Language.UpdateLanguage,
		http.MethodPatch,
		cfg,
	)

	HandleAdmin(
		languages,
		"/{code}",
		handlers.Language.DeleteLanguage,
		http.MethodDelete,
		cfg,
	)
}
