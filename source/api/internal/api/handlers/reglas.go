package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// ReglasHandler maneja endpoints de reglas
type ReglasHandler struct {
	db *gorm.DB
}

// NewReglasHandler crea un nuevo handler de reglas
func NewReglasHandler(db *gorm.DB) *ReglasHandler {
	return &ReglasHandler{db: db}
}

// GetAll obtiene todas las reglas
func (h *ReglasHandler) GetAll(c *fiber.Ctx) error {
	var reglas []models.Regla
	if err := h.db.Preload("Concepto").Preload("Accion").Order("nombre").Find(&reglas).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching rules"})
	}
	return c.JSON(reglas)
}

// GetByID obtiene una regla por ID
func (h *ReglasHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid rule ID"})
	}

	var regla models.Regla
	if err := h.db.Preload("Concepto").Preload("Accion").First(&regla, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Rule not found"})
	}

	return c.JSON(regla)
}

// ReglaRequest estructura para crear/actualizar regla
type ReglaRequest struct {
	Nombre     string          `json:"nombre"`
	ConceptoID uuid.UUID       `json:"concepto_id"`
	Condicion  json.RawMessage `json:"condicion"`
	AccionID   uuid.UUID       `json:"accion_id"`
	Activo     *bool           `json:"activo"`
}

// Create crea una nueva regla
func (h *ReglasHandler) Create(c *fiber.Ctx) error {
	var req ReglaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Nombre == "" || req.ConceptoID == uuid.Nil || req.AccionID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, concept_id and action_id are required"})
	}

	// Verificar que concepto y accion existen
	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", req.ConceptoID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Concept not found"})
	}

	var accion models.Accion
	if err := h.db.First(&accion, "id = ?", req.AccionID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Action not found"})
	}

	regla := models.Regla{
		Nombre:     req.Nombre,
		ConceptoID: req.ConceptoID,
		Condicion:  req.Condicion,
		AccionID:   req.AccionID,
		Activo:     true,
	}

	if err := h.db.Create(&regla).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating rule"})
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Accion").First(&regla, "id = ?", regla.ID)

	return c.Status(fiber.StatusCreated).JSON(regla)
}

// Update actualiza una regla
func (h *ReglasHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid rule ID"})
	}

	var regla models.Regla
	if err := h.db.First(&regla, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Rule not found"})
	}

	var req ReglaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Nombre != "" {
		regla.Nombre = req.Nombre
	}
	if req.ConceptoID != uuid.Nil {
		regla.ConceptoID = req.ConceptoID
	}
	if req.Condicion != nil {
		regla.Condicion = req.Condicion
	}
	if req.AccionID != uuid.Nil {
		regla.AccionID = req.AccionID
	}
	if req.Activo != nil {
		regla.Activo = *req.Activo
	}

	if err := h.db.Save(&regla).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating rule"})
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Accion").First(&regla, "id = ?", regla.ID)

	return c.JSON(regla)
}

// Delete elimina una regla (soft delete)
func (h *ReglasHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid rule ID"})
	}

	if err := h.db.Delete(&models.Regla{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting rule"})
	}

	return c.JSON(fiber.Map{"message": "Rule deleted"})
}
