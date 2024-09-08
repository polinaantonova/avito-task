package main

import (
	"log"
	"net/http"
	"polina.com/m/internal/handlers/ping"
)

const portNum string = "0.0.0.0:8080"

func main() {
	pingHandler := ping.NewPing()

	http.Handle("/api/ping", pingHandler)

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}

}
