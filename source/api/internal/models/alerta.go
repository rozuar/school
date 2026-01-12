package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Estados de alerta
const (
	AlertaAbierta = "abierta"
	AlertaCerrada = "cerrada"
)

// Alertas operativas: foco en asistencia/soporte, no castigo.
type Alerta struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Codigo    string         `gorm:"not null" json:"codigo"`    // codigo de accion o tipo
	Titulo    string         `gorm:"not null" json:"titulo"`
	Prioridad string         `gorm:"not null" json:"prioridad"` // baja, media, alta, critica
	Estado    string         `gorm:"not null;default:'abierta'" json:"estado"` // abierta, cerrada

	CursoID  *uuid.UUID `gorm:"type:uuid;index" json:"curso_id,omitempty"`
	AlumnoID *uuid.UUID `gorm:"type:uuid;index" json:"alumno_id,omitempty"`
	EventoID *uuid.UUID `gorm:"type:uuid;index" json:"evento_id,omitempty"`
	ReglaID  *uuid.UUID `gorm:"type:uuid;index" json:"regla_id,omitempty"`
	AccionID *uuid.UUID `gorm:"type:uuid;index" json:"accion_id,omitempty"`

	AsignadoA *uuid.UUID `gorm:"type:uuid;index" json:"asignado_a,omitempty"` // usuario inspector (opcional)
	CreadoPor *uuid.UUID `gorm:"type:uuid;index" json:"creado_por,omitempty"`
	CerradoPor *uuid.UUID `gorm:"type:uuid;index" json:"cerrado_por,omitempty"`
	CerradoEn  *time.Time `json:"cerrado_en,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Alerta) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Estado == "" {
		a.Estado = AlertaAbierta
	}
	return nil
}


