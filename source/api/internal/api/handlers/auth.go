package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// AuthHandler maneja endpoints de autenticacion
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler crea un nuevo handler de autenticacion
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// LoginRequest estructura de request de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse estructura de response de login
type LoginResponse struct {
	Token   string         `json:"token"`
	Usuario *models.Usuario `json:"usuario"`
}

// Login maneja el login de usuarios
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password required"})
	}

	// Buscar usuario por email
	var usuario models.Usuario
	if err := h.db.Where("email = ? AND activo = ?", req.Email, true).First(&usuario).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Verificar contrasena
	if !auth.CheckPassword(req.Password, usuario.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generar token JWT
	token, err := auth.GenerateToken(&usuario)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	// Responder con token y usuario
	response := LoginResponse{
		Token:   token,
		Usuario: &usuario,
	}

	return c.JSON(response)
}

// Logout maneja el logout (client-side, solo informativo)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// RefreshToken renueva un token JWT
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	newToken, err := auth.RefreshToken(req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": newToken})
}

// Me retorna el usuario actual
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	var usuario models.Usuario
	if err := h.db.First(&usuario, "id = ?", claims.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(usuario)
}

// Permisos retorna permisos del rol actual y matriz rol->permisos
func (h *AuthHandler) Permisos(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"rol":          claims.Rol,
		"mis_permisos": auth.ObtenerPermisos(claims.Rol),
		"por_rol":      auth.PermisosPorRol(),
	})
}
