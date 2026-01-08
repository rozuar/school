package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CursoEstado mantiene un "snapshot" operativo por curso (para monitor y dashboards)
type CursoEstado struct {
	CursoID              uuid.UUID     `gorm:"type:uuid;primaryKey" json:"curso_id"`
	Curso                *Curso        `gorm:"foreignKey:CursoID" json:"curso,omitempty"`
	UltimaAsistenciaEn   *time.Time    `json:"ultima_asistencia_en,omitempty"`
	UltimaAsistenciaPor  *uuid.UUID    `gorm:"type:uuid" json:"ultima_asistencia_por,omitempty"`
	UltimoHorarioID      *uuid.UUID    `gorm:"type:uuid" json:"ultimo_horario_id,omitempty"`
	UltimoBloqueID       *uuid.UUID    `gorm:"type:uuid" json:"ultimo_bloque_id,omitempty"`
	UltimoDiaSemana      *int          `json:"ultimo_dia_semana,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CursoEstado) TableName() string {
	return "cursos_estado"
}


