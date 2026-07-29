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

	dbConn := db.Connect()

	app := bootstrap.New(dbConn)

	router := api.NewRouter(dbConn)

	router.Use(middleware.CORS)
	router.Use(logger.Logging)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.Start(ctx)

	log.Println("InariOps running on :9142")

	if err := http.ListenAndServe(":9142", router); err != nil {
		log.Fatal(err)
	}
}
