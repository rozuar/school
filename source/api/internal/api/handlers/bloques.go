package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type BloquesHandler struct {
	db *gorm.DB
}

func NewBloquesHandler(db *gorm.DB) *BloquesHandler {
	return &BloquesHandler{db: db}
}

func (h *BloquesHandler) GetAll(c *fiber.Ctx) error {
	var bloques []models.BloqueHorario
	if err := h.db.Order("numero").Find(&bloques).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching blocks"})
	}
	return c.JSON(bloques)
}


