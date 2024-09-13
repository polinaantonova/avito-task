package tenderStatus

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"polina.com/m/internal/errorMessage"
	"polina.com/m/internal/tender"
)

func TenderStatus(ctx fiber.Ctx, db *sql.DB) error {
	tenderID := ctx.Params("tenderID", "")
	if tenderID == "" {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "Please specify tender id")
	}

	myTender := tender.NewTender()
	myTender.Id = tenderID

	err := db.QueryRow("SELECT status FROM tenders WHERE id = $1 ORDER BY version DESC LIMIT 1", tenderID).Scan(&myTender.Status)
	if errors.Is(sql.ErrNoRows, err) {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusNotFound, "No tenders with this id")
	}
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "Invalid tender format")
	}

	response, err := json.Marshal(myTender.Status)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	//??тут нужно сделать проверку, показывать ли тендер в зависимости от юзернейм

	ctx.Status(http.StatusOK)
	ctx.Write(response)
	return nil
}
