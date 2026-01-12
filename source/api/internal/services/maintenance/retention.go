package maintenance

import (
	"log"
	"os"
	"time"

	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// RunRetention ejecuta limpieza periódica por retención (MVP).
func RunRetention(db *gorm.DB, stop <-chan struct{}) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	// one-shot startup cleanup
	cleanup(db)

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			cleanup(db)
		}
	}
}

func cleanup(db *gorm.DB) {
	auditDays := envInt("AUDIT_RETENTION_DAYS", 30)
	execDays := envInt("EXEC_RETENTION_DAYS", 30)
	outboxDays := envInt("OUTBOX_RETENTION_DAYS", 30)
	alertDays := envInt("ALERT_RETENTION_DAYS", 90)

	now := time.Now()
	if auditDays > 0 {
		cut := now.AddDate(0, 0, -auditDays)
		if err := db.Where("created_at < ?", cut).Delete(&models.Auditoria{}).Error; err != nil {
			log.Printf("retention: auditorias: %v", err)
		}
	}
	if execDays > 0 {
		cut := now.AddDate(0, 0, -execDays)
		if err := db.Where("created_at < ?", cut).Delete(&models.AccionEjecucion{}).Error; err != nil {
			log.Printf("retention: acciones_ejecuciones: %v", err)
		}
	}
	if outboxDays > 0 {
		cut := now.AddDate(0, 0, -outboxDays)
		if err := db.Where("created_at < ?", cut).Delete(&models.NotificationOutbox{}).Error; err != nil {
			log.Printf("retention: notification_outboxes: %v", err)
		}
	}
	if alertDays > 0 {
		cut := now.AddDate(0, 0, -alertDays)
		// Solo limpiar cerradas antiguas
		if err := db.Where("estado = ? AND updated_at < ?", models.AlertaCerrada, cut).Delete(&models.Alerta{}).Error; err != nil {
			log.Printf("retention: alertas: %v", err)
		}
	}
}

func envInt(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n := 0
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return def
		}
		n = n*10 + int(ch-'0')
	}
	if n <= 0 {
		return def
	}
	return n
}


