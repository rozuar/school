package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type AlertasHandler struct {
	db *gorm.DB
}

func NewAlertasHandler(db *gorm.DB) *AlertasHandler {
	return &AlertasHandler{db: db}
}

// GET /alertas?estado=abierta&prioridad=&curso_id=&limit=&offset=
func (h *AlertasHandler) GetAll(c *fiber.Ctx) error {
	q := h.db.Model(&models.Alerta{})

	if estado := c.Query("estado"); estado != "" {
		q = q.Where("estado = ?", estado)
	}
	if prioridad := c.Query("prioridad"); prioridad != "" {
		q = q.Where("prioridad = ?", prioridad)
	}
	if cursoID := c.Query("curso_id"); cursoID != "" {
		if id, err := uuid.Parse(cursoID); err == nil {
			q = q.Where("curso_id = ?", id)
		}
	}

	if desde := c.Query("desde"); desde != "" {
		if t, err := time.Parse(time.RFC3339, desde); err == nil {
			q = q.Where("created_at >= ?", t)
		}
	}
	if hasta := c.Query("hasta"); hasta != "" {
		if t, err := time.Parse(time.RFC3339, hasta); err == nil {
			q = q.Where("created_at <= ?", t)
		}
	}

	limit := clamp(atoi(c.Query("limit")), 1, 200)
	offset := clamp(atoi(c.Query("offset")), 0, 1000000)

	var out []models.Alerta
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching alerts"})
	}
	return c.JSON(out)
}

// PUT /alertas/{id}/cerrar
func (h *AlertasHandler) Cerrar(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid alert ID"})
	}

	var alerta models.Alerta
	if err := h.db.First(&alerta, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alert not found"})
	}
	if alerta.Estado == models.AlertaCerrada {
		return c.JSON(alerta)
	}
	before := alerta
	now := time.Now()
	alerta.Estado = models.AlertaCerrada
	alerta.CerradoEn = &now
	if claims != nil {
		alerta.CerradoPor = &claims.UserID
	}

	if err := h.db.Save(&alerta).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error closing alert"})
	}

	_ = models.CrearAuditoria(h.db, "alertas", alerta.ID, models.AuditoriaUpdate, &before, &alerta, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	return c.JSON(alerta)
}


