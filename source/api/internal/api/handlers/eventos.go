package handlers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/models"
	"github.com/school-monitoring/backend/internal/services/orchestrator"
	"gorm.io/gorm"
)

// EventosHandler maneja endpoints de eventos
type EventosHandler struct {
	db   *gorm.DB
	orch *orchestrator.Orchestrator
}

// NewEventosHandler crea un nuevo handler de eventos
func NewEventosHandler(db *gorm.DB, orch *orchestrator.Orchestrator) *EventosHandler {
	return &EventosHandler{db: db, orch: orch}
}

// GetActivos obtiene todos los eventos activos
func (h *EventosHandler) GetActivos(c *fiber.Ctx) error {
	var eventos []models.Evento
	if err := h.db.Preload("Concepto").Preload("Alumno").Preload("Curso").
		Where("activo = ?", true).
		Order("created_at DESC").
		Find(&eventos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching events"})
	}
	return c.JSON(eventos)
}

// GetAll obtiene todos los eventos con filtros opcionales
func (h *EventosHandler) GetAll(c *fiber.Ctx) error {
	query := h.db.Preload("Concepto").Preload("Alumno").Preload("Curso")

	// Filtro por alumno
	if alumnoID := c.Query("alumno_id"); alumnoID != "" {
		query = query.Where("alumno_id = ?", alumnoID)
	}

	// Filtro por curso
	if cursoID := c.Query("curso_id"); cursoID != "" {
		query = query.Where("curso_id = ?", cursoID)
	}

	// Filtro por concepto
	if conceptoID := c.Query("concepto_id"); conceptoID != "" {
		query = query.Where("concepto_id = ?", conceptoID)
	}

	// Filtro por activo
	if activo := c.Query("activo"); activo != "" {
		query = query.Where("activo = ?", activo == "true")
	}

	var eventos []models.Evento
	if err := query.Order("created_at DESC").Find(&eventos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching events"})
	}
	return c.JSON(eventos)
}

// GetByAlumno obtiene eventos de un alumno
func (h *EventosHandler) GetByAlumno(c *fiber.Ctx) error {
	alumnoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	var eventos []models.Evento
	if err := h.db.Preload("Concepto").
		Where("alumno_id = ?", alumnoID).
		Order("created_at DESC").
		Find(&eventos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching events"})
	}
	return c.JSON(eventos)
}

// EventoRequest estructura para crear evento
type EventoRequest struct {
	ConceptoID uuid.UUID       `json:"concepto_id"`
	AlumnoID   *uuid.UUID      `json:"alumno_id,omitempty"`
	CursoID    *uuid.UUID      `json:"curso_id,omitempty"`
	Datos      json.RawMessage `json:"datos,omitempty"`
}

// Create crea un nuevo evento
func (h *EventosHandler) Create(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	var req EventoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.ConceptoID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Concept ID is required"})
	}

	// Verificar que el concepto existe
	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", req.ConceptoID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Concept not found"})
	}

	evento := models.Evento{
		ConceptoID:    &req.ConceptoID,
		AlumnoID:      req.AlumnoID,
		CursoID:       req.CursoID,
		Origen:        models.OrigenProfesor,
		OrigenUsuario: &claims.UserID,
		Datos:         req.Datos,
		Activo:        true,
	}

	if h.orch != nil {
		if err := h.orch.CreateEvento(&evento, &claims.UserID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating event"})
		}
	} else if err := h.db.Create(&evento).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating event"})
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Alumno").Preload("Curso").First(&evento, "id = ?", evento.ID)

	return c.Status(fiber.StatusCreated).JSON(evento)
}

// Cerrar cierra un evento
func (h *EventosHandler) Cerrar(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	var evento models.Evento
	if err := h.db.First(&evento, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	if !evento.Activo {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Event already closed"})
	}

	now := time.Now()
	evento.Activo = false
	evento.CerradoEn = &now
	evento.CerradoPor = &claims.UserID

	if err := h.db.Save(&evento).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error closing event"})
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Alumno").Preload("Curso").First(&evento, "id = ?", evento.ID)

	return c.JSON(evento)
}
