package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
func (h *AlertasHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := h.db.Model(&models.Alerta{})

	if estado := r.URL.Query().Get("estado"); estado != "" {
		q = q.Where("estado = ?", estado)
	}
	if prioridad := r.URL.Query().Get("prioridad"); prioridad != "" {
		q = q.Where("prioridad = ?", prioridad)
	}
	if cursoID := r.URL.Query().Get("curso_id"); cursoID != "" {
		if id, err := uuid.Parse(cursoID); err == nil {
			q = q.Where("curso_id = ?", id)
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

	var out []models.Alerta
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		http.Error(w, `{"error":"Error fetching alerts"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(out)
}

// PUT /alertas/{id}/cerrar
func (h *AlertasHandler) Cerrar(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid alert ID"}`, http.StatusBadRequest)
		return
	}

	var alerta models.Alerta
	if err := h.db.First(&alerta, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error":"Alert not found"}`, http.StatusNotFound)
		return
	}
	if alerta.Estado == models.AlertaCerrada {
		json.NewEncoder(w).Encode(alerta)
		return
	}
	before := alerta
	now := time.Now()
	alerta.Estado = models.AlertaCerrada
	alerta.CerradoEn = &now
	if claims != nil {
		alerta.CerradoPor = &claims.UserID
	}

	if err := h.db.Save(&alerta).Error; err != nil {
		http.Error(w, `{"error":"Error closing alert"}`, http.StatusInternalServerError)
		return
	}

	_ = models.CrearAuditoria(h.db, "alertas", alerta.ID, models.AuditoriaUpdate, &before, &alerta, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	json.NewEncoder(w).Encode(alerta)
}


