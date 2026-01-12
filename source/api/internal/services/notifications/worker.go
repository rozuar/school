package notifications

import (
	"log"
	"time"

	"github.com/school-monitoring/backend/internal/models"
	"github.com/school-monitoring/backend/internal/websocket"
	"gorm.io/gorm"
)

type Worker struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewWorker(db *gorm.DB, hub *websocket.Hub) *Worker {
	return &Worker{db: db, hub: hub}
}

// Run procesa en background la outbox. En MVP marcamos como enviada (stub).
func (w *Worker) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			w.processOnce()
		}
	}
}

func (w *Worker) processOnce() {
	var items []models.NotificationOutbox
	q := w.db.Where("estado = ?", models.NotificacionEstadoPendiente).
		Order("created_at ASC").
		Limit(50)
	// si hay siguiente intento, respetarlo
	q = q.Where("siguiente_intento_en IS NULL OR siguiente_intento_en <= ?", time.Now())

	if err := q.Find(&items).Error; err != nil {
		log.Printf("notifications: error fetching outbox: %v", err)
		return
	}
	if len(items) == 0 {
		return
	}

	for i := range items {
		item := items[i]
		now := time.Now()
		before := item

		// Stub: en demo marcamos enviada inmediatamente.
		item.Estado = models.NotificacionEstadoEnviada
		item.EnviadoEn = &now
		item.Intentos = item.Intentos + 1
		item.UltimoError = ""

		if err := w.db.Save(&item).Error; err != nil {
			log.Printf("notifications: error saving outbox: %v", err)
			continue
		}
		_ = models.CrearAuditoria(w.db, "notification_outboxes", item.ID, models.AuditoriaUpdate, &before, &item, item.CreadoPor)

		// WS: notificación enviada (para UI in-app)
		if w.hub != nil {
			_ = w.hub.BroadcastMessage(map[string]interface{}{
				"type": "notificacion_enviada",
				"ts":   now,
				"payload": item,
			})
		}
	}
}


