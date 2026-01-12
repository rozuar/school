package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Regla representa una condicion que dispara una accion
type Regla struct {
	ID         uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nombre     string          `gorm:"not null" json:"nombre"`
	ConceptoID uuid.UUID       `gorm:"type:uuid;not null" json:"concepto_id"`
	Concepto   *Concepto       `gorm:"foreignKey:ConceptoID" json:"concepto,omitempty"`
	Condicion  json.RawMessage `gorm:"type:jsonb;not null" json:"condicion"`
	AccionID   uuid.UUID       `gorm:"type:uuid;not null" json:"accion_id"`
	Accion     *Accion         `gorm:"foreignKey:AccionID" json:"accion,omitempty"`
	Activo     bool            `gorm:"default:true" json:"activo"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

// BeforeCreate genera UUID antes de crear
func (r *Regla) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// CondicionRegla estructura para definir condiciones
type CondicionRegla struct {
	// V1 (compat)
	Tipo     string `json:"tipo"`     // cantidad, tiempo, caso_especial
	Campo    string `json:"campo"`    // inasistencias, eventos
	Operador string `json:"operador"` // >=, <=, ==
	Valor    int    `json:"valor"`    // cantidad
	Dias     int    `json:"dias"`     // periodo en dias

	// V2 (extensiones)
	Scope         string `json:"scope,omitempty"`           // alumno, curso
	ConceptoCodigo string `json:"concepto_codigo,omitempty"` // si se quiere contar un concepto distinto al de la regla
	DistinctDias  bool   `json:"distinct_dias,omitempty"`   // cuenta dias distintos (reincidencia no consecutiva)
}

// ParseCondicion parsea la condicion JSON
func (r *Regla) ParseCondicion() (*CondicionRegla, error) {
	var cond CondicionRegla
	err := json.Unmarshal(r.Condicion, &cond)
	if err != nil {
		return nil, err
	}
	return &cond, nil
}
