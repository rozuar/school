package auth

import (
	"github.com/school-monitoring/backend/internal/models"
)

// Permisos del sistema
const (
	PermisoVerCursos          = "ver_cursos"
	PermisoVerAlumnos         = "ver_alumnos"
	PermisoRegistrarAsistencia = "registrar_asistencia"
	PermisoVerAsistencia      = "ver_asistencia"
	PermisoCrearEventos       = "crear_eventos"
	PermisoVerEventos         = "ver_eventos"
	PermisoCerrarEventos      = "cerrar_eventos"
	PermisoVerAlertas         = "ver_alertas"
	PermisoCerrarAlertas      = "cerrar_alertas"
	PermisoVerMonitor         = "ver_monitor"
	PermisoVerAuditoria       = "ver_auditoria"
	PermisoGestionarUsuarios  = "gestionar_usuarios"
	PermisoGestionarHorarios  = "gestionar_horarios"
	PermisoImportarDatos      = "importar_datos"
	PermisoVerCasos           = "ver_casos"
	PermisoGestionarCasos     = "gestionar_casos"
	PermisoGestionarConceptos = "gestionar_conceptos"
	PermisoGestionarAcciones  = "gestionar_acciones"
	PermisoGestionarReglas    = "gestionar_reglas"
	PermisoVerReportes        = "ver_reportes"
	PermisoAdministrar        = "administrar"
)

// permisosPorRol define los permisos de cada rol
var permisosPorRol = map[string][]string{
	models.RolAdmin: {
		PermisoVerCursos,
		PermisoVerAlumnos,
		PermisoRegistrarAsistencia,
		PermisoVerAsistencia,
		PermisoCrearEventos,
		PermisoVerEventos,
		PermisoCerrarEventos,
		PermisoVerAlertas,
		PermisoCerrarAlertas,
		PermisoVerMonitor,
		PermisoVerAuditoria,
		PermisoGestionarUsuarios,
		PermisoGestionarHorarios,
		PermisoImportarDatos,
		PermisoVerCasos,
		PermisoGestionarCasos,
		PermisoGestionarConceptos,
		PermisoGestionarAcciones,
		PermisoGestionarReglas,
		PermisoVerReportes,
		PermisoAdministrar,
	},
	models.RolProfesor: {
		PermisoVerCursos,
		PermisoVerAlumnos,
		PermisoRegistrarAsistencia,
		PermisoVerAsistencia,
		PermisoCrearEventos,
		PermisoVerEventos,
	},
	models.RolInspector: {
		PermisoVerCursos,
		PermisoVerAlumnos,
		PermisoVerAsistencia,
		PermisoVerEventos,
		PermisoCerrarEventos,
		PermisoVerAlertas,
		PermisoCerrarAlertas,
		PermisoVerMonitor,
		PermisoVerCasos,
	},
	models.RolAsistenteSocial: {
		PermisoVerCursos,
		PermisoVerAlumnos,
		PermisoVerAsistencia,
		PermisoVerEventos,
		PermisoVerCasos,
		PermisoGestionarCasos,
		PermisoVerReportes,
	},
	models.RolBackoffice: {
		PermisoVerCursos,
		PermisoVerAlumnos,
		PermisoVerAsistencia,
		PermisoVerEventos,
		PermisoVerCasos,
		PermisoGestionarConceptos,
		PermisoGestionarAcciones,
		PermisoGestionarReglas,
		PermisoVerAuditoria,
		PermisoGestionarUsuarios,
		PermisoGestionarHorarios,
		PermisoImportarDatos,
		PermisoVerMonitor,
		PermisoVerAlertas,
		PermisoCerrarAlertas,
		PermisoVerReportes,
	},
}

// TienePermiso verifica si un rol tiene un permiso especifico
func TienePermiso(rol string, permiso string) bool {
	permisos, ok := permisosPorRol[rol]
	if !ok {
		return false
	}

	for _, p := range permisos {
		if p == permiso {
			return true
		}
	}
	return false
}

// ObtenerPermisos retorna todos los permisos de un rol
func ObtenerPermisos(rol string) []string {
	permisos, ok := permisosPorRol[rol]
	if !ok {
		return []string{}
	}
	return permisos
}

// PermisosPorRol retorna una copia del mapa de permisos (para UI/admin)
func PermisosPorRol() map[string][]string {
	out := map[string][]string{}
	for k, v := range permisosPorRol {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

// TieneAlgunPermiso verifica si un rol tiene al menos uno de los permisos
func TieneAlgunPermiso(rol string, permisos ...string) bool {
	for _, permiso := range permisos {
		if TienePermiso(rol, permiso) {
			return true
		}
	}
	return false
}

// TieneTodosLosPermisos verifica si un rol tiene todos los permisos
func TieneTodosLosPermisos(rol string, permisos ...string) bool {
	for _, permiso := range permisos {
		if !TienePermiso(rol, permiso) {
			return false
		}
	}
	return true
}
