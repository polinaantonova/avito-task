package main

import (
	"log"
	"net/http"
	"os"
	"polina.com/m/internal/handlers/allTenders"
	"polina.com/m/internal/handlers/createTender"
	"polina.com/m/internal/handlers/dbconnect"
	"polina.com/m/internal/handlers/ping"
	"polina.com/m/internal/handlers/tendersByUser"
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
	tendersByService := allTenders.NewAllTenders(tenders)
	tenderByUser := tendersByUser.NewTendersByUser(tenders)
	//tenderStatus := tenderStatus2.NewTenderStatus(tenders)

	http.Handle("/api/ping", pingHandler)
	http.Handle("/api/dbconnect", dBConnector)
	http.Handle("/api/tenders/new", tenderCreator)
	http.Handle("/api/tenders", tendersByService)
	http.Handle("/api/tenders/my", tenderByUser)

	//router := httprouter.New()
	//router.GET("/tenders/:tenderID/status", tenderStatus)

	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
