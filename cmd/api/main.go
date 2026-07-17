package main

import (
	"context"
	"log"
	"net/http"

	"inariops/internal/api"
	"inariops/internal/db"
	"inariops/internal/modules/guides"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/middleware"
	"inariops/internal/workers"
)

func main() {
	dbConn := db.Connect()

	// Router
	router := api.NewRouter(dbConn)

	// Middlewares
	router.Use(middleware.CORS)
	router.Use(logger.Logging)

	// =====================
	// Workers
	// =====================

	guidesRepo := guides.NewRepository(dbConn)
	tourDayRepo := tourdays.NewRepository(dbConn)

	tourDayService := tourdays.NewService(tourDayRepo, guidesRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := workers.NewScheduler()

	scheduler.Register(
		workers.NewGuideAssignmentTimeoutWorker(tourDayService),
	)

	scheduler.Start(ctx)

	// =====================

	log.Println("InariOps running on :9142")

	if err := http.ListenAndServe(":9142", router); err != nil {
		log.Fatal(err)
	}
}
