package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type HorariosHandler struct {
	db *gorm.DB
}

func NewHorariosHandler(db *gorm.DB) *HorariosHandler {
	return &HorariosHandler{db: db}
}

// GetMis devuelve el horario del profesor autenticado
func (h *HorariosHandler) GetMis(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	q := h.db.Preload("Asignatura").Preload("Bloque").Preload("Curso")

	// Opcional: filtrar por día
	if dia := c.Query("dia_semana"); dia != "" {
		q = q.Where("dia_semana = ?", dia)
	}

	var horarios []models.Horario
	if err := q.Where("profesor_id = ?", claims.UserID).
		Order("dia_semana, bloque_id").
		Find(&horarios).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching schedules"})
	}

	return c.JSON(horarios)
}

// GetAll lista horarios con filtros opcionales (curso_id, dia_semana)
func (h *HorariosHandler) GetAll(c *fiber.Ctx) error {
	q := h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").Preload("Curso")

	if cursoID := c.Query("curso_id"); cursoID != "" {
		q = q.Where("curso_id = ?", cursoID)
	}
	if dia := c.Query("dia_semana"); dia != "" {
		q = q.Where("dia_semana = ?", dia)
	}

	var horarios []models.Horario
	if err := q.Order("curso_id, dia_semana, bloque_id").Find(&horarios).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching schedules"})
	}
	return c.JSON(horarios)
}

type HorarioRequest struct {
	CursoID      uuid.UUID `json:"curso_id"`
	AsignaturaID uuid.UUID `json:"asignatura_id"`
	ProfesorID   uuid.UUID `json:"profesor_id"`
	BloqueID     uuid.UUID `json:"bloque_id"`
	DiaSemana    int       `json:"dia_semana"`
}

// Upsert crea o actualiza un horario por (curso_id, dia_semana, bloque_id)
func (h *HorariosHandler) Upsert(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)

	var req HorarioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.CursoID == uuid.Nil || req.AsignaturaID == uuid.Nil || req.ProfesorID == uuid.Nil || req.BloqueID == uuid.Nil || req.DiaSemana < 1 || req.DiaSemana > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "curso_id, asignatura_id, profesor_id, bloque_id y dia_semana(1..5) son requeridos"})
	}

	var existing models.Horario
	err := h.db.Where("curso_id = ? AND dia_semana = ? AND bloque_id = ?", req.CursoID, req.DiaSemana, req.BloqueID).
		First(&existing).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying schedule"})
	}

	if err == gorm.ErrRecordNotFound {
		hh := models.Horario{
			CursoID:      req.CursoID,
			AsignaturaID: req.AsignaturaID,
			ProfesorID:   req.ProfesorID,
			BloqueID:     req.BloqueID,
			DiaSemana:    req.DiaSemana,
		}
		if err := h.db.Create(&hh).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating schedule"})
		}
		_ = models.CrearAuditoria(h.db, "horarios", hh.ID, models.AuditoriaInsert, nil, &hh, func() *uuid.UUID {
			if claims == nil {
				return nil
			}
			return &claims.UserID
		}())
		h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").Preload("Curso").First(&hh, "id = ?", hh.ID)
		return c.Status(fiber.StatusCreated).JSON(hh)
	}

	before := existing
	existing.AsignaturaID = req.AsignaturaID
	existing.ProfesorID = req.ProfesorID
	existing.CursoID = req.CursoID
	existing.BloqueID = req.BloqueID
	existing.DiaSemana = req.DiaSemana

	if err := h.db.Save(&existing).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating schedule"})
	}
	_ = models.CrearAuditoria(h.db, "horarios", existing.ID, models.AuditoriaUpdate, &before, &existing, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").Preload("Curso").First(&existing, "id = ?", existing.ID)
	return c.JSON(existing)
}

func (h *HorariosHandler) Delete(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid schedule ID"})
	}
	var horario models.Horario
	if err := h.db.First(&horario, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}
	before := horario

	if err := h.db.Delete(&horario).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting schedule"})
	}
	_ = models.CrearAuditoria(h.db, "horarios", id, models.AuditoriaDelete, &before, nil, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	return c.JSON(fiber.Map{"message": "Schedule deleted"})
}
