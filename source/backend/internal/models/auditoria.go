package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tipos de accion de auditoria
const (
	AuditoriaInsert = "INSERT"
	AuditoriaUpdate = "UPDATE"
	AuditoriaDelete = "DELETE"
)

// Auditoria representa un registro de auditoria
type Auditoria struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Tabla           string          `gorm:"not null" json:"tabla"`
	RegistroID      uuid.UUID       `gorm:"type:uuid;not null" json:"registro_id"`
	Accion          string          `gorm:"not null" json:"accion"` // INSERT, UPDATE, DELETE
	DatosAnteriores json.RawMessage `gorm:"type:jsonb" json:"datos_anteriores,omitempty"`
	DatosNuevos     json.RawMessage `gorm:"type:jsonb" json:"datos_nuevos,omitempty"`
	UsuarioID       *uuid.UUID      `gorm:"type:uuid" json:"usuario_id,omitempty"`
	Usuario         *Usuario        `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// BeforeCreate genera UUID antes de crear
func (a *Auditoria) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// CrearAuditoria crea un registro de auditoria
func CrearAuditoria(db *gorm.DB, tabla string, registroID uuid.UUID, accion string, datosAnteriores, datosNuevos interface{}, usuarioID *uuid.UUID) error {
	var anteriores, nuevos json.RawMessage
	var err error

	if datosAnteriores != nil {
		anteriores, err = json.Marshal(datosAnteriores)
		if err != nil {
			return err
		}
	}

	if datosNuevos != nil {
		nuevos, err = json.Marshal(datosNuevos)
		if err != nil {
			return err
		}
	}

	auditoria := &Auditoria{
		Tabla:           tabla,
		RegistroID:      registroID,
		Accion:          accion,
		DatosAnteriores: anteriores,
		DatosNuevos:     nuevos,
		UsuarioID:       usuarioID,
	}

	return db.Create(auditoria).Error
}
