package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
func (h *EventosHandler) GetActivos(w http.ResponseWriter, r *http.Request) {
	var eventos []models.Evento
	if err := h.db.Preload("Concepto").Preload("Alumno").Preload("Curso").
		Where("activo = ?", true).
		Order("created_at DESC").
		Find(&eventos).Error; err != nil {
		http.Error(w, `{"error": "Error fetching events"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(eventos)
}

// GetAll obtiene todos los eventos con filtros opcionales
func (h *EventosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := h.db.Preload("Concepto").Preload("Alumno").Preload("Curso")

	// Filtro por alumno
	if alumnoID := r.URL.Query().Get("alumno_id"); alumnoID != "" {
		query = query.Where("alumno_id = ?", alumnoID)
	}

	// Filtro por curso
	if cursoID := r.URL.Query().Get("curso_id"); cursoID != "" {
		query = query.Where("curso_id = ?", cursoID)
	}

	// Filtro por concepto
	if conceptoID := r.URL.Query().Get("concepto_id"); conceptoID != "" {
		query = query.Where("concepto_id = ?", conceptoID)
	}

	// Filtro por activo
	if activo := r.URL.Query().Get("activo"); activo != "" {
		query = query.Where("activo = ?", activo == "true")
	}

	var eventos []models.Evento
	if err := query.Order("created_at DESC").Find(&eventos).Error; err != nil {
		http.Error(w, `{"error": "Error fetching events"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(eventos)
}

// GetByAlumno obtiene eventos de un alumno
func (h *EventosHandler) GetByAlumno(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alumnoID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid student ID"}`, http.StatusBadRequest)
		return
	}

	var eventos []models.Evento
	if err := h.db.Preload("Concepto").
		Where("alumno_id = ?", alumnoID).
		Order("created_at DESC").
		Find(&eventos).Error; err != nil {
		http.Error(w, `{"error": "Error fetching events"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(eventos)
}

// EventoRequest estructura para crear evento
type EventoRequest struct {
	ConceptoID uuid.UUID       `json:"concepto_id"`
	AlumnoID   *uuid.UUID      `json:"alumno_id,omitempty"`
	CursoID    *uuid.UUID      `json:"curso_id,omitempty"`
	Datos      json.RawMessage `json:"datos,omitempty"`
}

// Create crea un nuevo evento
func (h *EventosHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	var req EventoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.ConceptoID == uuid.Nil {
		http.Error(w, `{"error": "Concept ID is required"}`, http.StatusBadRequest)
		return
	}

	// Verificar que el concepto existe
	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", req.ConceptoID).Error; err != nil {
		http.Error(w, `{"error": "Concept not found"}`, http.StatusBadRequest)
		return
	}

	evento := models.Evento{
		ConceptoID:    req.ConceptoID,
		AlumnoID:      req.AlumnoID,
		CursoID:       req.CursoID,
		Origen:        models.OrigenProfesor,
		OrigenUsuario: &claims.UserID,
		Datos:         req.Datos,
		Activo:        true,
	}

	if h.orch != nil {
		if err := h.orch.CreateEvento(&evento, &claims.UserID); err != nil {
			http.Error(w, `{"error": "Error creating event"}`, http.StatusInternalServerError)
			return
		}
	} else if err := h.db.Create(&evento).Error; err != nil {
		http.Error(w, `{"error": "Error creating event"}`, http.StatusInternalServerError)
		return
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Alumno").Preload("Curso").First(&evento, "id = ?", evento.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(evento)
}

// Cerrar cierra un evento
func (h *EventosHandler) Cerrar(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, `{"error": "User not authenticated"}`, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid event ID"}`, http.StatusBadRequest)
		return
	}

	var evento models.Evento
	if err := h.db.First(&evento, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Event not found"}`, http.StatusNotFound)
		return
	}

	if !evento.Activo {
		http.Error(w, `{"error": "Event already closed"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	evento.Activo = false
	evento.CerradoEn = &now
	evento.CerradoPor = &claims.UserID

	if err := h.db.Save(&evento).Error; err != nil {
		http.Error(w, `{"error": "Error closing event"}`, http.StatusInternalServerError)
		return
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Alumno").Preload("Curso").First(&evento, "id = ?", evento.ID)

	json.NewEncoder(w).Encode(evento)
}
