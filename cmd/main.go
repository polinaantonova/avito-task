package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"

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

	app := fiber.New()
	app.Get("/api/ping", adaptor.HTTPHandler(pingHandler))
	app.Get("/api/dbconnect", adaptor.HTTPHandler(dBConnector))
	app.Post("/api/tenders/new", adaptor.HTTPHandler(tenderCreator))
	app.Get("/api/tenders", adaptor.HTTPHandler(tendersByService))
	app.Get("/api/tenders/my", adaptor.HTTPHandler(tenderByUser))
	app.Get("/api/tenders/:tenderID/status", func(ctx fiber.Ctx) error {
		tenderID := ctx.Params("tenderID", "")
		fmt.Println("tenderID: ", tenderID)
		if tenderID == "" {
			return fiber.ErrBadRequest
		}

		for _, tender := range tenders.List() {
			if tender.Id == tenderID {
				ctx.SendString(tender.Status)
			}
		}
		return fiber.ErrNotFound
	})

	err := app.Listen(serverAddress)
	if err != nil {
		log.Fatal(err)
	}
}
