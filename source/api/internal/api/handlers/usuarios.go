package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// UsuariosHandler maneja endpoints de administracion de usuarios
type UsuariosHandler struct {
	db *gorm.DB
}

func NewUsuariosHandler(db *gorm.DB) *UsuariosHandler {
	return &UsuariosHandler{db: db}
}

// GetAll obtiene usuarios, con filtros opcionales (rol, activo)
func (h *UsuariosHandler) GetAll(c *fiber.Ctx) error {
	query := h.db.Model(&models.Usuario{})

	if rol := c.Query("rol"); rol != "" {
		query = query.Where("rol = ?", rol)
	}
	if activo := c.Query("activo"); activo != "" {
		query = query.Where("activo = ?", activo == "true")
	}

	var usuarios []models.Usuario
	if err := query.Order("rol, nombre").Find(&usuarios).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching users"})
	}
	return c.JSON(usuarios)
}

func (h *UsuariosHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var usuario models.Usuario
	if err := h.db.First(&usuario, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(usuario)
}

type UsuarioRequest struct {
	Email    string  `json:"email"`
	Nombre   string  `json:"nombre"`
	Rol      string  `json:"rol"`
	Password *string `json:"password,omitempty"`
	Activo   *bool   `json:"activo,omitempty"`
}

func (h *UsuariosHandler) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)

	var req UsuarioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email == "" || req.Nombre == "" || req.Rol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email, nombre y rol son requeridos"})
	}
	if !models.EsRolValido(req.Rol) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "rol invalido"})
	}

	pass := "changeme123"
	if req.Password != nil && *req.Password != "" {
		pass = *req.Password
	}
	hash, err := auth.HashPassword(pass)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	usuario := models.Usuario{
		Email:        req.Email,
		Nombre:       req.Nombre,
		Rol:          req.Rol,
		PasswordHash: hash,
		Activo:       true,
	}
	if req.Activo != nil {
		usuario.Activo = *req.Activo
	}

	if err := h.db.Create(&usuario).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
	}
	_ = models.CrearAuditoria(h.db, "usuarios", usuario.ID, models.AuditoriaInsert, nil, &usuario, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	return c.Status(fiber.StatusCreated).JSON(usuario)
}

func (h *UsuariosHandler) Update(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var usuario models.Usuario
	if err := h.db.First(&usuario, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	before := usuario

	var req UsuarioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Email != "" {
		usuario.Email = req.Email
	}
	if req.Nombre != "" {
		usuario.Nombre = req.Nombre
	}
	if req.Rol != "" {
		if !models.EsRolValido(req.Rol) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "rol invalido"})
		}
		usuario.Rol = req.Rol
	}
	if req.Activo != nil {
		usuario.Activo = *req.Activo
	}
	if req.Password != nil && *req.Password != "" {
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
		}
		usuario.PasswordHash = hash
	}

	if err := h.db.Save(&usuario).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating user"})
	}

	_ = models.CrearAuditoria(h.db, "usuarios", usuario.ID, models.AuditoriaUpdate, &before, &usuario, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	return c.JSON(usuario)
}


