package api

import (
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

func registerUserRoutes(
	router *mux.Router,
	handlers *bootstrap.Handlers,
	cfg config.Config,
) {
	users := router.PathPrefix("/users").Subrouter()

	// =========================================================
	// ADMIN ONLY
	// =========================================================

	// Get all users.
	HandleAdmin(
		users,
		"",
		handlers.User.GetAllUsers,
		http.MethodGet,
		cfg,
	)

	// Create user.
	HandleAdmin(
		users,
		"",
		handlers.User.CreateUser,
		http.MethodPost,
		cfg,
	)

	// Get all guides with their user information.
	HandleAdmin(
		users,
		"/guides",
		handlers.User.GetAllGuides,
		http.MethodGet,
		cfg,
	)

	// Deactivate user.
	HandleAdmin(
		users,
		"/{id}",
		handlers.User.DeactivateUser,
		http.MethodDelete,
		cfg,
	)

	// =========================================================
	// ADMIN + GUIDE (OWN USER)
	// =========================================================

	// ADMIN:
	//   can access any user.
	//
	// GUIDE:
	//   can only access the user represented by the JWT user_id.
	HandleAdminOrOwnUser(
		users,
		"/{id}",
		handlers.User.GetUserByID,
		http.MethodGet,
		cfg,
	)

	// ADMIN:
	//   can update any user.
	//
	// GUIDE:
	//   can only update their own user.
	HandleAdminOrOwnUser(
		users,
		"/{id}",
		handlers.User.UpdateUser,
		http.MethodPatch,
		cfg,
	)
}
