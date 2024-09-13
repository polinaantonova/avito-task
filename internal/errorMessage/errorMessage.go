package errorMessage

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

type ErrorMessage struct {
	Reason string `json:"reason"`
}

func NewErrorMessage(message string) ErrorMessage {
	return ErrorMessage{Reason: message}
}

func SendErrorMessage(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorMessage := NewErrorMessage(message)
	response, err := json.Marshal(errorMessage)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}
	w.Write(response)
}

func SendErrorMessageFiber(ctx fiber.Ctx, statusCode int, message string) error {
	errorMessage := NewErrorMessage(message)
	ctx.Status(statusCode).JSON(errorMessage)
	return nil
}
