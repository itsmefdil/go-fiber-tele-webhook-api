package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TelegramBot struct for storing Telegram Bot data
type TelegramBot struct {
	gorm.Model
	Token    string `json:"token"`
	RoomID   string `json:"room_id"`
	ThreadID string `json:"thread_id"`
}

// MessageRequest struct for parsing message request
type MessageRequest struct {
	Message string `json:"message"`
}

var DB *gorm.DB

// Initialize Database
func initDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("telegram_bot.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	fmt.Println("Database connected!")

	// Auto migrate model
	DB.AutoMigrate(&TelegramBot{})
	fmt.Println("Database Migrated")
}

// Load environment variables from .env file
func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Load environment variables
	loadEnv()

	// Inisialisasi aplikasi Fiber
	app := fiber.New()

	// Inisialisasi Database
	initDatabase()

	// Routes with Basic Auth middleware
	app.Get("/bots", BasicAuth(GetBots))
	app.Get("/bots/:id", BasicAuth(GetBot))
	app.Post("/bots", BasicAuth(CreateBot))
	app.Put("/bots/:id", BasicAuth(UpdateBot))
	app.Delete("/bots/:id", BasicAuth(DeleteBot))
	app.Post("/webhook/:id/send", BasicAuth(SendTelegramMessage)) // Route untuk webhook

	// Jalankan server di port 3000
	log.Fatal(app.Listen(":3000"))
}

// BasicAuth middleware for protecting routes
func BasicAuth(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := os.Getenv("BASIC_AUTH_USERNAME")
		password := os.Getenv("BASIC_AUTH_PASSWORD")

		// Get the basic auth credentials from request
		auth := c.Get("Authorization")
		if auth == "" {
			c.Status(fiber.StatusUnauthorized)
			return c.SendString("Unauthorized")
		}

		// Check if the authorization header is valid
		if !checkAuth(auth, username, password) {
			c.Status(fiber.StatusUnauthorized)
			return c.SendString("Unauthorized")
		}

		// Continue to next handler
		return next(c)
	}
}

// Helper function to check authorization credentials
func checkAuth(authHeader, validUsername, validPassword string) bool {
	// Decode Basic Auth Header (Authorization: Basic base64(username:password))
	const prefix = "Basic "
	if len(authHeader) < len(prefix) {
		return false
	}

	auth := authHeader[len(prefix):]
	decodedAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return false
	}

	// Split username and password
	credentials := string(decodedAuth)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return false
	}

	username, password := parts[0], parts[1]

	return username == validUsername && password == validPassword
}

// Handlers for CRUD and Webhook

// Get all bots
func GetBots(c *fiber.Ctx) error {
	var bots []TelegramBot
	DB.Find(&bots)
	return c.JSON(bots)
}

// Get a single bot by id
func GetBot(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot TelegramBot
	if result := DB.First(&bot, id); result.Error != nil {
		return c.Status(404).SendString("Bot not found")
	}
	return c.JSON(bot)
}

// Create a new bot
func CreateBot(c *fiber.Ctx) error {
	bot := new(TelegramBot)
	if err := c.BodyParser(bot); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	DB.Create(&bot)
	return c.JSON(bot)
}

// Update a bot
func UpdateBot(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot TelegramBot
	if result := DB.First(&bot, id); result.Error != nil {
		return c.Status(404).SendString("Bot not found")
	}

	if err := c.BodyParser(&bot); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	DB.Save(&bot)
	return c.JSON(bot)
}

// Delete a bot
func DeleteBot(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot TelegramBot
	if result := DB.First(&bot, id); result.Error != nil {
		return c.Status(404).SendString("Bot not found")
	}
	DB.Delete(&bot)
	return c.SendString("Bot deleted successfully")
}

// Send a message to Telegram
func SendTelegramMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot TelegramBot
	if result := DB.First(&bot, id); result.Error != nil {
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
		"chat_id": bot.RoomID,
		"text":    messageRequest.Message,
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
