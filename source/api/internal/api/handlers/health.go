package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/school-monitoring/backend/internal/api/middleware"
)

// Health responde para healthchecks (Railway / LB)
func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"ok":         true,
		"service":    "backend",
		"request_id": middleware.GetRequestID(c),
		"time":       time.Now().UTC().Format(time.RFC3339),
	})
}


