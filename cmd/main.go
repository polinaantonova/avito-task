package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"polina.com/m/internal/tender"

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
postgres_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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
		log.Fatal(err.Error())
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
			return fiber.NewError(fiber.StatusBadRequest, "Please specify tender id")
		}

		myTender := tender.NewTender()
		myTender.Id = tenderID

		err = db.QueryRow("SELECT status FROM Tenders WHERE id = $1 ORDER BY version DESC LIMIT 1", tenderID).Scan(&myTender.Status)
		if errors.Is(sql.ErrNoRows, err) {
			return fiber.NewError(fiber.StatusNotFound, "No tenders with this id")
		}
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid tender format")
		}

		response, err := json.Marshal(myTender.Status)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		//тут можно сделать проверку, показывать ли тендер в зависимости от юзернейм

		ctx.Status(http.StatusOK)
		ctx.Write(response)
		return nil
	})

	app.Put("/api/tenders/:tenderID/status", func(ctx fiber.Ctx) error {
		tenderID := ctx.Params("tenderID", "")
		if tenderID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Please specify tender id")
		}

		myTender := tender.NewTender()
		myTender.Id = tenderID

		body := ctx.Body()
		err = json.Unmarshal(body, &myTender)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		//проверка статуса на корректность
		err = myTender.ValidateStatus()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		query := `UPDATE tenders SET status = $1 WHERE id = $2`
		result, err := db.Exec(query, myTender.Status, tenderID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if rowsAffected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "no tender with this id found")
		}
		//добавить проверку юзера

		ctx.Status(http.StatusOK)
		ctx.Write([]byte("tender status updated, refresh page to see update"))
		return nil
	})

	app.Patch("/api/tenders/:tenderID/edit", func(ctx fiber.Ctx) error {
		tenderID := ctx.Params("tenderID", "")
		if tenderID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Please specify tender id")
		}

		myTender := tender.NewTender()
		myTender.Id = tenderID

		body := ctx.Body()
		err = json.Unmarshal(body, &myTender)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		//проверки, добавить проверку юзера

		err = myTender.ValidateTenderServiceType()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		err = myTender.ValidateStringFieldsLen()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		err = myTender.ValidateStatus()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		// Поиск записи с нужным ID
		existingTender := tender.NewTender()
		err = db.QueryRow("SELECT name, description, service_type, status, organization_id, creator_username, version FROM tenders WHERE id = $1 ORDER BY version DESC LIMIT 1", myTender.Id).Scan(
			&existingTender.Name,
			&existingTender.Description,
			&existingTender.ServiceType,
			&existingTender.Status,
			&existingTender.OrganizationId,
			&existingTender.CreatorUsername,
			&existingTender.Version,
		)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				return fiber.NewError(fiber.StatusNotFound, "tender with this id not found")
			}
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if myTender.Name == "" {
			myTender.Name = existingTender.Name
		}
		if myTender.Description == "" {
			myTender.Description = existingTender.Description
		}
		if myTender.ServiceType == "" {
			myTender.ServiceType = existingTender.ServiceType
		}

		myTender.OrganizationId = existingTender.OrganizationId

		myTender.CreatorUsername = existingTender.CreatorUsername

		// Увеличиваем версию
		myTender.Version = existingTender.Version + 1

		query := "INSERT INTO tenders (id, name, description, service_type, status, organization_id, creator_username, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		_, err := db.Exec(query, myTender.Id, myTender.Name, myTender.Description, myTender.ServiceType, myTender.Status, myTender.OrganizationId, myTender.CreatorUsername, myTender.Version)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		response, err := json.Marshal(myTender)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		ctx.Status(http.StatusOK)
		ctx.Write(response)
		return nil

	})

	err = app.Listen(serverAddress)
	if err != nil {
		log.Fatal(err)
	}
}
