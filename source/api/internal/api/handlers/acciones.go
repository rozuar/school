package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// AccionesHandler maneja endpoints de acciones
type AccionesHandler struct {
	db *gorm.DB
}

// NewAccionesHandler crea un nuevo handler de acciones
func NewAccionesHandler(db *gorm.DB) *AccionesHandler {
	return &AccionesHandler{db: db}
}

// GetAll obtiene todas las acciones
func (h *AccionesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var acciones []models.Accion
	if err := h.db.Order("codigo").Find(&acciones).Error; err != nil {
		http.Error(w, `{"error": "Error fetching actions"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(acciones)
}

// GetByID obtiene una accion por ID
func (h *AccionesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid action ID"}`, http.StatusBadRequest)
		return
	}

	var accion models.Accion
	if err := h.db.First(&accion, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Action not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(accion)
}

// AccionRequest estructura para crear/actualizar accion
type AccionRequest struct {
	Codigo     string          `json:"codigo"`
	Nombre     string          `json:"nombre"`
	Tipo       string          `json:"tipo"`
	Parametros json.RawMessage `json:"parametros"`
	Activo     *bool           `json:"activo"`
}

// Create crea una nueva accion
func (h *AccionesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req AccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Codigo == "" || req.Nombre == "" || req.Tipo == "" {
		http.Error(w, `{"error": "Code, name and type are required"}`, http.StatusBadRequest)
		return
	}

	accion := models.Accion{
		Codigo:     req.Codigo,
		Nombre:     req.Nombre,
		Tipo:       req.Tipo,
		Parametros: req.Parametros,
		Activo:     true,
	}

	if err := h.db.Create(&accion).Error; err != nil {
		http.Error(w, `{"error": "Error creating action"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(accion)
}

// Update actualiza una accion
func (h *AccionesHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid action ID"}`, http.StatusBadRequest)
		return
	}

	var accion models.Accion
	if err := h.db.First(&accion, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Action not found"}`, http.StatusNotFound)
		return
	}

	var req AccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Codigo != "" {
		accion.Codigo = req.Codigo
	}
	if req.Nombre != "" {
		accion.Nombre = req.Nombre
	}
	if req.Tipo != "" {
		accion.Tipo = req.Tipo
	}
	if req.Parametros != nil {
		accion.Parametros = req.Parametros
	}
	if req.Activo != nil {
		accion.Activo = *req.Activo
	}

	if err := h.db.Save(&accion).Error; err != nil {
		http.Error(w, `{"error": "Error updating action"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accion)
}

// Delete elimina una accion (soft delete)
func (h *AccionesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid action ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.Delete(&models.Accion{}, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error": "Error deleting action"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Action deleted"})
}
