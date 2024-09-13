package rollbackTender

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"polina.com/m/internal/errorMessage"
	"polina.com/m/internal/tender"
	"strconv"
)

func RollbackTender(ctx fiber.Ctx, db *sql.DB) error {
	tenderID := ctx.Params("tenderID", "")
	versionStr := ctx.Params("version", "")
	if tenderID == "" {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "Please specify tender id")
	}

	if versionStr == "" {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "Please specify version")
	}

	version, _ := strconv.Atoi(versionStr)

	myTender := tender.NewTender()
	myTender.Id = tenderID

	body := ctx.Body()
	err := json.Unmarshal(body, &myTender)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}
	//тут доделать проверку по юзеру

	err = db.QueryRow("SELECT name, description, service_type, status, organization_id, creator_username FROM tenders WHERE id = $1 AND version = $2", myTender.Id, version).Scan(
		&myTender.Name,
		&myTender.Description,
		&myTender.ServiceType,
		&myTender.Status,
		&myTender.OrganizationId,
		&myTender.CreatorUsername,
	)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusNotFound, "tender id or version not found")
		}
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	//ищем макс верию в таблице
	query := `SELECT version FROM tenders WHERE id = $1 ORDER BY version DESC LIMIT 1`
	err = db.QueryRow(query, tenderID).Scan(&myTender.Version)
	if err != nil {
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusNotFound, "tender id or version not found")
			}
			return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
		}
	}
	//и еще увеличиваем
	myTender.Version++
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
