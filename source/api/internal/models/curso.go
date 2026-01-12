package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Niveles educativos
const (
	NivelBasica = "basica"
	NivelMedia  = "media"
)

// Curso representa un curso del establecimiento
type Curso struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nombre    string         `gorm:"not null" json:"nombre"` // "1 Basico", "4 Medio"
	Nivel     string         `gorm:"not null" json:"nivel"`  // basica, media
	Alumnos   []Alumno       `gorm:"foreignKey:CursoID" json:"alumnos,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (c *Curso) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
