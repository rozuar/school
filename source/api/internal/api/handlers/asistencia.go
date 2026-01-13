package handlers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/models"
	"github.com/school-monitoring/backend/internal/services/orchestrator"
	"gorm.io/gorm"
)

// AsistenciaHandler maneja endpoints de asistencia
type AsistenciaHandler struct {
	db   *gorm.DB
	orch *orchestrator.Orchestrator
}

// NewAsistenciaHandler crea un nuevo handler de asistencia
func NewAsistenciaHandler(db *gorm.DB, orch *orchestrator.Orchestrator) *AsistenciaHandler {
	return &AsistenciaHandler{db: db, orch: orch}
}

// RegistroAsistenciaRequest estructura para registrar asistencia de bloque
type RegistroAsistenciaRequest struct {
	HorarioID uuid.UUID            `json:"horario_id"`
	Fecha     string               `json:"fecha"` // YYYY-MM-DD
	Registros []RegistroAlumno     `json:"registros"`
}

// RegistroAlumno estructura para cada alumno
type RegistroAlumno struct {
	AlumnoID uuid.UUID `json:"alumno_id"`
	Estado   string    `json:"estado"` // presente, ausente
}

// RegistrarBloque registra la asistencia de un bloque completo
func (h *AsistenciaHandler) RegistrarBloque(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	var req RegistroAsistenciaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Parsear fecha
	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, use YYYY-MM-DD"})
	}

	// Verificar que el horario existe
	var horario models.Horario
	if err := h.db.First(&horario, "id = ?", req.HorarioID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	// Registrar asistencia para cada alumno
	tx := h.db.Begin()

	// Concepto INASISTENCIA (si hay ausentes, se registran como Evento)
	var conceptoInasistencia models.Concepto
	_ = h.db.First(&conceptoInasistencia, "codigo = ?", models.ConceptoInasistencia).Error

	for _, registro := range req.Registros {
		// Detectar si ya existia asistencia para saber si cambio (para auditoria simple)
		var prev models.Asistencia
		prevFound := tx.Where("alumno_id = ? AND horario_id = ? AND fecha = ?",
			registro.AlumnoID, req.HorarioID, fecha).
			First(&prev).Error == nil

		asistencia := models.Asistencia{
			AlumnoID:      registro.AlumnoID,
			HorarioID:     req.HorarioID,
			Fecha:         fecha,
			Estado:        registro.Estado,
			RegistradoPor: claims.UserID,
		}

		// Upsert: actualizar si ya existe
		if err := tx.Where("alumno_id = ? AND horario_id = ? AND fecha = ?",
			registro.AlumnoID, req.HorarioID, fecha).
			Assign(asistencia).
			FirstOrCreate(&asistencia).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error registering attendance"})
		}

		// Auditoria basica (no guardamos "antes" real si no existia)
		if prevFound {
			_ = models.CrearAuditoria(tx, "asistencias", asistencia.ID, models.AuditoriaUpdate, &prev, &asistencia, &claims.UserID)
		} else {
			_ = models.CrearAuditoria(tx, "asistencias", asistencia.ID, models.AuditoriaInsert, nil, &asistencia, &claims.UserID)
		}

		// Orquestacion: si se marca AUSENTE => crear Evento INASISTENCIA (deduplicando por dia)
		if h.orch != nil && conceptoInasistencia.ID != uuid.Nil && registro.Estado == models.EstadoAusente {
			dayStart := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, time.UTC)
			dayEnd := dayStart.Add(24 * time.Hour)

			var existing int64
			tx.Model(&models.Evento{}).
				Where("concepto_id = ? AND alumno_id = ? AND created_at >= ? AND created_at < ?",
					conceptoInasistencia.ID, registro.AlumnoID, dayStart, dayEnd).
				Count(&existing)

			if existing == 0 {
				datos, _ := json.Marshal(map[string]interface{}{
					"horario_id": req.HorarioID,
					"fecha":      req.Fecha,
				})
				cid := conceptoInasistencia.ID
				evt := models.Evento{
					ConceptoID:    &cid,
					AlumnoID:      &registro.AlumnoID,
					CursoID:       &horario.CursoID,
					Origen:        models.OrigenProfesor,
					OrigenUsuario: &claims.UserID,
					Datos:         datos,
					Activo:        true,
				}
				_ = h.orch.CreateEventoTx(tx, &evt, &claims.UserID)
			}
		}

		// Si el alumno pasa a PRESENTE y antes estaba AUSENTE, cerrar evento INASISTENCIA activo del dia (si existe)
		if h.orch != nil && conceptoInasistencia.ID != uuid.Nil &&
			registro.Estado == models.EstadoPresente && prevFound && prev.Estado == models.EstadoAusente {
			dayStart := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, time.UTC)
			dayEnd := dayStart.Add(24 * time.Hour)

			var eventos []models.Evento
			tx.Where("concepto_id = ? AND alumno_id = ? AND activo = ? AND created_at >= ? AND created_at < ?",
				conceptoInasistencia.ID, registro.AlumnoID, true, dayStart, dayEnd).
				Find(&eventos)
			for _, e := range eventos {
				_, _ = h.orch.CloseEventoTx(tx, e.ID, claims.UserID)
			}
		}
	}
	tx.Commit()

	// WS: presencia del profesor por bloque (para monitor inspectoría)
	if h.orch != nil {
		// Persistir snapshot de presencia por curso (last attendance)
		now := time.Now()
		ce := models.CursoEstado{
			CursoID:             horario.CursoID,
			UltimaAsistenciaEn:  &now,
			UltimaAsistenciaPor: &claims.UserID,
			UltimoHorarioID:     &req.HorarioID,
			UltimoBloqueID:      &horario.BloqueID,
			UltimoDiaSemana:     &horario.DiaSemana,
		}
		// upsert por curso_id
		h.db.Where("curso_id = ?", horario.CursoID).Assign(ce).FirstOrCreate(&ce)

		// Snapshot por bloque+fecha (conteos)
		presentes, ausentes, justificados := 0, 0, 0
		for _, r := range req.Registros {
			switch r.Estado {
			case models.EstadoPresente:
				presentes++
			case models.EstadoJustificado:
				justificados++
			default:
				ausentes++
			}
		}
		hae := models.HorarioAsistenciaEstado{
			HorarioID:              req.HorarioID,
			Fecha:                  fecha,
			CursoID:                horario.CursoID,
			BloqueID:               horario.BloqueID,
			DiaSemana:              horario.DiaSemana,
			ProfesorID:             horario.ProfesorID,
			Presentes:              presentes,
			Ausentes:               ausentes,
			Justificados:           justificados,
			UltimaActualizacionEn:  &now,
			UltimaActualizacionPor: claims.UserID,
		}
		h.db.Where("horario_id = ? AND fecha = ?", req.HorarioID, fecha).Assign(hae).FirstOrCreate(&hae)
		_ = models.CrearAuditoria(h.db, "horarios_asistencia_estado", hae.ID, models.AuditoriaInsert, nil, &hae, &claims.UserID)

		h.orch.Notify("asistencia_bloque_registrada", map[string]interface{}{
			"curso_id":       horario.CursoID.String(),
			"horario_id":     req.HorarioID.String(),
			"registrado_por": claims.UserID.String(),
			"fecha":          req.Fecha,
			"presentes":      presentes,
			"ausentes":       ausentes,
			"justificados":   justificados,
		})
	}

	return c.JSON(fiber.Map{"message": "Attendance registered successfully"})
}

