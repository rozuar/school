package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// ReglasHandler maneja endpoints de reglas
type ReglasHandler struct {
	db *gorm.DB
}

// NewReglasHandler crea un nuevo handler de reglas
func NewReglasHandler(db *gorm.DB) *ReglasHandler {
	return &ReglasHandler{db: db}
}

// GetAll obtiene todas las reglas
func (h *ReglasHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var reglas []models.Regla
	if err := h.db.Preload("Concepto").Preload("Accion").Order("nombre").Find(&reglas).Error; err != nil {
		http.Error(w, `{"error": "Error fetching rules"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reglas)
}

// GetByID obtiene una regla por ID
func (h *ReglasHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid rule ID"}`, http.StatusBadRequest)
		return
	}

	var regla models.Regla
	if err := h.db.Preload("Concepto").Preload("Accion").First(&regla, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Rule not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(regla)
}

// ReglaRequest estructura para crear/actualizar regla
type ReglaRequest struct {
	Nombre     string          `json:"nombre"`
	ConceptoID uuid.UUID       `json:"concepto_id"`
	Condicion  json.RawMessage `json:"condicion"`
	AccionID   uuid.UUID       `json:"accion_id"`
	Activo     *bool           `json:"activo"`
}

// Create crea una nueva regla
func (h *ReglasHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req ReglaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.ConceptoID == uuid.Nil || req.AccionID == uuid.Nil {
		http.Error(w, `{"error": "Name, concept_id and action_id are required"}`, http.StatusBadRequest)
		return
	}

	// Verificar que concepto y accion existen
	var concepto models.Concepto
	if err := h.db.First(&concepto, "id = ?", req.ConceptoID).Error; err != nil {
		http.Error(w, `{"error": "Concept not found"}`, http.StatusBadRequest)
		return
	}

	var accion models.Accion
	if err := h.db.First(&accion, "id = ?", req.AccionID).Error; err != nil {
		http.Error(w, `{"error": "Action not found"}`, http.StatusBadRequest)
		return
	}

	regla := models.Regla{
		Nombre:     req.Nombre,
		ConceptoID: req.ConceptoID,
		Condicion:  req.Condicion,
		AccionID:   req.AccionID,
		Activo:     true,
	}

	if err := h.db.Create(&regla).Error; err != nil {
		http.Error(w, `{"error": "Error creating rule"}`, http.StatusInternalServerError)
		return
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Accion").First(&regla, "id = ?", regla.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(regla)
}

// Update actualiza una regla
func (h *ReglasHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid rule ID"}`, http.StatusBadRequest)
		return
	}

	var regla models.Regla
	if err := h.db.First(&regla, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Rule not found"}`, http.StatusNotFound)
		return
	}

	var req ReglaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Nombre != "" {
		regla.Nombre = req.Nombre
	}
	if req.ConceptoID != uuid.Nil {
		regla.ConceptoID = req.ConceptoID
	}
	if req.Condicion != nil {
		regla.Condicion = req.Condicion
	}
	if req.AccionID != uuid.Nil {
		regla.AccionID = req.AccionID
	}
	if req.Activo != nil {
		regla.Activo = *req.Activo
	}

	if err := h.db.Save(&regla).Error; err != nil {
		http.Error(w, `{"error": "Error updating rule"}`, http.StatusInternalServerError)
		return
	}

	// Cargar relaciones
	h.db.Preload("Concepto").Preload("Accion").First(&regla, "id = ?", regla.ID)

	json.NewEncoder(w).Encode(regla)
}

// Delete elimina una regla (soft delete)
func (h *ReglasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid rule ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.Delete(&models.Regla{}, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Error deleting rule"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Rule deleted"})
}
