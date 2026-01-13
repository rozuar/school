package database

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize inicializa la conexion a la base de datos
func Initialize() (*gorm.DB, error) {
	var dsn string

	// Prioridad 1: Usar DATABASE_URL si esta disponible (Railway, Heroku, etc.)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		// Parsear la URL de conexion
		parsedURL, err := url.Parse(databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
		}

		// Extraer componentes de la URL
		user := parsedURL.User.Username()
		password, _ := parsedURL.User.Password()
		host := parsedURL.Hostname()
		port := parsedURL.Port()
		if port == "" {
			port = "5432"
		}
		dbname := strings.TrimPrefix(parsedURL.Path, "/")

		// Construir DSN para GORM
		// Railway y otros servicios cloud requieren SSL
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			host, user, password, dbname, port)

		log.Printf("Using DATABASE_URL connection to %s:%s/%s", host, port, dbname)
	} else {
		// Prioridad 2: Usar variables individuales
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}

		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}

		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "postgres"
		}

		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "school_monitoring"
		}

		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		// Para conexiones locales, SSL puede estar deshabilitado
		sslMode := os.Getenv("DB_SSLMODE")
		if sslMode == "" {
			sslMode = "disable"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			host, user, password, dbname, port, sslMode)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")

	// Reset database if DB_RESET=true (development only)
	if os.Getenv("DB_RESET") == "true" {
		log.Println("DB_RESET=true: Dropping all tables...")
		DB.Exec("DROP TABLE IF EXISTS auditorias CASCADE")
		DB.Exec("DROP TABLE IF EXISTS notification_outboxes CASCADE")
		DB.Exec("DROP TABLE IF EXISTS alertas CASCADE")
		DB.Exec("DROP TABLE IF EXISTS horarios_asistencia_estado CASCADE")
		DB.Exec("DROP TABLE IF EXISTS cursos_estado CASCADE")
		DB.Exec("DROP TABLE IF EXISTS eventos CASCADE")
		DB.Exec("DROP TABLE IF EXISTS reglas CASCADE")
		DB.Exec("DROP TABLE IF EXISTS acciones CASCADE")
		DB.Exec("DROP TABLE IF EXISTS conceptos CASCADE")
		DB.Exec("DROP TABLE IF EXISTS estados_temporales CASCADE")
		DB.Exec("DROP TABLE IF EXISTS asistencias CASCADE")
		DB.Exec("DROP TABLE IF EXISTS horarios CASCADE")
		DB.Exec("DROP TABLE IF EXISTS bloque_horarios CASCADE")
		DB.Exec("DROP TABLE IF EXISTS asignaturas CASCADE")
		DB.Exec("DROP TABLE IF EXISTS alumnos CASCADE")
		DB.Exec("DROP TABLE IF EXISTS cursos CASCADE")
		DB.Exec("DROP TABLE IF EXISTS usuarios CASCADE")
		log.Println("Tables dropped successfully")
	}

	// Auto-migrate models (por defecto solo en local, para no bloquear arranque en Railway)
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	autoMigrate := appEnv == "local"
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("AUTO_MIGRATE"))); v != "" {
		autoMigrate = v == "1" || v == "true" || v == "yes" || v == "on"
	}

	if autoMigrate {
		err = DB.AutoMigrate(
			&models.Usuario{},
			&models.Curso{},
			&models.CursoEstado{},
			&models.Alumno{},
			&models.Asignatura{},
			&models.BloqueHorario{},
			&models.Horario{},
			&models.Asistencia{},
			&models.HorarioAsistenciaEstado{},
			&models.EstadoTemporal{},
			&models.Concepto{},
			&models.Accion{},
			&models.Regla{},
			&models.Evento{},
			&models.AccionEjecucion{},
			&models.Alerta{},
			&models.NotificationOutbox{},
			&models.Auditoria{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
		log.Println("Database migration completed")
	} else {
		log.Println("AUTO_MIGRATE disabled: skipping DB migrations")
	}

	return DB, nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}