// GetByHorarioFecha obtiene la asistencia de un bloque (horario) en una fecha
func (h *AsistenciaHandler) GetByHorarioFecha(c *fiber.Ctx) error {
	horarioID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid schedule ID"})
	}

	fechaStr := c.Params("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, use YYYY-MM-DD"})
	}

	var asistencias []models.Asistencia
	if err := h.db.Preload("Alumno").
		Where("horario_id = ? AND fecha = ?", horarioID, fecha).
		Find(&asistencias).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching attendance"})
	}

	return c.JSON(asistencias)
}

// GetByCursoFecha obtiene la asistencia de un curso en una fecha
func (h *AsistenciaHandler) GetByCursoFecha(c *fiber.Ctx) error {
	cursoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	fechaStr := c.Params("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, use YYYY-MM-DD"})
	}

	// Obtener horarios del curso
	var horarioIDs []uuid.UUID
	h.db.Model(&models.Horario{}).Where("curso_id = ?", cursoID).Pluck("id", &horarioIDs)

	// Obtener asistencias
	var asistencias []models.Asistencia
	if err := h.db.Preload("Alumno").Preload("Horario").Preload("Horario.Bloque").Preload("Horario.Asignatura").
		Where("horario_id IN ? AND fecha = ?", horarioIDs, fecha).
		Find(&asistencias).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching attendance"})
	}

	return c.JSON(asistencias)
}

// EstadoTemporalRequest estructura para cambiar estado temporal
type EstadoTemporalRequest struct {
	Tipo string `json:"tipo"` // bano, enfermeria, sos
}

