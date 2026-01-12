package orchestrator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/models"
	"github.com/school-monitoring/backend/internal/websocket"
	"gorm.io/gorm"
)

// Orchestrator centraliza: registrar eventos -> evaluar reglas -> ejecutar acciones -> auditar -> broadcast
type Orchestrator struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func New(db *gorm.DB, hub *websocket.Hub) *Orchestrator {
	return &Orchestrator{db: db, hub: hub}
}

type BroadcastEnvelope struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"ts"`
	Payload   interface{} `json:"payload"`
}

func (o *Orchestrator) broadcast(t string, payload interface{}) {
	if o.hub == nil {
		return
	}
	_ = o.hub.BroadcastMessage(BroadcastEnvelope{
		Type:      t,
		Timestamp: time.Now(),
		Payload:   payload,
	})
}

// Notify permite a otros handlers emitir eventos WS sin acoplarse al hub.
func (o *Orchestrator) Notify(t string, payload interface{}) {
	o.broadcast(t, payload)
}

// CreateEventoTx crea un evento dentro de una transaccion ya existente.
func (o *Orchestrator) CreateEventoTx(tx *gorm.DB, evt *models.Evento, usuarioID *uuid.UUID) error {
	if err := tx.Create(evt).Error; err != nil {
		return err
	}

	_ = models.CrearAuditoria(tx, "eventos", evt.ID, models.AuditoriaInsert, nil, evt, usuarioID)

	// Cargar relaciones para payload (best-effort)
	tx.Preload("Concepto").Preload("Alumno").Preload("Curso").First(evt, "id = ?", evt.ID)
	o.broadcast("evento_creado", evt)

	return o.EvaluateAndExecute(tx, evt, usuarioID)
}

// CreateEvento crea un evento, registra auditoria y dispara reglas+acciones.
func (o *Orchestrator) CreateEvento(evt *models.Evento, usuarioID *uuid.UUID) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		return o.CreateEventoTx(tx, evt, usuarioID)
	})
}

// CloseEventoTx cierra el evento dentro de una transaccion ya existente.
func (o *Orchestrator) CloseEventoTx(tx *gorm.DB, eventoID uuid.UUID, usuarioID uuid.UUID) (*models.Evento, error) {
	var out models.Evento
	if err := tx.First(&out, "id = ?", eventoID).Error; err != nil {
		return nil, err
	}
	if !out.Activo {
		return &out, nil
	}
	before := out
	out.Cerrar(usuarioID)
	if err := tx.Save(&out).Error; err != nil {
		return nil, err
	}
	_ = models.CrearAuditoria(tx, "eventos", out.ID, models.AuditoriaUpdate, &before, &out, &usuarioID)
	tx.Preload("Concepto").Preload("Alumno").Preload("Curso").First(&out, "id = ?", out.ID)
	o.broadcast("evento_cerrado", &out)
	return &out, nil
}

// CloseEvento cierra el evento (si esta activo), registra auditoria y broadcast.
func (o *Orchestrator) CloseEvento(eventoID uuid.UUID, usuarioID uuid.UUID) (*models.Evento, error) {
	var out models.Evento
	err := o.db.Transaction(func(tx *gorm.DB) error {
		_, err := o.CloseEventoTx(tx, eventoID, usuarioID)
		if err != nil {
			return err
		}
		return nil
	})
	return &out, err
}

// EvaluateAndExecute evalua reglas activas del concepto del evento y ejecuta acciones si corresponde.
func (o *Orchestrator) EvaluateAndExecute(tx *gorm.DB, evt *models.Evento, usuarioID *uuid.UUID) error {
	var reglas []models.Regla
	if err := tx.Preload("Accion").
		Where("concepto_id = ? AND activo = ?", evt.ConceptoID, true).
		Find(&reglas).Error; err != nil {
		return err
	}

	for _, regla := range reglas {
		ok, scopeKey, winStart, winEnd, detail, err := o.evalRegla(tx, &regla, evt)
		if err != nil {
			// Registrar auditoria de error sin cortar todo
			d, _ := json.Marshal(map[string]string{"error": err.Error()})
			_ = models.CrearAuditoria(tx, "reglas", regla.ID, models.AuditoriaUpdate, nil, json.RawMessage(d), usuarioID)
			continue
		}
		if !ok {
			continue
		}

		// Deduplicacion v2: por regla+accion+scope+ventana (si aplica). Fallback a regla+evento si no hay scope/ventana.
		var count int64
		q := tx.Model(&models.AccionEjecucion{}).
			Where("regla_id = ? AND accion_id = ?", regla.ID, regla.AccionID)
		if scopeKey != "" && winStart != nil && winEnd != nil {
			q = q.Where("scope_key = ? AND ventana_inicio = ? AND ventana_fin = ?", scopeKey, *winStart, *winEnd)
		} else {
			q = q.Where("evento_id = ?", evt.ID)
		}
		q.Count(&count)
		if count > 0 {
			continue
		}

		if err := o.executeAccion(tx, &regla, evt, scopeKey, winStart, winEnd, detail, usuarioID); err != nil {
			continue
		}
	}
	return nil
}

func (o *Orchestrator) evalRegla(tx *gorm.DB, regla *models.Regla, evt *models.Evento) (bool, string, *time.Time, *time.Time, map[string]interface{}, error) {
	cond, err := regla.ParseCondicion()
	if err != nil {
		return false, "", nil, nil, nil, err
	}

	detail := map[string]interface{}{
		"regla": regla.Nombre,
		"condicion": cond,
	}

	// Regla tipo caso_especial: si el alumno es caso especial, inhibir disparo (por defecto)
	if cond.Tipo == "caso_especial" {
		if evt.AlumnoID == nil {
			return false, "", nil, nil, detail, nil
		}
		var alumno models.Alumno
		if err := tx.First(&alumno, "id = ?", *evt.AlumnoID).Error; err != nil {
			return false, "", nil, nil, detail, err
		}
		detail["caso_especial"] = alumno.CasoEspecial
		// Por defecto: si es caso especial => NO disparar
		if alumno.CasoEspecial {
			return false, "", nil, nil, detail, nil
		}
		return true, "alumno:"+alumno.ID.String(), nil, nil, detail, nil
	}

	if cond.Tipo != "cantidad" {
		// Tipos no implementados aun (tiempo, etc.)
		return false, "", nil, nil, detail, nil
	}

	// Scope (default alumno)
	scope := cond.Scope
	if scope == "" {
		scope = "alumno"
	}

	// Determinar concepto a contar (por defecto el de la regla/evento)
	conceptoID := evt.ConceptoID
	if cond.ConceptoCodigo != "" {
		var c models.Concepto
		if err := tx.First(&c, "codigo = ?", cond.ConceptoCodigo).Error; err == nil {
			conceptoID = c.ID
		}
	}

	// Ventana temporal (dias)
	dias := cond.Dias
	if dias <= 0 {
		dias = 1
	}
	since := time.Now().AddDate(0, 0, -dias)
	until := time.Now()
	detail["since"] = since
	detail["until"] = until
	detail["distinct_dias"] = cond.DistinctDias

	// Construir scope key y consulta
	scopeKey := ""
	var n int64

	if scope == "curso" {
		if evt.CursoID == nil {
			return false, "", &since, &until, detail, nil
		}
		scopeKey = "curso:" + evt.CursoID.String()
		if cond.DistinctDias {
			// distinct days de ocurrencias
			type row struct{ D string }
			var rows []row
			if err := tx.Model(&models.Evento{}).
				Select("DATE(created_at) as d").
				Where("concepto_id = ? AND curso_id = ? AND created_at >= ?", conceptoID, *evt.CursoID, since).
				Group("DATE(created_at)").
				Scan(&rows).Error; err != nil {
				return false, scopeKey, &since, &until, detail, err
			}
			n = int64(len(rows))
		} else {
			if err := tx.Model(&models.Evento{}).
				Where("concepto_id = ? AND curso_id = ? AND created_at >= ?", conceptoID, *evt.CursoID, since).
				Count(&n).Error; err != nil {
				return false, scopeKey, &since, &until, detail, err
			}
		}
	} else {
		// alumno
		if evt.AlumnoID == nil {
			return false, "", &since, &until, detail, nil
		}
		scopeKey = "alumno:" + evt.AlumnoID.String()
		if cond.DistinctDias {
			type row struct{ D string }
			var rows []row
			if err := tx.Model(&models.Evento{}).
				Select("DATE(created_at) as d").
				Where("concepto_id = ? AND alumno_id = ? AND created_at >= ?", conceptoID, *evt.AlumnoID, since).
				Group("DATE(created_at)").
				Scan(&rows).Error; err != nil {
				return false, scopeKey, &since, &until, detail, err
			}
			n = int64(len(rows))
		} else {
			if err := tx.Model(&models.Evento{}).
				Where("concepto_id = ? AND alumno_id = ? AND created_at >= ?", conceptoID, *evt.AlumnoID, since).
				Count(&n).Error; err != nil {
				return false, scopeKey, &since, &until, detail, err
			}
		}
	}

	detail["conteo"] = n

	switch cond.Operador {
	case ">=":
		return n >= int64(cond.Valor), scopeKey, &since, &until, detail, nil
	case "<=":
		return n <= int64(cond.Valor), scopeKey, &since, &until, detail, nil
	case "==":
		return n == int64(cond.Valor), scopeKey, &since, &until, detail, nil
	default:
		return false, scopeKey, &since, &until, detail, fmt.Errorf("operador no soportado: %s", cond.Operador)
	}
}

func (o *Orchestrator) executeAccion(tx *gorm.DB, regla *models.Regla, evt *models.Evento, scopeKey string, winStart, winEnd *time.Time, detail map[string]interface{}, usuarioID *uuid.UUID) error {
	if regla.Accion == nil {
		var accion models.Accion
		if err := tx.First(&accion, "id = ?", regla.AccionID).Error; err != nil {
			return err
		}
		regla.Accion = &accion
	}

	detailBytes, _ := json.Marshal(detail)

	exec := models.AccionEjecucion{
		ReglaID:   regla.ID,
		AccionID:  regla.AccionID,
		EventoID:  evt.ID,
		AlumnoID:  evt.AlumnoID,
		CursoID:   evt.CursoID,
		Resultado: "ok",
		Detalle:   detailBytes,
		ScopeKey:  scopeKey,
		VentanaInicio: winStart,
		VentanaFin:    winEnd,
	}
	if err := tx.Create(&exec).Error; err != nil {
		return err
	}

	_ = models.CrearAuditoria(tx, "acciones_ejecuciones", exec.ID, models.AuditoriaInsert, nil, &exec, usuarioID)

	// Ejecutar "side effects" (MVP+): alertas operativas persistentes + broadcast
	if regla.Accion.Tipo == models.TipoAccionAlerta {
		var p models.ParametrosAlerta
		_ = json.Unmarshal(regla.Accion.Parametros, &p)

		prio := p.Prioridad
		if prio == "" {
			prio = "media"
		}

		// Dedup por evento+accion (no crear 2 alertas por el mismo disparo)
		var existing int64
		tx.Model(&models.Alerta{}).Where("evento_id = ? AND accion_id = ? AND estado = ?", evt.ID, regla.AccionID, models.AlertaAbierta).Count(&existing)
		if existing == 0 {
			alerta := models.Alerta{
				Codigo:    regla.Accion.Codigo,
				Titulo:    regla.Accion.Nombre,
				Prioridad: prio,
				Estado:    models.AlertaAbierta,
				CursoID:   evt.CursoID,
				AlumnoID:  evt.AlumnoID,
				EventoID:  &evt.ID,
				ReglaID:   &regla.ID,
				AccionID:  &regla.AccionID,
				CreadoPor: usuarioID,
			}
			_ = tx.Create(&alerta).Error
			_ = models.CrearAuditoria(tx, "alertas", alerta.ID, models.AuditoriaInsert, nil, &alerta, usuarioID)
			o.broadcast("alerta_creada", alerta)
		}
	}

	// Notificaciones: outbox (asíncrono) + broadcast
	if regla.Accion.Tipo == models.TipoAccionNotificacion {
		var p models.ParametrosNotificacion
		_ = json.Unmarshal(regla.Accion.Parametros, &p)

		dest := p.Destinatario
		if dest == "" {
			dest = "inspector"
		}
		asunto := p.Asunto
		if asunto == "" {
			asunto = regla.Accion.Nombre
		}

		payload := map[string]interface{}{
			"accion":  regla.Accion.Codigo,
			"evento":  evt.ID.String(),
			"curso_id": func() string {
				if evt.CursoID == nil {
					return ""
				}
				return evt.CursoID.String()
			}(),
			"alumno_id": func() string {
				if evt.AlumnoID == nil {
					return ""
				}
				return evt.AlumnoID.String()
			}(),
			"detalle": detail,
		}
		payloadBytes, _ := json.Marshal(payload)

		item := models.NotificationOutbox{
			Canal:       "in_app",
			Destinatario: dest,
			Asunto:      asunto,
			Payload:     payloadBytes,
			Estado:      models.NotificacionEstadoPendiente,
			CreadoPor:   usuarioID,
		}
		if err := tx.Create(&item).Error; err == nil {
			_ = models.CrearAuditoria(tx, "notification_outboxes", item.ID, models.AuditoriaInsert, nil, &item, usuarioID)
			o.broadcast("notificacion_creada", item)
		}
	}

	// Broadcast ejecucion (para trazabilidad realtime)
	o.broadcast("accion_ejecutada", map[string]interface{}{
		"regla":  regla,
		"accion": regla.Accion,
		"evento": evt,
		"exec":   exec,
	})

	return nil
}


