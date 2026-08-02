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
}
