package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// AccionesHandler maneja endpoints de acciones
type AccionesHandler struct {
	db *gorm.DB
}

// NewAccionesHandler crea un nuevo handler de acciones
func NewAccionesHandler(db *gorm.DB) *AccionesHandler {
	return &AccionesHandler{db: db}
}

// GetAll obtiene todas las acciones
func (h *AccionesHandler) GetAll(c *fiber.Ctx) error {
	var acciones []models.Accion
	if err := h.db.Order("codigo").Find(&acciones).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching actions"})
	}
	return c.JSON(acciones)
}

// GetByID obtiene una accion por ID
func (h *AccionesHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action ID"})
	}

	var accion models.Accion
	if err := h.db.First(&accion, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Action not found"})
	}

	return c.JSON(accion)
}

// AccionRequest estructura para crear/actualizar accion
type AccionRequest struct {
	Codigo     string          `json:"codigo"`
	Nombre     string          `json:"nombre"`
	Tipo       string          `json:"tipo"`
	Parametros json.RawMessage `json:"parametros"`
	Activo     *bool           `json:"activo"`
}

// Create crea una nueva accion
func (h *AccionesHandler) Create(c *fiber.Ctx) error {
	var req AccionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Codigo == "" || req.Nombre == "" || req.Tipo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code, name and type are required"})
	}

	accion := models.Accion{
		Codigo:     req.Codigo,
		Nombre:     req.Nombre,
		Tipo:       req.Tipo,
		Parametros: req.Parametros,
		Activo:     true,
	}

	if err := h.db.Create(&accion).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating action"})
	}

	return c.Status(fiber.StatusCreated).JSON(accion)
}

// Update actualiza una accion
func (h *AccionesHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action ID"})
	}

	var accion models.Accion
	if err := h.db.First(&accion, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Action not found"})
	}

	var req AccionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Codigo != "" {
		accion.Codigo = req.Codigo
	}
	if req.Nombre != "" {
		accion.Nombre = req.Nombre
	}
	if req.Tipo != "" {
		accion.Tipo = req.Tipo
	}
	if req.Parametros != nil {
		accion.Parametros = req.Parametros
	}
	if req.Activo != nil {
		accion.Activo = *req.Activo
	}

	if err := h.db.Save(&accion).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating action"})
	}

	return c.JSON(accion)
}

// Delete elimina una accion (soft delete)
func (h *AccionesHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action ID"})
	}

	if err := h.db.Delete(&models.Accion{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting action"})
	}

	return c.JSON(fiber.Map{"message": "Action deleted"})
}
