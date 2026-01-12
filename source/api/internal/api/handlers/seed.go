package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// SeedHandler maneja el seed de datos demo
type SeedHandler struct {
	db *gorm.DB
}

// NewSeedHandler crea un nuevo handler de seed
func NewSeedHandler(db *gorm.DB) *SeedHandler {
	return &SeedHandler{db: db}
}

// Seed crea datos demo
func (h *SeedHandler) Seed(w http.ResponseWriter, r *http.Request) {
	// Crear usuarios
	adminPassword, _ := auth.HashPassword("admin123")
	profPassword, _ := auth.HashPassword("profesor123")

	usuarios := []models.Usuario{
		{Email: "admin@escuela.cl", PasswordHash: adminPassword, Nombre: "Administrador", Rol: models.RolAdmin},
		{Email: "profesor1@escuela.cl", PasswordHash: profPassword, Nombre: "Juan Perez", Rol: models.RolProfesor},
		{Email: "profesor2@escuela.cl", PasswordHash: profPassword, Nombre: "Maria Garcia", Rol: models.RolProfesor},
		{Email: "profesor3@escuela.cl", PasswordHash: profPassword, Nombre: "Carlos Lopez", Rol: models.RolProfesor},
		{Email: "inspector@escuela.cl", PasswordHash: profPassword, Nombre: "Pedro Martinez", Rol: models.RolInspector},
		{Email: "asistente@escuela.cl", PasswordHash: profPassword, Nombre: "Ana Rodriguez", Rol: models.RolAsistenteSocial},
		{Email: "backoffice@escuela.cl", PasswordHash: profPassword, Nombre: "Luis Sanchez", Rol: models.RolBackoffice},
	}

	for i := range usuarios {
		h.db.FirstOrCreate(&usuarios[i], models.Usuario{Email: usuarios[i].Email})
	}

	// Crear cursos
	cursos := []models.Curso{
		{Nombre: "1 Basico", Nivel: models.NivelBasica},
		{Nombre: "2 Basico", Nivel: models.NivelBasica},
		{Nombre: "3 Basico", Nivel: models.NivelBasica},
		{Nombre: "4 Basico", Nivel: models.NivelBasica},
		{Nombre: "5 Basico", Nivel: models.NivelBasica},
		{Nombre: "6 Basico", Nivel: models.NivelBasica},
		{Nombre: "7 Basico", Nivel: models.NivelBasica},
		{Nombre: "8 Basico", Nivel: models.NivelBasica},
		{Nombre: "1 Medio", Nivel: models.NivelMedia},
		{Nombre: "2 Medio", Nivel: models.NivelMedia},
		{Nombre: "3 Medio", Nivel: models.NivelMedia},
		{Nombre: "4 Medio", Nivel: models.NivelMedia},
	}

	for i := range cursos {
		h.db.FirstOrCreate(&cursos[i], models.Curso{Nombre: cursos[i].Nombre})
	}

	// Recargar cursos con IDs
	h.db.Find(&cursos)

	// Crear alumnos para cada curso
	nombres := []string{"Santiago", "Martina", "Mateo", "Sofia", "Benjamin", "Valentina", "Lucas", "Isabella", "Agustin", "Emma"}
	apellidos := []string{"Gonzalez", "Rodriguez", "Martinez", "Lopez", "Garcia", "Hernandez", "Perez", "Sanchez", "Ramirez", "Torres"}

	for _, curso := range cursos {
		for i := 0; i < 10; i++ {
			alumno := models.Alumno{
				CursoID:  curso.ID,
				Nombre:   nombres[i],
				Apellido: apellidos[i],
				Rut:      generarRut(curso.Nombre, i),
				Activo:   true,
			}
			// Upsert por RUT para asegurar que el alumno quede asociado al curso correcto
			h.db.Where("rut = ?", alumno.Rut).Assign(alumno).FirstOrCreate(&alumno)
		}
	}

	// Crear asignaturas
	asignaturas := []models.Asignatura{
		{Nombre: "Matematica"},
		{Nombre: "Lenguaje"},
		{Nombre: "Ciencias Naturales"},
		{Nombre: "Historia"},
		{Nombre: "Ingles"},
		{Nombre: "Educacion Fisica"},
		{Nombre: "Arte"},
		{Nombre: "Musica"},
		{Nombre: "Tecnologia"},
	}

	for i := range asignaturas {
		h.db.FirstOrCreate(&asignaturas[i], models.Asignatura{Nombre: asignaturas[i].Nombre})
	}

	// Crear bloques horarios
	bloques := []models.BloqueHorario{
		{Numero: 1, HoraInicio: "09:00", HoraFin: "10:00"},
		{Numero: 2, HoraInicio: "10:30", HoraFin: "11:30"},
		{Numero: 3, HoraInicio: "12:00", HoraFin: "13:00"},
		{Numero: 4, HoraInicio: "14:00", HoraFin: "15:00"},
		{Numero: 5, HoraInicio: "15:30", HoraFin: "16:30"},
	}

	for i := range bloques {
		h.db.FirstOrCreate(&bloques[i], models.BloqueHorario{Numero: bloques[i].Numero})
	}

	// Recargar asignaturas y bloques con IDs (por si ya existian)
	h.db.Find(&asignaturas)
	h.db.Find(&bloques)

	// Obtener profesores para asignar horarios
	var profesores []models.Usuario
	h.db.Where("rol = ? AND activo = ?", models.RolProfesor, true).Order("email").Find(&profesores)

	// Crear horarios demo (5 dias x 5 bloques) para cada curso
	// Si ya existe un horario para curso+dia+bloque, se mantiene estable (o se actualiza con asignatura/profesor del seed)
	if len(profesores) > 0 && len(asignaturas) > 0 && len(bloques) > 0 {
		for ci, curso := range cursos {
			for dia := 1; dia <= 5; dia++ {
				for bi, bloque := range bloques {
					asig := asignaturas[(ci+dia+bi)%len(asignaturas)]
					prof := profesores[(ci+dia+bi)%len(profesores)]

					hh := models.Horario{
						CursoID:      curso.ID,
						AsignaturaID: asig.ID,
						ProfesorID:   prof.ID,
						BloqueID:     bloque.ID,
						DiaSemana:    dia,
					}

					// Unico por (curso, dia, bloque)
					h.db.Where("curso_id = ? AND dia_semana = ? AND bloque_id = ?", curso.ID, dia, bloque.ID).
						Assign(hh).
						FirstOrCreate(&hh)
				}
			}
		}
	}

	// Crear conceptos
	conceptos := []models.Concepto{
		{Codigo: models.ConceptoInasistencia, Nombre: "Inasistencia", Descripcion: "Alumno ausente a clases"},
		{Codigo: models.ConceptoBano, Nombre: "Bano", Descripcion: "Alumno fue al bano"},
		{Codigo: models.ConceptoEnfermeria, Nombre: "Enfermeria", Descripcion: "Alumno en enfermeria"},
		{Codigo: models.ConceptoSOS, Nombre: "SOS", Descripcion: "Situacion de emergencia"},
		{Codigo: models.ConceptoComportamiento, Nombre: "Comportamiento", Descripcion: "Incidente de comportamiento"},
		{Codigo: models.ConceptoDisciplinario, Nombre: "Disciplinario", Descripcion: "Problema disciplinario"},
	}

	for i := range conceptos {
		h.db.FirstOrCreate(&conceptos[i], models.Concepto{Codigo: conceptos[i].Codigo})
	}

	// Recargar conceptos con IDs
	h.db.Find(&conceptos)

	// Crear acciones
	acciones := []models.Accion{
		{Codigo: "NOTIFICAR_APODERADO", Nombre: "Notificar Apoderado", Tipo: models.TipoAccionNotificacion, Parametros: []byte(`{"destinatario": "apoderado"}`)},
		{Codigo: "ALERTA_INSPECTOR", Nombre: "Alerta a Inspector", Tipo: models.TipoAccionAlerta, Parametros: []byte(`{"destinatario": "inspector", "prioridad": "media"}`)},
		{Codigo: "ALERTA_ASISTENTE", Nombre: "Alerta a Asistente Social", Tipo: models.TipoAccionAlerta, Parametros: []byte(`{"destinatario": "asistente_social", "prioridad": "alta"}`)},
	}

	for i := range acciones {
		h.db.FirstOrCreate(&acciones[i], models.Accion{Codigo: acciones[i].Codigo})
	}

	// Recargar acciones con IDs
	h.db.Find(&acciones)

	// Crear reglas
	var conceptoInasistencia models.Concepto
	h.db.First(&conceptoInasistencia, "codigo = ?", models.ConceptoInasistencia)

	var accionNotificar models.Accion
	h.db.First(&accionNotificar, "codigo = ?", "NOTIFICAR_APODERADO")

	var accionAlertaAsistente models.Accion
	h.db.First(&accionAlertaAsistente, "codigo = ?", "ALERTA_ASISTENTE")

	reglas := []models.Regla{
		{
			Nombre:     "Notificar por inasistencia",
			ConceptoID: conceptoInasistencia.ID,
			Condicion:  []byte(`{"tipo": "cantidad", "campo": "inasistencias", "operador": ">=", "valor": 1, "dias": 1}`),
			AccionID:   accionNotificar.ID,
		},
		{
			Nombre:     "Alerta por 2 inasistencias en 7 dias",
			ConceptoID: conceptoInasistencia.ID,
			Condicion:  []byte(`{"tipo": "cantidad", "campo": "inasistencias", "operador": ">=", "valor": 2, "dias": 7}`),
			AccionID:   accionAlertaAsistente.ID,
		},
	}

	for i := range reglas {
		h.db.FirstOrCreate(&reglas[i], models.Regla{Nombre: reglas[i].Nombre})
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Seed completed successfully",
		"usuarios":    "7 usuarios creados",
		"cursos":      "12 cursos creados",
		"alumnos":     "120 alumnos creados",
		"asignaturas": "9 asignaturas creadas",
		"bloques":     "5 bloques horarios creados",
		"conceptos":   "6 conceptos creados",
		"acciones":    "3 acciones creadas",
		"reglas":      "2 reglas creadas",
	})
}

func generarRut(curso string, indice int) string {
	// RUT demo: debe ser estable y único (hay uniqueIndex en DB).
	// Usamos un identificador determinístico por curso + índice (no es un RUT real).
	key := strings.ToUpper(strings.ReplaceAll(curso, " ", ""))
	return fmt.Sprintf("DEMO-%s-%02d", key, indice+1)
}
