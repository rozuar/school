package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccionEjecucion registra una accion ejecutada por una regla (trazabilidad y deduplicacion)
type AccionEjecucion struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ReglaID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"regla_id"`
	Regla     *Regla         `gorm:"foreignKey:ReglaID" json:"regla,omitempty"`
	AccionID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"accion_id"`
	Accion    *Accion        `gorm:"foreignKey:AccionID" json:"accion,omitempty"`
	EventoID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"evento_id"`
	Evento    *Evento        `gorm:"foreignKey:EventoID" json:"evento,omitempty"`
	AlumnoID  *uuid.UUID     `gorm:"type:uuid;index" json:"alumno_id,omitempty"`
	CursoID   *uuid.UUID     `gorm:"type:uuid;index" json:"curso_id,omitempty"`
	Resultado string         `gorm:"not null;default:'ok'" json:"resultado"` // ok, error
	Detalle   json.RawMessage `gorm:"type:jsonb" json:"detalle,omitempty"`
	// Para deduplicación por alcance/ventana (reglas v2)
	ScopeKey    string     `gorm:"index" json:"scope_key,omitempty"`       // ej: "alumno:<uuid>" o "curso:<uuid>"
	VentanaInicio *time.Time `json:"ventana_inicio,omitempty"`
	VentanaFin    *time.Time `json:"ventana_fin,omitempty"`
	EjecutadoEn time.Time    `gorm:"not null" json:"ejecutado_en"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *AccionEjecucion) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.EjecutadoEn.IsZero() {
		a.EjecutadoEn = time.Now()
	}
	return nil
}


