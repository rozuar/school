package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Conceptos predefinidos
const (
	ConceptoInasistencia   = "INASISTENCIA"
	ConceptoBano           = "BANO"
	ConceptoEnfermeria     = "ENFERMERIA"
	ConceptoSOS            = "SOS"
	ConceptoComportamiento = "COMPORTAMIENTO"
	ConceptoDisciplinario  = "DISCIPLINARIO"
)

// Concepto representa un evento estandarizado del establecimiento
type Concepto struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Codigo      string         `gorm:"uniqueIndex;not null" json:"codigo"` // "INASISTENCIA", "BANO", etc
	Nombre      string         `gorm:"not null" json:"nombre"`
	Descripcion string         `gorm:"type:text" json:"descripcion"`
	Activo      bool           `gorm:"default:true" json:"activo"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (c *Concepto) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
