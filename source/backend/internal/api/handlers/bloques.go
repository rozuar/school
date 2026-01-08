package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type BloquesHandler struct {
	db *gorm.DB
}

func NewBloquesHandler(db *gorm.DB) *BloquesHandler {
	return &BloquesHandler{db: db}
}

func (h *BloquesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var bloques []models.BloqueHorario
	if err := h.db.Order("numero").Find(&bloques).Error; err != nil {
		http.Error(w, `{"error":"Error fetching blocks"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bloques)
}


