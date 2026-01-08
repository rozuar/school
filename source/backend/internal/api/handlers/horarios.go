package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

// GetAll lista horarios con filtros opcionales (curso_id, dia_semana)
func (h *HorariosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").Preload("Curso")

	if cursoID := r.URL.Query().Get("curso_id"); cursoID != "" {
		q = q.Where("curso_id = ?", cursoID)
	}
	if dia := r.URL.Query().Get("dia_semana"); dia != "" {
		q = q.Where("dia_semana = ?", dia)
	}

	var horarios []models.Horario
	if err := q.Order("curso_id, dia_semana, bloque_id").Find(&horarios).Error; err != nil {
		http.Error(w, `{"error":"Error fetching schedules"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(horarios)
}

type HorarioRequest struct {
	CursoID      uuid.UUID `json:"curso_id"`
	AsignaturaID uuid.UUID `json:"asignatura_id"`
	ProfesorID   uuid.UUID `json:"profesor_id"`
	BloqueID     uuid.UUID `json:"bloque_id"`
	DiaSemana    int       `json:"dia_semana"`
}

// Upsert crea o actualiza un horario por (curso_id, dia_semana, bloque_id)
func (h *HorariosHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var req HorarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.CursoID == uuid.Nil || req.AsignaturaID == uuid.Nil || req.ProfesorID == uuid.Nil || req.BloqueID == uuid.Nil || req.DiaSemana < 1 || req.DiaSemana > 5 {
		http.Error(w, `{"error":"curso_id, asignatura_id, profesor_id, bloque_id y dia_semana(1..5) son requeridos"}`, http.StatusBadRequest)
		return
	}

	var existing models.Horario
	err := h.db.Where("curso_id = ? AND dia_semana = ? AND bloque_id = ?", req.CursoID, req.DiaSemana, req.BloqueID).
		First(&existing).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"error":"Error querying schedule"}`, http.StatusInternalServerError)
		return
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
			http.Error(w, `{"error":"Error creating schedule"}`, http.StatusInternalServerError)
			return
		}
		_ = models.CrearAuditoria(h.db, "horarios", hh.ID, models.AuditoriaInsert, nil, &hh, func() *uuid.UUID {
			if claims == nil {
				return nil
			}
			return &claims.UserID
		}())
		h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").Preload("Curso").First(&hh, "id = ?", hh.ID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(hh)
		return
	}

	before := existing
	existing.AsignaturaID = req.AsignaturaID
	existing.ProfesorID = req.ProfesorID
	existing.CursoID = req.CursoID
	existing.BloqueID = req.BloqueID
	existing.DiaSemana = req.DiaSemana

	if err := h.db.Save(&existing).Error; err != nil {
		http.Error(w, `{"error":"Error updating schedule"}`, http.StatusInternalServerError)
		return
	}
	_ = models.CrearAuditoria(h.db, "horarios", existing.ID, models.AuditoriaUpdate, &before, &existing, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	h.db.Preload("Asignatura").Preload("Profesor").Preload("Bloque").Preload("Curso").First(&existing, "id = ?", existing.ID)
	json.NewEncoder(w).Encode(existing)
}

func (h *HorariosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid schedule ID"}`, http.StatusBadRequest)
		return
	}
	var horario models.Horario
	if err := h.db.First(&horario, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error":"Schedule not found"}`, http.StatusNotFound)
		return
	}
	before := horario

	if err := h.db.Delete(&horario).Error; err != nil {
		http.Error(w, `{"error":"Error deleting schedule"}`, http.StatusInternalServerError)
		return
	}
	_ = models.CrearAuditoria(h.db, "horarios", id, models.AuditoriaDelete, &before, nil, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	json.NewEncoder(w).Encode(map[string]string{"message": "Schedule deleted"})
}


