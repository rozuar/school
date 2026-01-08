package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	NotificacionEstadoPendiente = "pendiente"
	NotificacionEstadoEnviada   = "enviada"
	NotificacionEstadoError     = "error"
)

// NotificationOutbox permite envío asíncrono (push/email/in-app) con trazabilidad.
type NotificationOutbox struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Canal       string          `gorm:"not null;index" json:"canal"` // in_app, push, email
	Destinatario string         `gorm:"not null;index" json:"destinatario"` // rol/email/user_id (simple para MVP)
	Asunto      string          `gorm:"not null" json:"asunto"`
	Payload     json.RawMessage `gorm:"type:jsonb" json:"payload"`
	Estado      string          `gorm:"not null;index;default:'pendiente'" json:"estado"` // pendiente, enviada, error
	Intentos    int             `gorm:"not null;default:0" json:"intentos"`
	SiguienteIntentoEn *time.Time `json:"siguiente_intento_en,omitempty"`
	UltimoError string          `gorm:"type:text" json:"ultimo_error,omitempty"`
	CreadoPor   *uuid.UUID      `gorm:"type:uuid;index" json:"creado_por,omitempty"`
	EnviadoEn   *time.Time      `json:"enviado_en,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (n *NotificationOutbox) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	if n.Estado == "" {
		n.Estado = NotificacionEstadoPendiente
	}
	return nil
}


