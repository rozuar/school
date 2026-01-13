package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware asegura un request id para correlación (logs/headers).
func RequestIDMiddleware(c *fiber.Ctx) error {
	reqID := c.Get("X-Request-Id")
	if reqID == "" {
		reqID = uuid.New().String()
	}

	c.Set("X-Request-Id", reqID)
	c.Locals(LocalsRequestIDKey, reqID)
	return c.Next()
}

// GetRequestID obtiene el request id desde locals (si existe).
func GetRequestID(c *fiber.Ctx) string {
	if v := c.Locals(LocalsRequestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}


