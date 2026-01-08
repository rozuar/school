package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// ConceptosHandler maneja endpoints de conceptos
type ConceptosHandler struct {
	db *gorm.DB
}

// NewConceptosHandler crea un nuevo handler de conceptos
func NewConceptosHandler(db *gorm.DB) *ConceptosHandler {
	return &ConceptosHandler{db: db}
}

// GetAll obtiene todos los conceptos
func (h *ConceptosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var conceptos []models.Concepto
	if err := h.db.Order("codigo").Find(&conceptos).Error; err != nil {
		http.Error(w, `{"error": "Error fetching concepts"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(conceptos)
}

// GetByID obtiene un concepto por ID
func (h *ConceptosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid concept ID"}`, http.StatusBadRequest)
		return
	}

	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Concept not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(concepto)
}

// ConceptoRequest estructura para crear/actualizar concepto
type ConceptoRequest struct {
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Activo      *bool  `json:"activo"`
}

// Create crea un nuevo concepto
func (h *ConceptosHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req ConceptoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Codigo == "" || req.Nombre == "" {
		http.Error(w, `{"error": "Code and name are required"}`, http.StatusBadRequest)
		return
	}

	concepto := models.Concepto{
		Codigo:      req.Codigo,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Activo:      true,
	}

	if err := h.db.Create(&concepto).Error; err != nil {
		http.Error(w, `{"error": "Error creating concept"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(concepto)
}

// Update actualiza un concepto
func (h *ConceptosHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid concept ID"}`, http.StatusBadRequest)
		return
	}

	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Concept not found"}`, http.StatusNotFound)
		return
	}

	var req ConceptoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Codigo != "" {
		concepto.Codigo = req.Codigo
	}
	if req.Nombre != "" {
		concepto.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		concepto.Descripcion = req.Descripcion
	}
	if req.Activo != nil {
		concepto.Activo = *req.Activo
	}

	if err := h.db.Save(&concepto).Error; err != nil {
		http.Error(w, `{"error": "Error updating concept"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(concepto)
}

// Delete elimina un concepto (soft delete)
func (h *ConceptosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid concept ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.Delete(&models.Concepto{}, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Error deleting concept"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Concept deleted"})
}
