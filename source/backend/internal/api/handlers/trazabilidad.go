package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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
func (h *TrazabilidadHandler) Auditorias(w http.ResponseWriter, r *http.Request) {
	q := h.db.Preload("Usuario").Model(&models.Auditoria{})

	if tabla := r.URL.Query().Get("tabla"); tabla != "" {
		q = q.Where("tabla = ?", tabla)
	}
	if registroID := r.URL.Query().Get("registro_id"); registroID != "" {
		if id, err := uuid.Parse(registroID); err == nil {
			q = q.Where("registro_id = ?", id)
		}
	}
	if usuarioID := r.URL.Query().Get("usuario_id"); usuarioID != "" {
		if id, err := uuid.Parse(usuarioID); err == nil {
			q = q.Where("usuario_id = ?", id)
		}
	}

	if desde := r.URL.Query().Get("desde"); desde != "" {
		if t, err := time.Parse(time.RFC3339, desde); err == nil {
			q = q.Where("created_at >= ?", t)
		}
	}
	if hasta := r.URL.Query().Get("hasta"); hasta != "" {
		if t, err := time.Parse(time.RFC3339, hasta); err == nil {
			q = q.Where("created_at <= ?", t)
		}
	}

	limit := clamp(atoi(r.URL.Query().Get("limit")), 1, 200)
	offset := clamp(atoi(r.URL.Query().Get("offset")), 0, 1000000)

	var out []models.Auditoria
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		http.Error(w, `{"error":"Error fetching audit logs"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(out)
}

// GET /acciones-ejecuciones?evento_id=&alumno_id=&curso_id=&regla_id=&accion_id=&limit=&offset=
func (h *TrazabilidadHandler) AccionesEjecuciones(w http.ResponseWriter, r *http.Request) {
	q := h.db.Preload("Regla").Preload("Accion").Preload("Evento").Model(&models.AccionEjecucion{})

	if eventoID := r.URL.Query().Get("evento_id"); eventoID != "" {
		if id, err := uuid.Parse(eventoID); err == nil {
			q = q.Where("evento_id = ?", id)
		}
	}
	if alumnoID := r.URL.Query().Get("alumno_id"); alumnoID != "" {
		if id, err := uuid.Parse(alumnoID); err == nil {
			q = q.Where("alumno_id = ?", id)
		}
	}
	if cursoID := r.URL.Query().Get("curso_id"); cursoID != "" {
		if id, err := uuid.Parse(cursoID); err == nil {
			q = q.Where("curso_id = ?", id)
		}
	}
	if reglaID := r.URL.Query().Get("regla_id"); reglaID != "" {
		if id, err := uuid.Parse(reglaID); err == nil {
			q = q.Where("regla_id = ?", id)
		}
	}
	if accionID := r.URL.Query().Get("accion_id"); accionID != "" {
		if id, err := uuid.Parse(accionID); err == nil {
			q = q.Where("accion_id = ?", id)
		}
	}

	if desde := r.URL.Query().Get("desde"); desde != "" {
		if t, err := time.Parse(time.RFC3339, desde); err == nil {
			q = q.Where("ejecutado_en >= ?", t)
		}
	}
	if hasta := r.URL.Query().Get("hasta"); hasta != "" {
		if t, err := time.Parse(time.RFC3339, hasta); err == nil {
			q = q.Where("ejecutado_en <= ?", t)
		}
	}

	limit := clamp(atoi(r.URL.Query().Get("limit")), 1, 200)
	offset := clamp(atoi(r.URL.Query().Get("offset")), 0, 1000000)

	var out []models.AccionEjecucion
	if err := q.Order("ejecutado_en DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		http.Error(w, `{"error":"Error fetching action executions"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(out)
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


