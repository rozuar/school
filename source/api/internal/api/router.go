package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/school-monitoring/backend/internal/api/handlers"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"github.com/school-monitoring/backend/internal/services/orchestrator"
	"github.com/school-monitoring/backend/internal/websocket"
	"gorm.io/gorm"
)

// NewRouter crea y configura el router principal
func NewRouter(db *gorm.DB, hub *websocket.Hub) *mux.Router {
	router := mux.NewRouter()

	// Middleware global
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.JSONMiddleware)

	// Inicializar handlers
	orch := orchestrator.New(db, hub)
	authHandler := handlers.NewAuthHandler(db)
	cursosHandler := handlers.NewCursosHandler(db)
	asistenciaHandler := handlers.NewAsistenciaHandler(db, orch)
	conceptosHandler := handlers.NewConceptosHandler(db)
	accionesHandler := handlers.NewAccionesHandler(db)
	reglasHandler := handlers.NewReglasHandler(db)
	eventosHandler := handlers.NewEventosHandler(db, orch)
	seedHandler := handlers.NewSeedHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)
	usuariosHandler := handlers.NewUsuariosHandler(db)
	asignaturasHandler := handlers.NewAsignaturasHandler(db)
	bloquesHandler := handlers.NewBloquesHandler(db)
	horariosHandler := handlers.NewHorariosHandler(db)
	importHandler := handlers.NewImportHandler(db)
	trazabilidadHandler := handlers.NewTrazabilidadHandler(db)
	monitorHandler := handlers.NewMonitorHandler(db)
	alertasHandler := handlers.NewAlertasHandler(db)

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// Rutas publicas (sin autenticacion)
	api.HandleFunc("/health", handlers.Health).Methods("GET", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/seed", seedHandler.Seed).Methods("POST", "OPTIONS")

	// Rutas protegidas
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Auth
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET", "OPTIONS")
	protected.HandleFunc("/auth/permisos", authHandler.Permisos).Methods("GET", "OPTIONS")

	// Cursos (todos los roles autenticados pueden ver)
	cursosRoutes := protected.PathPrefix("/cursos").Subrouter()
	cursosRoutes.Use(middleware.PermissionMiddleware(auth.PermisoVerCursos, auth.PermisoVerAlumnos))
	cursosRoutes.HandleFunc("", cursosHandler.GetAll).Methods("GET", "OPTIONS")
	cursosRoutes.HandleFunc("/{id}", cursosHandler.GetByID).Methods("GET", "OPTIONS")
	cursosRoutes.HandleFunc("/{id}/alumnos", cursosHandler.GetAlumnos).Methods("GET", "OPTIONS")
	cursosRoutes.HandleFunc("/{id}/horario", cursosHandler.GetHorario).Methods("GET", "OPTIONS")

	// Horarios del profesor autenticado
	misHorarios := protected.PathPrefix("").Subrouter()
	misHorarios.Use(middleware.PermissionMiddleware(auth.PermisoVerCursos))
	misHorarios.HandleFunc("/horarios/mis", horariosHandler.GetMis).Methods("GET", "OPTIONS")

	// Catalogos
	catalogos := protected.PathPrefix("").Subrouter()
	catalogos.Use(middleware.PermissionMiddleware(auth.PermisoVerCursos))
	catalogos.HandleFunc("/asignaturas", asignaturasHandler.GetAll).Methods("GET", "OPTIONS")
	catalogos.HandleFunc("/bloques", bloquesHandler.GetAll).Methods("GET", "OPTIONS")

	// Asistencia
	asistenciaRoutes := protected.PathPrefix("/asistencia").Subrouter()
	asistenciaRoutes.Use(middleware.PermissionMiddleware(auth.PermisoRegistrarAsistencia, auth.PermisoVerAsistencia))
	asistenciaRoutes.HandleFunc("/bloque", asistenciaHandler.RegistrarBloque).Methods("POST", "OPTIONS")
	asistenciaRoutes.HandleFunc("/curso/{id}/fecha/{fecha}", asistenciaHandler.GetByCursoFecha).Methods("GET", "OPTIONS")
	asistenciaRoutes.HandleFunc("/horario/{id}/fecha/{fecha}", asistenciaHandler.GetByHorarioFecha).Methods("GET", "OPTIONS")

	// Estados temporales de alumnos
	estTemp := protected.PathPrefix("").Subrouter()
	estTemp.Use(middleware.PermissionMiddleware(auth.PermisoRegistrarAsistencia, auth.PermisoCrearEventos, auth.PermisoVerAsistencia))
	estTemp.HandleFunc("/alumnos/{id}/estado-temporal", asistenciaHandler.SetEstadoTemporal).Methods("PUT", "OPTIONS")
	estTemp.HandleFunc("/alumnos/{id}/estado-temporal", asistenciaHandler.ClearEstadoTemporal).Methods("DELETE", "OPTIONS")
	estTemp.HandleFunc("/estados-temporales", asistenciaHandler.GetEstadosTemporalesActivos).Methods("GET", "OPTIONS")

	// Conceptos (backoffice)
	conceptosRoutes := protected.PathPrefix("/conceptos").Subrouter()
	conceptosRoutes.HandleFunc("", conceptosHandler.GetAll).Methods("GET", "OPTIONS")
	conceptosRoutes.HandleFunc("/{id}", conceptosHandler.GetByID).Methods("GET", "OPTIONS")

	conceptosAdmin := conceptosRoutes.PathPrefix("").Subrouter()
	conceptosAdmin.Use(middleware.RoleMiddleware(models.RolAdmin, models.RolBackoffice))
	conceptosAdmin.HandleFunc("", conceptosHandler.Create).Methods("POST", "OPTIONS")
	conceptosAdmin.HandleFunc("/{id}", conceptosHandler.Update).Methods("PUT", "OPTIONS")
	conceptosAdmin.HandleFunc("/{id}", conceptosHandler.Delete).Methods("DELETE", "OPTIONS")

	// Acciones (backoffice)
	accionesRoutes := protected.PathPrefix("/acciones").Subrouter()
	accionesRoutes.HandleFunc("", accionesHandler.GetAll).Methods("GET", "OPTIONS")
	accionesRoutes.HandleFunc("/{id}", accionesHandler.GetByID).Methods("GET", "OPTIONS")

	accionesAdmin := accionesRoutes.PathPrefix("").Subrouter()
	accionesAdmin.Use(middleware.RoleMiddleware(models.RolAdmin, models.RolBackoffice))
	accionesAdmin.HandleFunc("", accionesHandler.Create).Methods("POST", "OPTIONS")
	accionesAdmin.HandleFunc("/{id}", accionesHandler.Update).Methods("PUT", "OPTIONS")
	accionesAdmin.HandleFunc("/{id}", accionesHandler.Delete).Methods("DELETE", "OPTIONS")

	// Reglas (backoffice)
	reglasRoutes := protected.PathPrefix("/reglas").Subrouter()
	reglasRoutes.HandleFunc("", reglasHandler.GetAll).Methods("GET", "OPTIONS")
	reglasRoutes.HandleFunc("/{id}", reglasHandler.GetByID).Methods("GET", "OPTIONS")

	reglasAdmin := reglasRoutes.PathPrefix("").Subrouter()
	reglasAdmin.Use(middleware.RoleMiddleware(models.RolAdmin, models.RolBackoffice))
	reglasAdmin.HandleFunc("", reglasHandler.Create).Methods("POST", "OPTIONS")
	reglasAdmin.HandleFunc("/{id}", reglasHandler.Update).Methods("PUT", "OPTIONS")
	reglasAdmin.HandleFunc("/{id}", reglasHandler.Delete).Methods("DELETE", "OPTIONS")

	// Eventos
	eventosRoutes := protected.PathPrefix("/eventos").Subrouter()
	eventosRoutes.Use(middleware.PermissionMiddleware(auth.PermisoVerEventos, auth.PermisoCrearEventos, auth.PermisoCerrarEventos))
	eventosRoutes.HandleFunc("", eventosHandler.GetAll).Methods("GET", "OPTIONS")
	eventosRoutes.HandleFunc("/activos", eventosHandler.GetActivos).Methods("GET", "OPTIONS")
	eventosRoutes.HandleFunc("/alumno/{id}", eventosHandler.GetByAlumno).Methods("GET", "OPTIONS")
	eventosRoutes.HandleFunc("", eventosHandler.Create).Methods("POST", "OPTIONS")
	eventosRoutes.HandleFunc("/{id}/cerrar", eventosHandler.Cerrar).Methods("PUT", "OPTIONS")

	// Dashboard
	dash := protected.PathPrefix("").Subrouter()
	dash.Use(middleware.PermissionMiddleware(auth.PermisoVerReportes, auth.PermisoVerEventos))
	dash.HandleFunc("/dashboard", dashboardHandler.Get).Methods("GET", "OPTIONS")

	// Monitor snapshot (inspectoría/admin/backoffice)
	monitor := protected.PathPrefix("/monitor").Subrouter()
	monitor.Use(middleware.PermissionMiddleware(auth.PermisoVerMonitor))
	monitor.HandleFunc("/snapshot", monitorHandler.Snapshot).Methods("GET", "OPTIONS")

	// Alertas operativas (inspectoría/admin/backoffice)
	alertas := protected.PathPrefix("/alertas").Subrouter()
	alertas.Use(middleware.PermissionMiddleware(auth.PermisoVerAlertas, auth.PermisoCerrarAlertas))
	alertas.HandleFunc("", alertasHandler.GetAll).Methods("GET", "OPTIONS")
	alertas.HandleFunc("/{id}/cerrar", alertasHandler.Cerrar).Methods("PUT", "OPTIONS")

	// Admin (usuarios + horarios)
	admin := protected.PathPrefix("").Subrouter()
	admin.Use(middleware.PermissionMiddleware(auth.PermisoAdministrar, auth.PermisoGestionarUsuarios, auth.PermisoGestionarHorarios, auth.PermisoImportarDatos, auth.PermisoVerAuditoria))

	admin.HandleFunc("/usuarios", usuariosHandler.GetAll).Methods("GET", "OPTIONS")
	admin.HandleFunc("/usuarios/{id}", usuariosHandler.GetByID).Methods("GET", "OPTIONS")
	admin.HandleFunc("/usuarios", usuariosHandler.Create).Methods("POST", "OPTIONS")
	admin.HandleFunc("/usuarios/{id}", usuariosHandler.Update).Methods("PUT", "OPTIONS")

	admin.HandleFunc("/horarios", horariosHandler.GetAll).Methods("GET", "OPTIONS")
	admin.HandleFunc("/horarios", horariosHandler.Upsert).Methods("POST", "OPTIONS")
	admin.HandleFunc("/horarios/{id}", horariosHandler.Delete).Methods("DELETE", "OPTIONS")

	// Importaciones
	admin.HandleFunc("/import/horarios", importHandler.ImportHorariosCSV).Methods("POST", "OPTIONS")

	// Trazabilidad
	admin.HandleFunc("/auditorias", trazabilidadHandler.Auditorias).Methods("GET", "OPTIONS")
	admin.HandleFunc("/acciones-ejecuciones", trazabilidadHandler.AccionesEjecuciones).Methods("GET", "OPTIONS")

	// WebSocket
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleConnections(hub, w, r)
	})

	return router
}
