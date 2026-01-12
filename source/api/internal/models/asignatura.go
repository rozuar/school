package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Asignatura representa una asignatura/materia
type Asignatura struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nombre    string         `gorm:"not null;uniqueIndex" json:"nombre"` // "Matematica", "Lenguaje"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (a *Asignatura) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
