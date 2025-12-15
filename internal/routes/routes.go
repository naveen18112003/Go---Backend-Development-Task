package routes

import (
	"user-age-api/internal/handler"

	"github.com/gofiber/fiber/v2"
)

// Register registers API routes.
func Register(app *fiber.App, userHandler *handler.UserHandler) {
	// The handler expects to bind directly to /users endpoints.
	userHandler.Register(app)
}

