package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

type MonitorHandler struct {
	db *gorm.DB
}

func NewMonitorHandler(db *gorm.DB) *MonitorHandler {
	return &MonitorHandler{db: db}
}

type MonitorCurso struct {
	CursoID                 string     `json:"curso_id"`
	Nombre                  string     `json:"nombre"`
	Nivel                   string     `json:"nivel"`
	SalaSemaforo            string     `json:"sala_semaforo"`    // verde, amarillo, rojo, gris
	ProfesorSemaforo        string     `json:"profesor_semaforo"` // verde, amarillo, gris
	EventosActivos          int        `json:"eventos_activos"`
	EstadosTemporalesActivos int       `json:"estados_temporales_activos"`
	UltimaAsistenciaEn      *time.Time `json:"ultima_asistencia_en,omitempty"`
}

type MonitorSnapshot struct {
	UpdatedAt time.Time     `json:"updated_at"`
	Cursos    []MonitorCurso `json:"cursos"`
}

func (h *MonitorHandler) Snapshot(w http.ResponseWriter, r *http.Request) {
	var cursos []models.Curso
	if err := h.db.Order("nivel, nombre").Find(&cursos).Error; err != nil {
		http.Error(w, `{"error":"Error fetching courses"}`, http.StatusInternalServerError)
		return
	}

	// Eventos activos agrupados por curso y severidad
	type evAgg struct {
		CursoID string
		Codigo  string
		Cnt     int
	}
	var evRows []evAgg
	h.db.Table("eventos").
		Select("curso_id as curso_id, conceptos.codigo as codigo, COUNT(*) as cnt").
		Joins("JOIN conceptos ON eventos.concepto_id = conceptos.id").
		Where("eventos.activo = true AND eventos.curso_id IS NOT NULL").
		Group("curso_id, conceptos.codigo").
		Scan(&evRows)

	byCurso := map[string]map[string]int{}
	for _, row := range evRows {
		if _, ok := byCurso[row.CursoID]; !ok {
			byCurso[row.CursoID] = map[string]int{}
		}
		byCurso[row.CursoID][row.Codigo] = row.Cnt
	}

	// Estados temporales activos agrupados por curso
	type stAgg struct {
		CursoID string
		Cnt     int
	}
	var stRows []stAgg
	h.db.Table("estado_temporals").
		Select("alumnos.curso_id as curso_id, COUNT(*) as cnt").
		Joins("JOIN alumnos ON estado_temporals.alumno_id = alumnos.id").
		Where("estado_temporals.fin IS NULL").
		Group("alumnos.curso_id").
		Scan(&stRows)
	stByCurso := map[string]int{}
	for _, row := range stRows {
		stByCurso[row.CursoID] = row.Cnt
	}

	// Estado curso (presencia profesor)
	var estados []models.CursoEstado
	h.db.Find(&estados)
	ceByCurso := map[string]*models.CursoEstado{}
	for i := range estados {
		ce := estados[i]
		id := ce.CursoID.String()
		ceByCurso[id] = &ce
	}

	out := MonitorSnapshot{UpdatedAt: time.Now()}
	for _, c := range cursos {
		cid := c.ID.String()
		// sala semaforo
		codes := byCurso[cid]
		sala := "verde"
		if codes != nil {
			if codes[models.ConceptoSOS] > 0 || codes[models.ConceptoComportamiento] > 0 || codes[models.ConceptoDisciplinario] > 0 {
				sala = "rojo"
			} else if codes[models.ConceptoBano] > 0 || codes[models.ConceptoEnfermeria] > 0 || codes[models.ConceptoInasistencia] > 0 {
				sala = "amarillo"
			}
		}

		// profesor semaforo por ultima asistencia
		prof := "gris"
		var last *time.Time
		if ce := ceByCurso[cid]; ce != nil && ce.UltimaAsistenciaEn != nil {
			last = ce.UltimaAsistenciaEn
			min := time.Since(*ce.UltimaAsistenciaEn).Minutes()
			if min <= 75 {
				prof = "verde"
			} else if min <= 150 {
				prof = "amarillo"
			} else {
				prof = "gris"
			}
		}

		evtCount := 0
		if codes != nil {
			for _, n := range codes {
				evtCount += n
			}
		}

		stCount := stByCurso[cid]

		// Regla de "sin data": si no hay eventos, no hay estados temporales y no hay asistencia registrada => gris
		if evtCount == 0 && stCount == 0 && last == nil {
			sala = "gris"
		}

		out.Cursos = append(out.Cursos, MonitorCurso{
			CursoID:                 cid,
			Nombre:                  c.Nombre,
			Nivel:                   c.Nivel,
			SalaSemaforo:            sala,
			ProfesorSemaforo:        prof,
			EventosActivos:          evtCount,
			EstadosTemporalesActivos: stCount,
			UltimaAsistenciaEn:      last,
		})
	}

	json.NewEncoder(w).Encode(out)
}


