package api

import (
	"inariops/internal/config"
	"inariops/internal/domain"
	"inariops/internal/shared/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleAdmin(
	router *mux.Router,
	path string,
	handler http.HandlerFunc,
	method string,
	cfg config.Config,
) {
	router.Handle(
		path,
		middleware.RequireAdmin(cfg)(
			http.HandlerFunc(handler),
		),
	).Methods(method)
}

func HandleGuide(
	router *mux.Router,
	path string,
	handler http.HandlerFunc,
	method string,
	cfg config.Config,
) {
	router.Handle(
		path,
		middleware.RequireGuide(cfg)(
			http.HandlerFunc(handler),
		),
	).Methods(method)
}

func HandleAdminOrGuide(
	router *mux.Router,
	path string,
	handler http.HandlerFunc,
	method string,
	cfg config.Config,
) {
	router.Handle(
		path,
		middleware.RequireAdminOrGuide(cfg)(
			http.HandlerFunc(handler),
		),
	).Methods(method)
}

func HandleRoles(
	router *mux.Router,
	path string,
	handler http.HandlerFunc,
	method string,
	cfg config.Config,
	roles ...domain.UserRole,
) {
	router.Handle(
		path,
		middleware.RequireRole(cfg, roles...)(
			http.HandlerFunc(handler),
		),
	).Methods(method)
}

func HandleAdminOrOwnUser(
	router *mux.Router,
	path string,
	handler http.HandlerFunc,
	method string,
	cfg config.Config,
) {
	router.Handle(
		path,
		middleware.RequireAdminOrOwnUser(cfg)(
			http.HandlerFunc(handler),
		),
	).Methods(method)
}

func HandleAdminOrOwnGuide(
	router *mux.Router,
	path string,
	handler http.HandlerFunc,
	method string,
	cfg config.Config,
) {
	router.Handle(
		path,
		middleware.RequireAdminOrOwnUser(cfg)(
			http.HandlerFunc(handler),
		),
	).Methods(method)
}
