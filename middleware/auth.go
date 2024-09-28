package middleware

import (
	"encoding/base64"
	"strings"

	"go-tele-webhook/config"

	"github.com/gofiber/fiber/v2"
)

// BasicAuth middleware for protecting routes
func BasicAuth(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil username dan password dari variabel lingkungan
		username := config.GetEnv("BASIC_AUTH_USERNAME", "defaultUsername")
		password := config.GetEnv("BASIC_AUTH_PASSWORD", "defaultPassword")

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
