package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// ConceptosHandler maneja endpoints de conceptos
type ConceptosHandler struct {
	db *gorm.DB
}

// NewConceptosHandler crea un nuevo handler de conceptos
func NewConceptosHandler(db *gorm.DB) *ConceptosHandler {
	return &ConceptosHandler{db: db}
}

// GetAll obtiene todos los conceptos
func (h *ConceptosHandler) GetAll(c *fiber.Ctx) error {
	var conceptos []models.Concepto
	if err := h.db.Order("codigo").Find(&conceptos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching concepts"})
	}
	return c.JSON(conceptos)
}

// GetByID obtiene un concepto por ID
func (h *ConceptosHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid concept ID"})
	}

	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Concept not found"})
	}

	return c.JSON(concepto)
}

// ConceptoRequest estructura para crear/actualizar concepto
type ConceptoRequest struct {
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Activo      *bool  `json:"activo"`
}

// Create crea un nuevo concepto
func (h *ConceptosHandler) Create(c *fiber.Ctx) error {
	var req ConceptoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Codigo == "" || req.Nombre == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code and name are required"})
	}

	concepto := models.Concepto{
		Codigo:      req.Codigo,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Activo:      true,
	}

	if err := h.db.Create(&concepto).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating concept"})
	}

	return c.Status(fiber.StatusCreated).JSON(concepto)
}

// Update actualiza un concepto
func (h *ConceptosHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid concept ID"})
	}

	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Concept not found"})
	}

	var req ConceptoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Codigo != "" {
		concepto.Codigo = req.Codigo
	}
	if req.Nombre != "" {
		concepto.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		concepto.Descripcion = req.Descripcion
	}
	if req.Activo != nil {
		concepto.Activo = *req.Activo
	}

	if err := h.db.Save(&concepto).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating concept"})
	}

	return c.JSON(concepto)
}

// Delete elimina un concepto (soft delete)
func (h *ConceptosHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid concept ID"})
	}

	if err := h.db.Delete(&models.Concepto{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting concept"})
	}

	return c.JSON(fiber.Map{"message": "Concept deleted"})
}
