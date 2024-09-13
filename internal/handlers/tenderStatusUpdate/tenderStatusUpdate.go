package tenderStatusUpdate

import (
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"polina.com/m/internal/errorMessage"
	"polina.com/m/internal/tender"
)

func TenderStatusUpdate(ctx fiber.Ctx, db *sql.DB) error {
	tenderID := ctx.Params("tenderID", "")
	if tenderID == "" {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, "Please specify tender id")
	}

	myTender := tender.NewTender()
	myTender.Id = tenderID

	body := ctx.Body()
	err := json.Unmarshal(body, &myTender)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	//проверка статуса на корректность
	err = myTender.ValidateStatus()
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	query := `UPDATE tenders SET status = $1 WHERE id = $2`
	result, err := db.Exec(query, myTender.Status, tenderID)
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusBadRequest, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if rowsAffected == 0 {
		return errorMessage.SendErrorMessageFiber(ctx, fiber.StatusNotFound, "no tender with this id found")
	}
	//добавить проверку юзера

	ctx.Status(http.StatusOK)
	ctx.Write([]byte("tender status updated, refresh page to see update"))
	return nil
}
