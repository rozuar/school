package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// CORSMiddleware middleware para manejar CORS (Fiber).
func CORSMiddleware(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-Id")
	c.Set("Access-Control-Max-Age", "86400")

	if c.Method() == fiber.MethodOptions {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Next()
}

// JSONMiddleware middleware para establecer Content-Type JSON (Fiber).
func JSONMiddleware(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Next()
}
