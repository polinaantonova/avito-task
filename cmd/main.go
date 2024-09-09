package main

import (
	"log"
	"net/http"
	"os"
	"polina.com/m/internal/handlers/createTender"
	"polina.com/m/internal/handlers/dbconnect"
	"polina.com/m/internal/handlers/ping"
)

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "0.0.0.0:8080"
	}

	pingHandler := ping.NewPing()
	dBConnector := dbconnect.NewDBConnector()
	tenderCreator := createTender.NewTenderCreator()

	http.Handle("/api/ping", pingHandler)
	http.Handle("/api/dbconnect", dBConnector)
	http.Handle("/api/tenders/new", tenderCreator)

	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}
