package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// DashboardHandler maneja endpoints del dashboard
type DashboardHandler struct {
	db *gorm.DB
}

// NewDashboardHandler crea un nuevo handler de dashboard
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// DashboardResponse estructura de respuesta del dashboard
type DashboardResponse struct {
	TotalCursos          int64                `json:"total_cursos"`
	TotalAlumnos         int64                `json:"total_alumnos"`
	EventosActivos       int64                `json:"eventos_activos"`
	EstadosTemporales    int64                `json:"estados_temporales_activos"`
	AsistenciaHoy        AsistenciaResumen    `json:"asistencia_hoy"`
	EventosPorTipo       []EventoPorTipo      `json:"eventos_por_tipo"`
	UltimosEventos       []models.Evento      `json:"ultimos_eventos"`
	AlumnosEstadoTemp    []models.EstadoTemporal `json:"alumnos_estado_temporal"`
	CursosEstado         []models.CursoEstado `json:"cursos_estado"`
}

// AsistenciaResumen resumen de asistencia
type AsistenciaResumen struct {
	Presentes   int64 `json:"presentes"`
	Ausentes    int64 `json:"ausentes"`
	Justificados int64 `json:"justificados"`
}

// EventoPorTipo conteo de eventos por tipo de concepto
type EventoPorTipo struct {
	Concepto string `json:"concepto"`
	Cantidad int64  `json:"cantidad"`
}

// Get obtiene los datos del dashboard
func (h *DashboardHandler) Get(c *fiber.Ctx) error {
	var response DashboardResponse

	// Total cursos
	h.db.Model(&models.Curso{}).Count(&response.TotalCursos)

	// Total alumnos activos
	h.db.Model(&models.Alumno{}).Where("activo = ?", true).Count(&response.TotalAlumnos)

	// Eventos activos
	h.db.Model(&models.Evento{}).Where("activo = ?", true).Count(&response.EventosActivos)

	// Estados temporales activos
	h.db.Model(&models.EstadoTemporal{}).Where("fin IS NULL").Count(&response.EstadosTemporales)

	// Asistencia de hoy
	hoy := time.Now().Format("2006-01-02")
	h.db.Model(&models.Asistencia{}).Where("fecha = ? AND estado = ?", hoy, models.EstadoPresente).Count(&response.AsistenciaHoy.Presentes)
	h.db.Model(&models.Asistencia{}).Where("fecha = ? AND estado = ?", hoy, models.EstadoAusente).Count(&response.AsistenciaHoy.Ausentes)
	h.db.Model(&models.Asistencia{}).Where("fecha = ? AND estado = ?", hoy, models.EstadoJustificado).Count(&response.AsistenciaHoy.Justificados)

	// Eventos por tipo (ultimos 7 dias)
	hace7Dias := time.Now().AddDate(0, 0, -7)
	rows, err := h.db.Model(&models.Evento{}).
		Select("conceptos.nombre as concepto, COUNT(*) as cantidad").
		Joins("JOIN conceptos ON eventos.concepto_id = conceptos.id").
		Where("eventos.created_at >= ?", hace7Dias).
		Group("conceptos.nombre").
		Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var evt EventoPorTipo
			rows.Scan(&evt.Concepto, &evt.Cantidad)
			response.EventosPorTipo = append(response.EventosPorTipo, evt)
		}
	}

	// Ultimos 10 eventos
	h.db.Preload("Concepto").Preload("Alumno").Preload("Curso").
		Where("activo = ?", true).
		Order("created_at DESC").
		Limit(10).
		Find(&response.UltimosEventos)

	// Alumnos con estado temporal activo
	h.db.Preload("Alumno").Preload("Alumno.Curso").
		Where("fin IS NULL").
		Find(&response.AlumnosEstadoTemp)

	// Snapshot de presencia por curso (última asistencia)
	h.db.Order("updated_at DESC").Find(&response.CursosEstado)

	return c.JSON(response)
}
