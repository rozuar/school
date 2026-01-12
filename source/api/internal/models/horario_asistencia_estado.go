package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HorarioAsistenciaEstado es un snapshot por bloque (horario) y fecha.
// Permite saber rápidamente si el profesor ya registró asistencia en ese bloque,
// y entregar conteos sin recalcular sobre todas las filas.
type HorarioAsistenciaEstado struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	HorarioID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"horario_id"`
	Horario       *Horario       `gorm:"foreignKey:HorarioID" json:"horario,omitempty"`
	Fecha         time.Time      `gorm:"type:date;not null;index" json:"fecha"`

	CursoID    uuid.UUID `gorm:"type:uuid;not null;index" json:"curso_id"`
	BloqueID   uuid.UUID `gorm:"type:uuid;not null;index" json:"bloque_id"`
	DiaSemana  int       `gorm:"not null" json:"dia_semana"`
	ProfesorID uuid.UUID `gorm:"type:uuid;not null;index" json:"profesor_id"`

	Presentes    int `gorm:"not null;default:0" json:"presentes"`
	Ausentes     int `gorm:"not null;default:0" json:"ausentes"`
	Justificados int `gorm:"not null;default:0" json:"justificados"`

	UltimaActualizacionEn *time.Time `json:"ultima_actualizacion_en,omitempty"`
	UltimaActualizacionPor uuid.UUID `gorm:"type:uuid" json:"ultima_actualizacion_por"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HorarioAsistenciaEstado) TableName() string {
	return "horarios_asistencia_estado"
}

func (h *HorarioAsistenciaEstado) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}


