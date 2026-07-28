package main

import (
	"context"
	"log"
	"net/http"

	"inariops/internal/api"
	"inariops/internal/db"

	reservationsapp "inariops/internal/application/reservations"
	tourdaysapp "inariops/internal/application/tour_days"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/guides"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	tourdays "inariops/internal/modules/tours/tour_days"

	"inariops/internal/shared/logger"
	"inariops/internal/shared/middleware"
	"inariops/internal/workers"
)

func main() {

	dbConn := db.Connect()

	// =====================
	// Router
	// =====================

	router := api.NewRouter(dbConn)

	// =====================
	// Middlewares
	// =====================

	router.Use(middleware.CORS)
	router.Use(logger.Logging)

	// =====================
	// Repositories
	// =====================

	tourDayappRepo := tourdaysapp.NewRepository(dbConn)
	tourDayRepo := tourdays.NewRepository(dbConn)
	guideRepo := guides.NewRepository(dbConn)

	reservationRepo := reservations.NewRepository(dbConn)

	customerRepo := customers.NewRepository(dbConn)

	reservationCustomerRepo := reservationscustomers.NewRepository(dbConn)

	// =====================
	// Services
	// =====================

	tourDayAppService := tourdaysapp.NewService(
		tourDayRepo,
		guideRepo,
		tourDayappRepo,
	)

	reservationAppService := reservationsapp.NewService(
		reservationRepo,
		customerRepo,
		reservationCustomerRepo,
		tourDayRepo,
	)

	// =====================
	// Workers
	// =====================

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	defer cancel()

	scheduler := workers.NewScheduler()

	scheduler.Register(
		workers.NewGuideAssignmentTimeoutWorker(
			tourDayAppService,
		),

		workers.NewReservationStatusWorker(
			reservationAppService,
			reservationRepo,
		),
	)

	scheduler.Start(ctx)

	// =====================
	// Server
	// =====================

	log.Println("InariOps running on :9142")

	if err := http.ListenAndServe(":9142", router); err != nil {
		log.Fatal(err)
	}
}
