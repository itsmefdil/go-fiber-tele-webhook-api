package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-tele-webhook/database"
	"go-tele-webhook/models"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type MessageRequest struct {
	Message string `json:"message"`
}

// SendTelegramMessage sends a message to Telegram using bot token and room ID
func SendTelegramMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot models.TelegramBot
	if result := database.DB.First(&bot, id); result.Error != nil {
		return c.Status(404).SendString("Bot not found")
	}

	// Parse the message request
	messageRequest := new(MessageRequest)
	if err := c.BodyParser(messageRequest); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Prepare the Telegram API request
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot.Token)

	// Prepare the payload
	payload := map[string]interface{}{
		"chat_id":           bot.RoomID,
		"message_thread_id": bot.ThreadID,
		"text":              messageRequest.Message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return c.Status(500).SendString("Failed to marshal payload")
	}

	// Send the request to the Telegram API
	resp, err := http.Post(telegramURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(500).SendString("Failed to send message to Telegram")
	}

	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response body:", string(respBody))

	return c.SendString("Message sent successfully")
}
