package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Origenes de eventos
const (
	OrigenProfesor = "profesor"
	OrigenSistema  = "sistema"
)

// Evento representa la ocurrencia concreta de un concepto
type Evento struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConceptoID    uuid.UUID       `gorm:"type:uuid;not null" json:"concepto_id"`
	Concepto      *Concepto       `gorm:"foreignKey:ConceptoID" json:"concepto,omitempty"`
	AlumnoID      *uuid.UUID      `gorm:"type:uuid" json:"alumno_id,omitempty"`
	Alumno        *Alumno         `gorm:"foreignKey:AlumnoID" json:"alumno,omitempty"`
	CursoID       *uuid.UUID      `gorm:"type:uuid" json:"curso_id,omitempty"`
	Curso         *Curso          `gorm:"foreignKey:CursoID" json:"curso,omitempty"`
	Origen        string          `gorm:"not null" json:"origen"` // profesor, sistema
	OrigenUsuario *uuid.UUID      `gorm:"type:uuid" json:"origen_usuario_id,omitempty"`
	Usuario       *Usuario        `gorm:"foreignKey:OrigenUsuario" json:"usuario,omitempty"`
	Datos         json.RawMessage `gorm:"type:jsonb" json:"datos,omitempty"`
	Activo        bool            `gorm:"default:true" json:"activo"`
	CerradoEn     *time.Time      `json:"cerrado_en,omitempty"`
	CerradoPor    *uuid.UUID      `gorm:"type:uuid" json:"cerrado_por,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (e *Evento) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// Cerrar cierra el evento
func (e *Evento) Cerrar(usuarioID uuid.UUID) {
	now := time.Now()
	e.Activo = false
	e.CerradoEn = &now
	e.CerradoPor = &usuarioID
}
