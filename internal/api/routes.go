package api

import (
	"database/sql"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/reservations"
	"inariops/internal/modules/users"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	authRepo := auth.NewRepository(db)
	userRepo := users.NewRepository(db)
	reservationRepo := reservations.NewRepository(db)

	authService := auth.NewService(authRepo)
	userService := users.NewService(userRepo, authService)
	reservationService := reservations.NewService(reservationRepo)

	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)
	reservationHandler := reservations.NewHandler(reservationService)

	router.HandleFunc("/auth", authHandler.HandleAuth)
	router.HandleFunc("/users", userHandler.HandleUsers)
	router.HandleFunc("/reservations", reservationHandler.HandleReservations)

	return router
}
