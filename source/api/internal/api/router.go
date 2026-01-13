package api

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/school-monitoring/backend/internal/api/handlers"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"github.com/school-monitoring/backend/internal/services/orchestrator"
	"github.com/school-monitoring/backend/internal/websocket"
	"gorm.io/gorm"
)

// NewRouter crea y configura el router principal (Fiber).
func NewRouter(db *gorm.DB, hub *websocket.Hub) *fiber.App {
	app := fiber.New()

	// Middleware global
	app.Use(middleware.RequestIDMiddleware)
	app.Use(middleware.CORSMiddleware)
	app.Use(middleware.JSONMiddleware)

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
	api := app.Group("/api/v1")

	// Rutas publicas (sin autenticacion)
	api.Get("/health", handlers.Health)
	api.Post("/auth/login", authHandler.Login)
	api.Post("/seed", seedHandler.Seed)

	// Rutas protegidas
	protected := api.Group("", middleware.AuthMiddleware)

	// Auth
	protected.Post("/auth/logout", authHandler.Logout)
	protected.Post("/auth/refresh", authHandler.RefreshToken)
	protected.Get("/auth/me", authHandler.Me)
	protected.Get("/auth/permisos", authHandler.Permisos)

	// Cursos (todos los roles autenticados pueden ver)
	cursosRoutes := protected.Group("/cursos", middleware.PermissionMiddleware(auth.PermisoVerCursos, auth.PermisoVerAlumnos))
	cursosRoutes.Get("", cursosHandler.GetAll)
	cursosRoutes.Get("/:id", cursosHandler.GetByID)
	cursosRoutes.Get("/:id/alumnos", cursosHandler.GetAlumnos)
	cursosRoutes.Get("/:id/horario", cursosHandler.GetHorario)

	// Horarios del profesor autenticado
	misHorarios := protected.Group("", middleware.PermissionMiddleware(auth.PermisoVerCursos))
	misHorarios.Get("/horarios/mis", horariosHandler.GetMis)

	// Catalogos
	catalogos := protected.Group("", middleware.PermissionMiddleware(auth.PermisoVerCursos))
	catalogos.Get("/asignaturas", asignaturasHandler.GetAll)
	catalogos.Get("/bloques", bloquesHandler.GetAll)

	// Asistencia
	asistenciaRoutes := protected.Group("/asistencia", middleware.PermissionMiddleware(auth.PermisoRegistrarAsistencia, auth.PermisoVerAsistencia))
	asistenciaRoutes.Post("/bloque", asistenciaHandler.RegistrarBloque)
	asistenciaRoutes.Get("/curso/:id/fecha/:fecha", asistenciaHandler.GetByCursoFecha)
	asistenciaRoutes.Get("/horario/:id/fecha/:fecha", asistenciaHandler.GetByHorarioFecha)

	// Estados temporales de alumnos
	estTemp := protected.Group("", middleware.PermissionMiddleware(auth.PermisoRegistrarAsistencia, auth.PermisoCrearEventos, auth.PermisoVerAsistencia))
	estTemp.Put("/alumnos/:id/estado-temporal", asistenciaHandler.SetEstadoTemporal)
	estTemp.Delete("/alumnos/:id/estado-temporal", asistenciaHandler.ClearEstadoTemporal)
	estTemp.Get("/estados-temporales", asistenciaHandler.GetEstadosTemporalesActivos)

	// Conceptos (backoffice)
	conceptosRoutes := protected.Group("/conceptos")
	conceptosRoutes.Get("", conceptosHandler.GetAll)
	conceptosRoutes.Get("/:id", conceptosHandler.GetByID)

	conceptosAdmin := conceptosRoutes.Group("", middleware.RoleMiddleware(models.RolAdmin, models.RolBackoffice))
	conceptosAdmin.Post("", conceptosHandler.Create)
	conceptosAdmin.Put("/:id", conceptosHandler.Update)
	conceptosAdmin.Delete("/:id", conceptosHandler.Delete)

	// Acciones (backoffice)
	accionesRoutes := protected.Group("/acciones")
	accionesRoutes.Get("", accionesHandler.GetAll)
	accionesRoutes.Get("/:id", accionesHandler.GetByID)

	accionesAdmin := accionesRoutes.Group("", middleware.RoleMiddleware(models.RolAdmin, models.RolBackoffice))
	accionesAdmin.Post("", accionesHandler.Create)
	accionesAdmin.Put("/:id", accionesHandler.Update)
	accionesAdmin.Delete("/:id", accionesHandler.Delete)

	// Reglas (backoffice)
	reglasRoutes := protected.Group("/reglas")
	reglasRoutes.Get("", reglasHandler.GetAll)
	reglasRoutes.Get("/:id", reglasHandler.GetByID)

	reglasAdmin := reglasRoutes.Group("", middleware.RoleMiddleware(models.RolAdmin, models.RolBackoffice))
	reglasAdmin.Post("", reglasHandler.Create)
	reglasAdmin.Put("/:id", reglasHandler.Update)
	reglasAdmin.Delete("/:id", reglasHandler.Delete)

	// Eventos
	eventosRoutes := protected.Group("/eventos", middleware.PermissionMiddleware(auth.PermisoVerEventos, auth.PermisoCrearEventos, auth.PermisoCerrarEventos))
	eventosRoutes.Get("", eventosHandler.GetAll)
	eventosRoutes.Get("/activos", eventosHandler.GetActivos)
	eventosRoutes.Get("/alumno/:id", eventosHandler.GetByAlumno)
	eventosRoutes.Post("", eventosHandler.Create)
	eventosRoutes.Put("/:id/cerrar", eventosHandler.Cerrar)

	// Dashboard
	dash := protected.Group("", middleware.PermissionMiddleware(auth.PermisoVerReportes, auth.PermisoVerEventos))
	dash.Get("/dashboard", dashboardHandler.Get)

	// Monitor snapshot (inspectoría/admin/backoffice)
	monitor := protected.Group("/monitor", middleware.PermissionMiddleware(auth.PermisoVerMonitor))
	monitor.Get("/snapshot", monitorHandler.Snapshot)

	// Alertas operativas (inspectoría/admin/backoffice)
	alertas := protected.Group("/alertas", middleware.PermissionMiddleware(auth.PermisoVerAlertas, auth.PermisoCerrarAlertas))
	alertas.Get("", alertasHandler.GetAll)
	alertas.Put("/:id/cerrar", alertasHandler.Cerrar)

	// Admin (usuarios + horarios)
	admin := protected.Group("", middleware.PermissionMiddleware(auth.PermisoAdministrar, auth.PermisoGestionarUsuarios, auth.PermisoGestionarHorarios, auth.PermisoImportarDatos, auth.PermisoVerAuditoria))

	admin.Get("/usuarios", usuariosHandler.GetAll)
	admin.Get("/usuarios/:id", usuariosHandler.GetByID)
	admin.Post("/usuarios", usuariosHandler.Create)
	admin.Put("/usuarios/:id", usuariosHandler.Update)

	admin.Get("/horarios", horariosHandler.GetAll)
	admin.Post("/horarios", horariosHandler.Upsert)
	admin.Delete("/horarios/:id", horariosHandler.Delete)

	// Importaciones
	admin.Post("/import/horarios", importHandler.ImportHorariosCSV)

	// Trazabilidad
	admin.Get("/auditorias", trazabilidadHandler.Auditorias)
	admin.Get("/acciones-ejecuciones", trazabilidadHandler.AccionesEjecuciones)

	// WebSocket
	app.Get("/ws", fiberws.New(func(conn *fiberws.Conn) {
		websocket.HandleConnections(hub, conn)
	}))

	return app
}
