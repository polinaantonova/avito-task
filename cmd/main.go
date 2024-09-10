package main

import (
	"log"
	"net/http"
	"os"
	"polina.com/m/internal/handlers/createTender"
	"polina.com/m/internal/handlers/dbconnect"
	"polina.com/m/internal/handlers/ping"
	"polina.com/m/internal/tender"
)

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "0.0.0.0:8080"
	}
	tenders := tender.NewTenderList()

	pingHandler := ping.NewPing()
	dBConnector := dbconnect.NewDBConnector()
	tenderCreator := createTender.NewTenderCreator(tenders)

	http.Handle("/api/ping", pingHandler)
	http.Handle("/api/dbconnect", dBConnector)
	http.Handle("/api/tenders/new", tenderCreator)

	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}
