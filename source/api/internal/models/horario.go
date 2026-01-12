package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BloqueHorario representa un bloque de tiempo en el horario
type BloqueHorario struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Numero     int            `gorm:"not null" json:"numero"`      // 1, 2, 3, 4, 5
	HoraInicio string         `gorm:"not null" json:"hora_inicio"` // "09:00"
	HoraFin    string         `gorm:"not null" json:"hora_fin"`    // "10:00"
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Horario representa la asignacion de una asignatura a un curso en un bloque
type Horario struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CursoID      uuid.UUID      `gorm:"type:uuid;not null" json:"curso_id"`
	Curso        *Curso         `gorm:"foreignKey:CursoID" json:"curso,omitempty"`
	AsignaturaID uuid.UUID      `gorm:"type:uuid;not null" json:"asignatura_id"`
	Asignatura   *Asignatura    `gorm:"foreignKey:AsignaturaID" json:"asignatura,omitempty"`
	ProfesorID   uuid.UUID      `gorm:"type:uuid;not null" json:"profesor_id"`
	Profesor     *Usuario       `gorm:"foreignKey:ProfesorID" json:"profesor,omitempty"`
	BloqueID     uuid.UUID      `gorm:"type:uuid;not null" json:"bloque_id"`
	Bloque       *BloqueHorario `gorm:"foreignKey:BloqueID" json:"bloque,omitempty"`
	DiaSemana    int            `gorm:"not null" json:"dia_semana"` // 1=lunes, 5=viernes
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (b *BloqueHorario) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// BeforeCreate genera UUID antes de crear
func (h *Horario) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// DiaSemanaTexto retorna el nombre del dia
func DiaSemanaTexto(dia int) string {
	dias := map[int]string{
		1: "Lunes",
		2: "Martes",
		3: "Miercoles",
		4: "Jueves",
		5: "Viernes",
	}
	return dias[dia]
}
