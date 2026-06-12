package main

import (
	"log"
	"net/http"

	"inariops/internal/api"
	"inariops/internal/db"
	"inariops/internal/shared/logger"
)

func main() {
	dbConn := db.Connect()

	router := api.NewRouter(dbConn)
	router.Use(logger.Logging) // Middleware para logging

	log.Println("InariOps running on :9142")

	err := http.ListenAndServe(":9142", router)
	if err != nil {
		log.Fatal(err)
	}
}
