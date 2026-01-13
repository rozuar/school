package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// CursosHandler maneja endpoints de cursos
type CursosHandler struct {
	db *gorm.DB
}

// NewCursosHandler crea un nuevo handler de cursos
func NewCursosHandler(db *gorm.DB) *CursosHandler {
	return &CursosHandler{db: db}
}

// GetAll obtiene todos los cursos
func (h *CursosHandler) GetAll(c *fiber.Ctx) error {
	var cursos []models.Curso
	if err := h.db.Order("nivel, nombre").Find(&cursos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching courses"})
	}
	return c.JSON(cursos)
}

// GetByID obtiene un curso por ID
func (h *CursosHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var curso models.Curso
	if err := h.db.First(&curso, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	}

	return c.JSON(curso)
}

// GetAlumnos obtiene los alumnos de un curso
func (h *CursosHandler) GetAlumnos(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var alumnos []models.Alumno
	if err := h.db.Where("curso_id = ? AND activo = ?", id, true).Order("apellido, nombre").Find(&alumnos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching students"})
	}

	return c.JSON(alumnos)
}

// GetHorario obtiene el horario de un curso
func (h *CursosHandler) GetHorario(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var horarios []models.Horario
	if err := h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").
		Where("curso_id = ?", id).
		Order("dia_semana, bloque_id").
		Find(&horarios).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching schedule"})
	}

	return c.JSON(horarios)
}

// GetHorarioActual obtiene el horario actual del curso (bloque actual)
func (h *CursosHandler) GetHorarioActual(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	// Obtener dia de la semana actual (1=lunes, 5=viernes)
	// TODO: Implementar logica para obtener bloque actual basado en hora

	var horario models.Horario
	if err := h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").
		Where("curso_id = ?", id).
		First(&horario).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No current schedule found"})
	}

	return c.JSON(horario)
}
