package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"

	"polina.com/m/internal/handlers/allTenders"
	"polina.com/m/internal/handlers/createTender"
	"polina.com/m/internal/handlers/dbconnect"
	"polina.com/m/internal/handlers/ping"
	"polina.com/m/internal/handlers/tendersByUser"
)

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "0.0.0.0:8080"
	}

	//подключаюсь к postgres
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DATABASE")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {

		//errorText := fmt.Sprintf("empty env variables\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
		//return errors.New(errorText)

		user = "cnrprod1725725190-team-78136"
		password = "cnrprod1725725190-team-78136"
		host = "rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net"
		port = "6432"
		dbName = "cnrprod1725725190-team-78136"
	}

	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		errorText := fmt.Sprintf("cannot connect to database\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
		log.Fatal(errorText)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		errorText := fmt.Sprintf("cannot ping database\n host: %v\n port: %v\n user: %v\n password: %v\n dbName: %v\n", host, port, user, password, dbName)
		log.Fatal(errorText)
	}

	fmt.Println("Successfully connected!")

	sqlStatement := `CREATE TABLE IF NOT EXISTS tenders(
id UUID,
name VARCHAR(100),
description VARCHAR(500), 
service_type VARCHAR(50),
status VARCHAR(50) DEFAULT 'Created',
organization_id UUID,
creator_username VARCHAR(100),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
version INT DEFAULT 1
);
`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Fatal("error while creating table tenders")
	}

	//----------

	pingHandler := ping.NewPing()
	dBConnector := dbconnect.NewDBConnector(db)
	tenderCreator := createTender.NewTenderCreator(db)
	tendersByService := allTenders.NewAllTenders(db)
	tenderByUser := tendersByUser.NewTendersByUser(db)

	app := fiber.New()

	app.Get("/api/ping", adaptor.HTTPHandler(pingHandler))
	app.Get("/api/dbconnect", adaptor.HTTPHandler(dBConnector))
	app.Post("/api/tenders/new", adaptor.HTTPHandler(tenderCreator))
	app.Get("/api/tenders", adaptor.HTTPHandler(tendersByService))
	app.Get("/api/tenders/my", adaptor.HTTPHandler(tenderByUser))
	app.Get("/api/tenders/:tenderID/status", func(ctx fiber.Ctx) error {
		tenderID := ctx.Params("tenderID", "")
		if tenderID == "" {
			return fiber.ErrBadRequest
		}

		var status string

		err = db.QueryRow("SELECT status FROM Tenders WHERE id = $1", tenderID).Scan(&status)
		if err == sql.ErrNoRows {
			return fiber.ErrNotFound
		}
		if err != nil {
			return fiber.ErrInternalServerError
		}

		response, err := json.Marshal(status)
		if err != nil {
			return fiber.ErrBadRequest
		}

		ctx.Status(http.StatusOK)
		ctx.Write(response)
		return nil
	})

	//app.Put("/api/tenders/:tenderID/status", func(ctx fiber.Ctx) error {
	//	tenderID := ctx.Params("tenderID", "")
	//	if tenderID == "" {
	//		return fiber.ErrBadRequest
	//	}
	//
	//
	//})

	err = app.Listen(serverAddress)
	if err != nil {
		log.Fatal(err)
	}
}
