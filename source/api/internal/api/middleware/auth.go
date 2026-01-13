package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/school-monitoring/backend/internal/auth"
)

const (
	LocalsUserKey      = "user"
	LocalsRequestIDKey = "request_id"
)

// AuthMiddleware middleware de autenticacion JWT (Fiber).
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization format"})
	}

	claims, err := auth.ValidateToken(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	c.Locals(LocalsUserKey, claims)
	return c.Next()
}

// RoleMiddleware middleware que verifica el rol del usuario (Fiber).
func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := GetUserFromContext(c)
		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found in context"})
		}

		roleAllowed := false
		for _, role := range roles {
			if claims.Rol == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
		}

		return c.Next()
	}
}

// PermissionMiddleware middleware que verifica permisos especificos (Fiber).
func PermissionMiddleware(permisos ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := GetUserFromContext(c)
		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found in context"})
		}

		if !auth.TieneAlgunPermiso(claims.Rol, permisos...) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
		}

		return c.Next()
	}
}

// GetUserFromContext obtiene los claims del usuario desde locals (Fiber).
func GetUserFromContext(c *fiber.Ctx) *auth.Claims {
	if v := c.Locals(LocalsUserKey); v != nil {
		if claims, ok := v.(*auth.Claims); ok {
			return claims
		}
	}
	return nil
}
