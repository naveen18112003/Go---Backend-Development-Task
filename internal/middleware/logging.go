package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RequestLogger logs request/response lifecycle with duration.
func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		logger.Info("request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration_ms", duration),
		)

		return err
	}
}



