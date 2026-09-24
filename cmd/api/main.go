package main

import (
	"context"
	"log"
	"net/http"

	"inariops/internal/api"
	"inariops/internal/bootstrap"
	"inariops/internal/config"
	"inariops/internal/db"
	"inariops/internal/modules/auth"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/middleware"
)

func main() {

	// =====================
	// Configuration
	// =====================

	cfg := config.Load()

	// =====================
	// Authentication
	// =====================

	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
	)

	// =====================
	// Database
	// =====================

	dbConn := db.Connect()

	// =====================
	// Application
	// =====================

	app := bootstrap.New(
		dbConn,
		jwtManager,
	)

	// =====================
	// Router
	// =====================

	router := api.NewRouter(
		app.Handlers,
		cfg,
		jwtManager,
	)

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

	log.Println("InariOps running on :9142 on " + cfg.AppEnv + " mode")

	server := &http.Server{
		Addr:    ":9142",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