// SetEstadoTemporal establece un estado temporal para un alumno
func (h *AsistenciaHandler) SetEstadoTemporal(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	alumnoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	var req EstadoTemporalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validar tipo
	if req.Tipo != models.EstadoTemporalBano &&
	   req.Tipo != models.EstadoTemporalEnfermeria &&
	   req.Tipo != models.EstadoTemporalSOS {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid state type"})
	}

	// Verificar que el alumno existe
	var alumno models.Alumno
	if err := h.db.First(&alumno, "id = ?", alumnoID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	// Cerrar estados temporales activos del mismo alumno
	h.db.Model(&models.EstadoTemporal{}).
		Where("alumno_id = ? AND fin IS NULL", alumnoID).
		Update("fin", time.Now())

	// Crear nuevo estado temporal
	estado := models.EstadoTemporal{
		AlumnoID:      alumnoID,
		Tipo:          req.Tipo,
		Inicio:        time.Now(),
		RegistradoPor: claims.UserID,
	}

	if err := h.db.Create(&estado).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating temporary state"})
	}

	// WS: notificar estado temporal creado
	if h.orch != nil {
		h.orch.Notify("estado_temporal_creado", estado)
	}

	// Orquestacion: reflejar como Evento (BANO/ENFERMERIA/SOS)
	if h.orch != nil {
		var conceptoCodigo string
		switch req.Tipo {
		case models.EstadoTemporalBano:
			conceptoCodigo = models.ConceptoBano
		case models.EstadoTemporalEnfermeria:
			conceptoCodigo = models.ConceptoEnfermeria
		case models.EstadoTemporalSOS:
			conceptoCodigo = models.ConceptoSOS
		}

		if conceptoCodigo != "" {
			var concepto models.Concepto
			if err := h.db.First(&concepto, "codigo = ?", conceptoCodigo).Error; err == nil {
				// Cerrar eventos activos previos del mismo set (bano/enfermeria/sos)
				var conceptoIDs []uuid.UUID
				h.db.Model(&models.Concepto{}).
					Where("codigo IN ?", []string{models.ConceptoBano, models.ConceptoEnfermeria, models.ConceptoSOS}).
					Pluck("id", &conceptoIDs)

				var activos []models.Evento
				h.db.Where("alumno_id = ? AND activo = ? AND concepto_id IN ?", alumnoID, true, conceptoIDs).Find(&activos)
				for _, e := range activos {
					_, _ = h.orch.CloseEvento(e.ID, claims.UserID)
				}

				datos, _ := json.Marshal(map[string]interface{}{
					"estado_temporal_id": estado.ID,
					"tipo":              req.Tipo,
				})
				cid := concepto.ID
				evt := models.Evento{
					ConceptoID:    &cid,
					AlumnoID:      &alumnoID,
					CursoID:       &alumno.CursoID,
					Origen:        models.OrigenProfesor,
					OrigenUsuario: &claims.UserID,
					Datos:         datos,
					Activo:        true,
				}
				_ = h.orch.CreateEvento(&evt, &claims.UserID)
			}
		}
	}

	return c.JSON(estado)
}

// ClearEstadoTemporal cierra el estado temporal de un alumno
func (h *AsistenciaHandler) ClearEstadoTemporal(c *fiber.Ctx) error {
	claims := middleware.GetUserFromContext(c)

	alumnoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	// Cerrar estados temporales activos
	result := h.db.Model(&models.EstadoTemporal{}).
		Where("alumno_id = ? AND fin IS NULL", alumnoID).
		Update("fin", time.Now())

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error clearing temporary state"})
	}

	// WS: notificar estado temporal cerrado (por alumno)
	if h.orch != nil {
		h.orch.Notify("estado_temporal_cerrado", map[string]string{"alumno_id": alumnoID.String()})
	}

	// Cerrar eventos activos asociados (bano/enfermeria/sos)
	if h.orch != nil && claims != nil {
		var conceptoIDs []uuid.UUID
		h.db.Model(&models.Concepto{}).
			Where("codigo IN ?", []string{models.ConceptoBano, models.ConceptoEnfermeria, models.ConceptoSOS}).
			Pluck("id", &conceptoIDs)

		var activos []models.Evento
		h.db.Where("alumno_id = ? AND activo = ? AND concepto_id IN ?", alumnoID, true, conceptoIDs).Find(&activos)
		for _, e := range activos {
			_, _ = h.orch.CloseEvento(e.ID, claims.UserID)
		}
	}

	return c.JSON(fiber.Map{"message": "Temporary state cleared"})
}

// GetEstadosTemporalesActivos obtiene los estados temporales activos
func (h *AsistenciaHandler) GetEstadosTemporalesActivos(c *fiber.Ctx) error {
	var estados []models.EstadoTemporal
	if err := h.db.Preload("Alumno").Preload("Alumno.Curso").
		Where("fin IS NULL").
		Find(&estados).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching temporary states"})
	}

	return c.JSON(estados)
}
