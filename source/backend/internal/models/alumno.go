package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Alumno representa un alumno del establecimiento
type Alumno struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CursoID      uuid.UUID      `gorm:"type:uuid;not null" json:"curso_id"`
	Curso        *Curso         `gorm:"foreignKey:CursoID" json:"curso,omitempty"`
	Nombre       string         `gorm:"not null" json:"nombre"`
	Apellido     string         `gorm:"not null" json:"apellido"`
	Rut          string         `gorm:"uniqueIndex" json:"rut"`
	CasoEspecial bool           `gorm:"default:false" json:"caso_especial"`
	Activo       bool           `gorm:"default:true" json:"activo"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (a *Alumno) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// NombreCompleto retorna el nombre completo del alumno
func (a *Alumno) NombreCompleto() string {
	return a.Nombre + " " + a.Apellido
}
