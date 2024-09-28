package main

import (
	"log"

	"go-tele-webhook/config"
	"go-tele-webhook/database"
	"go-tele-webhook/handlers"
	"go-tele-webhook/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize the database
	database.InitDatabase()

	// Initialize Fiber app
	app := fiber.New()

	// Routes with Basic Auth middleware
	app.Get("/bots", middleware.BasicAuth(handlers.GetBots))
	app.Get("/bots/:id", middleware.BasicAuth(handlers.GetBot))
	app.Post("/bots", middleware.BasicAuth(handlers.CreateBot))
	app.Put("/bots/:id", middleware.BasicAuth(handlers.UpdateBot))
	app.Delete("/bots/:id", middleware.BasicAuth(handlers.DeleteBot))

	// Webhook route
	app.Post("/webhook/:id/send", middleware.BasicAuth(handlers.SendTelegramMessage))

	// Start server on port 3000
	log.Fatal(app.Listen(":3000"))
}
