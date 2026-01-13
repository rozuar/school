package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type TrazabilidadHandler struct {
	db *gorm.DB
}

func NewTrazabilidadHandler(db *gorm.DB) *TrazabilidadHandler {
	return &TrazabilidadHandler{db: db}
}

// GET /auditorias?tabla=&registro_id=&usuario_id=&limit=&offset=
func (h *TrazabilidadHandler) Auditorias(c *fiber.Ctx) error {
	q := h.db.Preload("Usuario").Model(&models.Auditoria{})

	if tabla := c.Query("tabla"); tabla != "" {
		q = q.Where("tabla = ?", tabla)
	}
	if registroID := c.Query("registro_id"); registroID != "" {
		if id, err := uuid.Parse(registroID); err == nil {
			q = q.Where("registro_id = ?", id)
		}
	}
	if usuarioID := c.Query("usuario_id"); usuarioID != "" {
		if id, err := uuid.Parse(usuarioID); err == nil {
			q = q.Where("usuario_id = ?", id)
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

	var out []models.Auditoria
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching audit logs"})
	}
	return c.JSON(out)
}

// GET /acciones-ejecuciones?evento_id=&alumno_id=&curso_id=&regla_id=&accion_id=&limit=&offset=
func (h *TrazabilidadHandler) AccionesEjecuciones(c *fiber.Ctx) error {
	q := h.db.Preload("Regla").Preload("Accion").Preload("Evento").Model(&models.AccionEjecucion{})

	if eventoID := c.Query("evento_id"); eventoID != "" {
		if id, err := uuid.Parse(eventoID); err == nil {
			q = q.Where("evento_id = ?", id)
		}
	}
	if alumnoID := c.Query("alumno_id"); alumnoID != "" {
		if id, err := uuid.Parse(alumnoID); err == nil {
			q = q.Where("alumno_id = ?", id)
		}
	}
	if cursoID := c.Query("curso_id"); cursoID != "" {
		if id, err := uuid.Parse(cursoID); err == nil {
			q = q.Where("curso_id = ?", id)
		}
	}
	if reglaID := c.Query("regla_id"); reglaID != "" {
		if id, err := uuid.Parse(reglaID); err == nil {
			q = q.Where("regla_id = ?", id)
		}
	}
	if accionID := c.Query("accion_id"); accionID != "" {
		if id, err := uuid.Parse(accionID); err == nil {
			q = q.Where("accion_id = ?", id)
		}
	}

	if desde := c.Query("desde"); desde != "" {
		if t, err := time.Parse(time.RFC3339, desde); err == nil {
			q = q.Where("ejecutado_en >= ?", t)
		}
	}
	if hasta := c.Query("hasta"); hasta != "" {
		if t, err := time.Parse(time.RFC3339, hasta); err == nil {
			q = q.Where("ejecutado_en <= ?", t)
		}
	}

	limit := clamp(atoi(c.Query("limit")), 1, 200)
	offset := clamp(atoi(c.Query("offset")), 0, 1000000)

	var out []models.AccionEjecucion
	if err := q.Order("ejecutado_en DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching action executions"})
	}
	return c.JSON(out)
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}


