package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// UsuariosHandler maneja endpoints de administracion de usuarios
type UsuariosHandler struct {
	db *gorm.DB
}

func NewUsuariosHandler(db *gorm.DB) *UsuariosHandler {
	return &UsuariosHandler{db: db}
}

// GetAll obtiene usuarios, con filtros opcionales (rol, activo)
func (h *UsuariosHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := h.db.Model(&models.Usuario{})

	if rol := r.URL.Query().Get("rol"); rol != "" {
		query = query.Where("rol = ?", rol)
	}
	if activo := r.URL.Query().Get("activo"); activo != "" {
		query = query.Where("activo = ?", activo == "true")
	}

	var usuarios []models.Usuario
	if err := query.Order("rol, nombre").Find(&usuarios).Error; err != nil {
		http.Error(w, `{"error":"Error fetching users"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(usuarios)
}

func (h *UsuariosHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := h.db.First(&usuario, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(usuario)
}

type UsuarioRequest struct {
	Email    string  `json:"email"`
	Nombre   string  `json:"nombre"`
	Rol      string  `json:"rol"`
	Password *string `json:"password,omitempty"`
	Activo   *bool   `json:"activo,omitempty"`
}

func (h *UsuariosHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var req UsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Nombre == "" || req.Rol == "" {
		http.Error(w, `{"error":"email, nombre y rol son requeridos"}`, http.StatusBadRequest)
		return
	}
	if !models.EsRolValido(req.Rol) {
		http.Error(w, `{"error":"rol invalido"}`, http.StatusBadRequest)
		return
	}

	pass := "changeme123"
	if req.Password != nil && *req.Password != "" {
		pass = *req.Password
	}
	hash, err := auth.HashPassword(pass)
	if err != nil {
		http.Error(w, `{"error":"Error hashing password"}`, http.StatusInternalServerError)
		return
	}

	usuario := models.Usuario{
		Email:        req.Email,
		Nombre:       req.Nombre,
		Rol:          req.Rol,
		PasswordHash: hash,
		Activo:       true,
	}
	if req.Activo != nil {
		usuario.Activo = *req.Activo
	}

	if err := h.db.Create(&usuario).Error; err != nil {
		http.Error(w, `{"error":"Error creating user"}`, http.StatusInternalServerError)
		return
	}
	_ = models.CrearAuditoria(h.db, "usuarios", usuario.ID, models.AuditoriaInsert, nil, &usuario, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

func (h *UsuariosHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := h.db.First(&usuario, "id = ?", id).Error; err != nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}
	before := usuario

	var req UsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Email != "" {
		usuario.Email = req.Email
	}
	if req.Nombre != "" {
		usuario.Nombre = req.Nombre
	}
	if req.Rol != "" {
		if !models.EsRolValido(req.Rol) {
			http.Error(w, `{"error":"rol invalido"}`, http.StatusBadRequest)
			return
		}
		usuario.Rol = req.Rol
	}
	if req.Activo != nil {
		usuario.Activo = *req.Activo
	}
	if req.Password != nil && *req.Password != "" {
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			http.Error(w, `{"error":"Error hashing password"}`, http.StatusInternalServerError)
			return
		}
		usuario.PasswordHash = hash
	}

	if err := h.db.Save(&usuario).Error; err != nil {
		http.Error(w, `{"error":"Error updating user"}`, http.StatusInternalServerError)
		return
	}

	_ = models.CrearAuditoria(h.db, "usuarios", usuario.ID, models.AuditoriaUpdate, &before, &usuario, func() *uuid.UUID {
		if claims == nil {
			return nil
		}
		return &claims.UserID
	}())

	json.NewEncoder(w).Encode(usuario)
}


