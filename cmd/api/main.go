package main

import (
	"context"
	"log"
	"net/http"

	"inariops/internal/api"
	"inariops/internal/bootstrap"
	"inariops/internal/db"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/middleware"
)

func main() {

	// =====================
	// Database
	// =====================

	dbConn := db.Connect()

	// =====================
	// Application
	// =====================

	app := bootstrap.New(dbConn)

	// =====================
	// Router
	// =====================

	router := api.NewRouter(app.Handlers)

	router.Use(middleware.CORS)
	router.Use(logger.Logging)

	// =====================
	// Background Workers
	// =====================

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.StartWorkers(ctx)

	// =====================
	// HTTP Server
	// =====================

	log.Println("InariOps running on :9142")

	server := &http.Server{
		Addr:    ":9142",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
