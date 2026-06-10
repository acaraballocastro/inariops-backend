package main

import (
	"log"
	"net/http"

	"inariops/internal/api"
	"inariops/internal/db"
)

func main() {
	dbConn := db.Connect()

	router := api.NewRouter(dbConn)

	log.Println("InariOps running on :9142")

	err := http.ListenAndServe(":9142", router)
	if err != nil {
		log.Fatal(err)
	}
}
