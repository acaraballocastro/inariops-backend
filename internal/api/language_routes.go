package api

import (
	"inariops/internal/bootstrap"
	"net/http"

	"github.com/gorilla/mux"
)

func registerLanguageRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
) {

	languages := router.PathPrefix("/languages").Subrouter()

	languages.HandleFunc("", handlers.Language.ListLanguages).
		Methods(http.MethodGet)

	languages.HandleFunc("", handlers.Language.CreateLanguage).
		Methods(http.MethodPost)

	languages.HandleFunc("/{code}", handlers.Language.GetLanguageByCode).
		Methods(http.MethodGet)

	languages.HandleFunc("/{code}", handlers.Language.UpdateLanguage).
		Methods(http.MethodPatch)

	languages.HandleFunc("/{code}", handlers.Language.DeleteLanguage).
		Methods(http.MethodDelete)
}
