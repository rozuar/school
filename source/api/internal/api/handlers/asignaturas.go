package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type AsignaturasHandler struct {
	db *gorm.DB
}

func NewAsignaturasHandler(db *gorm.DB) *AsignaturasHandler {
	return &AsignaturasHandler{db: db}
}

func (h *AsignaturasHandler) GetAll(c *fiber.Ctx) error {
	var asignaturas []models.Asignatura
	if err := h.db.Order("nombre").Find(&asignaturas).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching subjects"})
	}
	return c.JSON(asignaturas)
}


