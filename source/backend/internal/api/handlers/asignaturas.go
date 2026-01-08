package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type AsignaturasHandler struct {
	db *gorm.DB
}

func NewAsignaturasHandler(db *gorm.DB) *AsignaturasHandler {
	return &AsignaturasHandler{db: db}
}

func (h *AsignaturasHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var asignaturas []models.Asignatura
	if err := h.db.Order("nombre").Find(&asignaturas).Error; err != nil {
		http.Error(w, `{"error":"Error fetching subjects"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(asignaturas)
}


