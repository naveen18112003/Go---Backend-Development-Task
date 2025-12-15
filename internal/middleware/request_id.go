package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestID injects a unique request id into the response headers.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(requestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}
		c.Set(requestIDHeader, id)
		return c.Next()
	}
}



