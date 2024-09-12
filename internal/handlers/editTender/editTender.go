package editTender

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"polina.com/m/internal/tender"
)

func EditTender(ctx fiber.Ctx, db *sql.DB) error {
	tenderID := ctx.Params("tenderID", "")
	if tenderID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please specify tender id")
	}

	myTender := tender.NewTender()
	myTender.Id = tenderID

	body := ctx.Body()
	err := json.Unmarshal(body, &myTender)
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
	myTender.Status = existingTender.Status
	myTender.CreatorUsername = existingTender.CreatorUsername

	// Увеличиваем версию
	myTender.Version = existingTender.Version + 1

	query := "INSERT INTO tenders (id, name, description, service_type, status, organization_id, creator_username, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = db.Exec(query, myTender.Id, myTender.Name, myTender.Description, myTender.ServiceType, myTender.Status, myTender.OrganizationId, myTender.CreatorUsername, myTender.Version)

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
}
