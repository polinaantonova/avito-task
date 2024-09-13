package editTender

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"polina.com/m/internal/errorMessage"
	"polina.com/m/internal/tender"
	"polina.com/m/internal/validations"
)

func EditTender(ctx fiber.Ctx, db *sql.DB) error {
	tenderID := ctx.Params("tenderID", "")
	if tenderID == "" {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "Please specify tender id")
	}

	tenderEditor := validations.NewTenderEditValidator()
	body := ctx.Body()
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&tenderEditor)

	//проверка json
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, http.StatusBadRequest, "check JSON schema: you can update only name, description, service_type")
	}

	//проверки
	if tenderEditor.ServiceType != "" {
		err = tender.ValidateTenderServiceType(tenderEditor.ServiceType)
		if err != nil {
			return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
		}
	}

	err = tenderEditor.ValidateStringFieldsLen()
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	//проверка username
	if tenderEditor.Username == "" {
		return errorMessage.SendErrorMessageFiber(ctx, http.StatusUnauthorized, "please specify username")
	}
	var exists bool

	//есть ли в таблице тендер с таким id?
	query := `SELECT EXISTS (
    SELECT 1
    FROM tenders
    WHERE id = $1
);`
	err = db.QueryRow(query, tenderID).Scan(&exists)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	if !exists {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "tender does not exist")
	}

	//проверка прав. не понимаю, что конкретно надо сделать, пусть изменять тендер сможет только тот, кто его создал
	query = `SELECT EXISTS (
		SELECT 1
		FROM tenders
		WHERE creator_username = $1 AND id = $2
	)`
	err = db.QueryRow(query, tenderEditor.Username, tenderID).Scan(&exists)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	if !exists {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusForbidden, "this user is not allowed to modify tender fields")
	}
	myTender := tender.NewTender()
	myTender.Id = tenderID
	myTender.Name = tenderEditor.Name
	myTender.Description = tenderEditor.Description
	myTender.ServiceType = tenderEditor.ServiceType

	// Поиск записи с нужным ID
	existingTender := tender.NewTender()
	err = db.QueryRow("SELECT name, description, service_type, status, organization_id, creator_username, version FROM tenders WHERE id = $1 ORDER BY version DESC LIMIT 1", tenderID).Scan(
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
			return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusNotFound, "tender with this id not found")
		}
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
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

	query = "INSERT INTO tenders (id, name, description, service_type, status, organization_id, creator_username, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = db.Exec(query, myTender.Id, myTender.Name, myTender.Description, myTender.ServiceType, myTender.Status, myTender.OrganizationId, myTender.CreatorUsername, myTender.Version)

	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusInternalServerError, err.Error())
	}

	response, err := json.Marshal(myTender)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	ctx.Status(http.StatusOK)
	ctx.Write(response)
	return nil
}
