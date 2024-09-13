package tenderStatusUpdate

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"net/http"
	"polina.com/m/internal/errorMessage"
	"polina.com/m/internal/tender"
)

func TenderStatusUpdate(ctx fiber.Ctx, db *sql.DB) error {
	tenderIDStr := ctx.Params("tenderID", "")
	if tenderIDStr == "" {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "Please specify tender id")
	}

	tenderID, err := uuid.Parse(tenderIDStr)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "unvalid tender format")
	}

	statusUpdateValidator := tender.NewStatusUpdateValidator()
	body := ctx.Body()
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&statusUpdateValidator)

	//проверка json
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, http.StatusBadRequest, "check JSON schema: you can update only status")
	}

	//проверка корректности статуса
	if statusUpdateValidator.Status == "" {
		return errorMessage.SendErrorMessageFiber(ctx, http.StatusBadRequest, "please specify status")
	}

	err = tender.ValidateStatus(statusUpdateValidator.Status)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, http.StatusBadRequest, err.Error())
	}

	//проверка username
	if statusUpdateValidator.Username == "" {
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
	err = db.QueryRow(query, statusUpdateValidator.Username, tenderID).Scan(&exists)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	if !exists {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusForbidden, "this user is not allowed to modify tender status")
	}

	query = `UPDATE tenders SET status = $1 WHERE id = $2`
	_, err = db.Exec(query, statusUpdateValidator.Status, tenderID)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	myTender := tender.NewTender()
	query = `SELECT id, name, description, service_type, status, organization_id, creator_username, created_at, version FROM tenders`

	err = db.QueryRow("SELECT name, description, service_type, status, organization_id, creator_username, version FROM tenders "+
		"WHERE id = $1 ORDER BY version DESC LIMIT 1", tenderID).Scan(
		&myTender.Name,
		&myTender.Description,
		&myTender.ServiceType,
		&myTender.Status,
		&myTender.OrganizationId,
		&myTender.CreatorUsername,
		&myTender.Version,
	)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}
	response, err := json.Marshal(myTender)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, http.StatusBadRequest, "bad JSON")
	}

	ctx.Status(http.StatusOK)
	ctx.Write([]byte("tender status updated, refresh page to see update\n"))
	ctx.Write(response)
	return nil
}
