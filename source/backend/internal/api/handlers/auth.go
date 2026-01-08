package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// AuthHandler maneja endpoints de autenticacion
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler crea un nuevo handler de autenticacion
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// LoginRequest estructura de request de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse estructura de response de login
type LoginResponse struct {
	Token   string         `json:"token"`
	Usuario *models.Usuario `json:"usuario"`
}

// Login maneja el login de usuarios
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password required"}`, http.StatusBadRequest)
		return
	}

	// Buscar usuario por email
	var usuario models.Usuario
	if err := h.db.Where("email = ? AND activo = ?", req.Email, true).First(&usuario).Error; err != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Verificar contrasena
	if !auth.CheckPassword(req.Password, usuario.PasswordHash) {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generar token JWT
	token, err := auth.GenerateToken(&usuario)
	if err != nil {
		http.Error(w, `{"error": "Error generating token"}`, http.StatusInternalServerError)
		return
	}

	// Responder con token y usuario
	response := LoginResponse{
		Token:   token,
		Usuario: &usuario,
	}

	json.NewEncoder(w).Encode(response)
}

// Logout maneja el logout (client-side, solo informativo)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// RefreshToken renueva un token JWT
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	newToken, err := auth.RefreshToken(req.Token)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": newToken})
}

// Me retorna el usuario actual
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
		return
	}

	var usuario models.Usuario
	if err := h.db.First(&usuario, "id = ?", claims.UserID).Error; err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(usuario)
}

// Permisos retorna permisos del rol actual y matriz rol->permisos
func (h *AuthHandler) Permisos(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"rol":          claims.Rol,
		"mis_permisos": auth.ObtenerPermisos(claims.Rol),
		"por_rol":      auth.PermisosPorRol(),
	})
}
