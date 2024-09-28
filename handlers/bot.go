package handlers

import (
	"go-tele-webhook/database"
	"go-tele-webhook/models"

	"github.com/gofiber/fiber/v2"
)

// GetBots retrieves all bots
func GetBots(c *fiber.Ctx) error {
	var bots []models.TelegramBot
	database.DB.Find(&bots)
	return c.JSON(bots)
}

// GetBot retrieves a bot by id
func GetBot(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot models.TelegramBot
	if result := database.DB.First(&bot, id); result.Error != nil {
		return c.Status(404).SendString("Bot not found")
	}
	return c.JSON(bot)
}

// CreateBot creates a new bot
func CreateBot(c *fiber.Ctx) error {
	bot := new(models.TelegramBot)
	if err := c.BodyParser(bot); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	database.DB.Create(&bot)
	return c.JSON(bot)
}

// UpdateBot updates an existing bot
func UpdateBot(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot models.TelegramBot
	if result := database.DB.First(&bot, id); result.Error != nil {
		return c.Status(404).SendString("Bot not found")
	}

	if err := c.BodyParser(&bot); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	database.DB.Save(&bot)
	return c.JSON(bot)
}

// DeleteBot deletes a bot by id
func DeleteBot(c *fiber.Ctx) error {
	id := c.Params("id")
	var bot models.TelegramBot
	if result := database.DB.First(&bot, id); result.Error != nil {
		return c.Status(404).SendString("Bot not found")
	}
	database.DB.Delete(&bot)
	return c.SendString("Bot deleted successfully")
}
