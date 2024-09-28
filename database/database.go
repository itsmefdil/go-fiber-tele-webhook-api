package database

import (
	"fmt"
	"go-tele-webhook/models" // Import models package
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a global variable to hold the database connection
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("telegram_bot.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	fmt.Println("Database connected!")

	// Auto migrate tables
	MigrateDatabase()
}

// MigrateDatabase migrates all models
func MigrateDatabase() {
	DB.AutoMigrate(&models.TelegramBot{}) // Gunakan models.TelegramBot
	fmt.Println("Database Migrated")
}
