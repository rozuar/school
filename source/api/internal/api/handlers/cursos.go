package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
func (h *CursosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var cursos []models.Curso
	if err := h.db.Order("nivel, nombre").Find(&cursos).Error; err != nil {
		http.Error(w, `{"error": "Error fetching courses"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cursos)
}

// GetByID obtiene un curso por ID
func (h *CursosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid course ID"}`, http.StatusBadRequest)
		return
	}

	var curso models.Curso
	if err := h.db.First(&curso, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Course not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(curso)
}

// GetAlumnos obtiene los alumnos de un curso
func (h *CursosHandler) GetAlumnos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid course ID"}`, http.StatusBadRequest)
		return
	}

	var alumnos []models.Alumno
	if err := h.db.Where("curso_id = ? AND activo = ?", id, true).Order("apellido, nombre").Find(&alumnos).Error; err != nil {
		http.Error(w, `{"error": "Error fetching students"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(alumnos)
}

// GetHorario obtiene el horario de un curso
func (h *CursosHandler) GetHorario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid course ID"}`, http.StatusBadRequest)
		return
	}

	var horarios []models.Horario
	if err := h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").
		Where("curso_id = ?", id).
		Order("dia_semana, bloque_id").
		Find(&horarios).Error; err != nil {
		http.Error(w, `{"error": "Error fetching schedule"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(horarios)
}

// GetHorarioActual obtiene el horario actual del curso (bloque actual)
func (h *CursosHandler) GetHorarioActual(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid course ID"}`, http.StatusBadRequest)
		return
	}

	// Obtener dia de la semana actual (1=lunes, 5=viernes)
	// TODO: Implementar logica para obtener bloque actual basado en hora

	var horario models.Horario
	if err := h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").
		Where("curso_id = ?", id).
		First(&horario).Error; err != nil {
		http.Error(w, `{"error": "No current schedule found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(horario)
}
