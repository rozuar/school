package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tipos de acciones
const (
	TipoAccionNotificacion = "notificacion"
	TipoAccionAlerta       = "alerta"
	TipoAccionCambioEstado = "cambio_estado"
	TipoAccionRegistro     = "registro"
)

// Accion representa una respuesta automatica configurable
type Accion struct {
	ID         uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Codigo     string          `gorm:"uniqueIndex;not null" json:"codigo"`
	Nombre     string          `gorm:"not null" json:"nombre"`
	Tipo       string          `gorm:"not null" json:"tipo"` // notificacion, alerta, cambio_estado, registro
	Parametros json.RawMessage `gorm:"type:jsonb" json:"parametros"`
	Activo     bool            `gorm:"default:true" json:"activo"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (a *Accion) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// ParametrosNotificacion estructura para parametros de notificacion
type ParametrosNotificacion struct {
	Destinatario string `json:"destinatario"` // apoderado, inspector, asistente_social
	Asunto       string `json:"asunto"`
	Plantilla    string `json:"plantilla"`
}

// ParametrosAlerta estructura para parametros de alerta
type ParametrosAlerta struct {
	Destinatario string `json:"destinatario"` // inspector, asistente_social
	Prioridad    string `json:"prioridad"`    // baja, media, alta, critica
}
