package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Roles del sistema
const (
	RolAdmin          = "admin"
	RolProfesor       = "profesor"
	RolInspector      = "inspector"
	RolAsistenteSocial = "asistente_social"
	RolBackoffice     = "backoffice"
)

// Usuario representa un usuario del sistema
type Usuario struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Nombre       string         `gorm:"not null" json:"nombre"`
	Rol          string         `gorm:"not null" json:"rol"`
	Activo       bool           `gorm:"default:true" json:"activo"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (u *Usuario) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// EsRolValido verifica si el rol es valido
func EsRolValido(rol string) bool {
	switch rol {
	case RolAdmin, RolProfesor, RolInspector, RolAsistenteSocial, RolBackoffice:
		return true
	}
	return false
}
