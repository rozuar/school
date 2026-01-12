package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Estados de asistencia
const (
	EstadoPresente    = "presente"
	EstadoAusente     = "ausente"
	EstadoJustificado = "justificado"
)

// Tipos de estado temporal
const (
	EstadoTemporalBano       = "bano"
	EstadoTemporalEnfermeria = "enfermeria"
	EstadoTemporalSOS        = "sos"
)

// Asistencia representa el registro de asistencia de un alumno en un bloque
type Asistencia struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AlumnoID      uuid.UUID      `gorm:"type:uuid;not null" json:"alumno_id"`
	Alumno        *Alumno        `gorm:"foreignKey:AlumnoID" json:"alumno,omitempty"`
	HorarioID     uuid.UUID      `gorm:"type:uuid;not null" json:"horario_id"`
	Horario       *Horario       `gorm:"foreignKey:HorarioID" json:"horario,omitempty"`
	Fecha         time.Time      `gorm:"type:date;not null" json:"fecha"`
	Estado        string         `gorm:"not null;default:'ausente'" json:"estado"` // presente, ausente, justificado
	RegistradoPor uuid.UUID      `gorm:"type:uuid" json:"registrado_por"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// EstadoTemporal representa un estado temporal de un alumno (bano, enfermeria, sos)
type EstadoTemporal struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AlumnoID      uuid.UUID      `gorm:"type:uuid;not null" json:"alumno_id"`
	Alumno        *Alumno        `gorm:"foreignKey:AlumnoID" json:"alumno,omitempty"`
	Tipo          string         `gorm:"not null" json:"tipo"` // bano, enfermeria, sos
	Inicio        time.Time      `gorm:"not null" json:"inicio"`
	Fin           *time.Time     `json:"fin,omitempty"`
	RegistradoPor uuid.UUID      `gorm:"type:uuid" json:"registrado_por"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (a *Asistencia) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// BeforeCreate genera UUID antes de crear
func (e *EstadoTemporal) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// EstaActivo verifica si el estado temporal esta activo
func (e *EstadoTemporal) EstaActivo() bool {
	return e.Fin == nil
}
