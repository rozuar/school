package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/school-monitoring/backend/internal/api/middleware"
	"github.com/school-monitoring/backend/internal/auth"
	"github.com/school-monitoring/backend/internal/models"
	"gorm.io/gorm"
)

// ImportHandler agrupa importaciones
type ImportHandler struct {
	db *gorm.DB
}

func NewImportHandler(db *gorm.DB) *ImportHandler {
	return &ImportHandler{db: db}
}

type ImportHorariosRequest struct {
	Formato string `json:"formato"` // "csv"
	CSV     string `json:"csv"`
	// Por si quieren defaults
	DefaultProfesorPassword string `json:"default_profesor_password,omitempty"`
}

type ImportHorariosResponse struct {
	RowsTotal      int `json:"rows_total"`
	RowsOK         int `json:"rows_ok"`
	RowsError      int `json:"rows_error"`
	HorariosCreados int `json:"horarios_creados"`
	HorariosActualizados int `json:"horarios_actualizados"`
	ProfesoresCreados int `json:"profesores_creados"`
	AsignaturasCreadas int `json:"asignaturas_creadas"`
	Errores        []string `json:"errores,omitempty"`
}

// ImportHorariosCSV importa horarios desde CSV pegado.
// Formato esperado (con o sin header):
// curso,dia_semana,bloque_numero,asignatura,profesor_email
func (h *ImportHandler) ImportHorariosCSV(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r)

	var req ImportHorariosRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.CSV) == "" {
		http.Error(w, `{"error":"csv es requerido"}`, http.StatusBadRequest)
		return
	}

	defaultPass := strings.TrimSpace(req.DefaultProfesorPassword)
	if defaultPass == "" {
		defaultPass = "profesor123"
	}

	reader := csv.NewReader(strings.NewReader(req.CSV))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, `{"error":"CSV invalido"}`, http.StatusBadRequest)
		return
	}

	resp := ImportHorariosResponse{
		RowsTotal: len(records),
	}

	// Cache de bloques por numero
	var bloques []models.BloqueHorario
	h.db.Order("numero").Find(&bloques)
	bloqueByNumero := map[string]models.BloqueHorario{}
	for _, b := range bloques {
		bloqueByNumero[strings.TrimSpace(string(rune('0'+b.Numero)))] = b
	}

	// Mejor: map por int como string
	bloqueByNum := map[int]models.BloqueHorario{}
	for _, b := range bloques {
		bloqueByNum[b.Numero] = b
	}

	// Ejecutar en transaccion
	err = h.db.Transaction(func(tx *gorm.DB) error {
		for idx, row := range records {
			// header detection
			if idx == 0 && len(row) >= 5 {
				if strings.Contains(strings.ToLower(row[0]), "curso") && strings.Contains(strings.ToLower(row[1]), "dia") {
					resp.RowsTotal--
					continue
				}
			}

			if len(row) < 5 {
				resp.RowsError++
				resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": columnas insuficientes (esperadas 5)")
				continue
			}

			cursoNombre := strings.TrimSpace(row[0])
			diaStr := strings.TrimSpace(row[1])
			bloqueStr := strings.TrimSpace(row[2])
			asigNombre := strings.TrimSpace(row[3])
			profEmail := strings.TrimSpace(strings.ToLower(row[4]))

			if cursoNombre == "" || diaStr == "" || bloqueStr == "" || asigNombre == "" || profEmail == "" {
				resp.RowsError++
				resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": campos vacios")
				continue
			}

			dia := atoi(diaStr)
			if dia < 1 || dia > 5 {
				resp.RowsError++
				resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": dia_semana invalido (1..5)")
				continue
			}
			bloqueNum := atoi(bloqueStr)
			bloque, ok := bloqueByNum[bloqueNum]
			if !ok {
				resp.RowsError++
				resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": bloque_numero no existe")
				continue
			}

			// Curso por nombre (case-insensitive)
			var curso models.Curso
			if err := tx.Where("LOWER(nombre) = LOWER(?)", cursoNombre).First(&curso).Error; err != nil {
				resp.RowsError++
				resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": curso no existe: "+cursoNombre)
				continue
			}

			// Asignatura (crear si no existe)
			var asig models.Asignatura
			if err := tx.Where("LOWER(nombre) = LOWER(?)", asigNombre).First(&asig).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					asig = models.Asignatura{Nombre: asigNombre}
					if err := tx.Create(&asig).Error; err != nil {
						resp.RowsError++
						resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": no se pudo crear asignatura")
						continue
					}
					resp.AsignaturasCreadas++
					_ = models.CrearAuditoria(tx, "asignaturas", asig.ID, models.AuditoriaInsert, nil, &asig, userIDPtr(claims))
				} else {
					resp.RowsError++
					resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": error buscando asignatura")
					continue
				}
			}

			// Profesor (crear si no existe)
			var prof models.Usuario
			if err := tx.Where("email = ?", profEmail).First(&prof).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					hash, _ := auth.HashPassword(defaultPass)
					prof = models.Usuario{
						Email:        profEmail,
						Nombre:       strings.Split(profEmail, "@")[0],
						Rol:          models.RolProfesor,
						PasswordHash: hash,
						Activo:       true,
					}
					if err := tx.Create(&prof).Error; err != nil {
						resp.RowsError++
						resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": no se pudo crear profesor")
						continue
					}
					resp.ProfesoresCreados++
					_ = models.CrearAuditoria(tx, "usuarios", prof.ID, models.AuditoriaInsert, nil, &prof, userIDPtr(claims))
				} else {
					resp.RowsError++
					resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": error buscando profesor")
					continue
				}
			}

			// Upsert Horario por (curso, dia, bloque)
			var existing models.Horario
			err := tx.Where("curso_id = ? AND dia_semana = ? AND bloque_id = ?", curso.ID, dia, bloque.ID).
				First(&existing).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				resp.RowsError++
				resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": error buscando horario")
				continue
			}

			if err == gorm.ErrRecordNotFound {
				hh := models.Horario{
					CursoID:      curso.ID,
					AsignaturaID: asig.ID,
					ProfesorID:   prof.ID,
					BloqueID:     bloque.ID,
					DiaSemana:    dia,
				}
				if err := tx.Create(&hh).Error; err != nil {
					resp.RowsError++
					resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": error creando horario")
					continue
				}
				resp.HorariosCreados++
				_ = models.CrearAuditoria(tx, "horarios", hh.ID, models.AuditoriaInsert, nil, &hh, userIDPtr(claims))
				resp.RowsOK++
				continue
			}

			before := existing
			existing.AsignaturaID = asig.ID
			existing.ProfesorID = prof.ID
			if err := tx.Save(&existing).Error; err != nil {
				resp.RowsError++
				resp.Errores = append(resp.Errores, "fila "+itoa(idx+1)+": error actualizando horario")
				continue
			}
			resp.HorariosActualizados++
			_ = models.CrearAuditoria(tx, "horarios", existing.ID, models.AuditoriaUpdate, &before, &existing, userIDPtr(claims))
			resp.RowsOK++
		}
		return nil
	})
	if err != nil {
		http.Error(w, `{"error":"Error importando horarios"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func userIDPtr(claims *auth.Claims) *uuid.UUID {
	if claims == nil {
		return nil
	}
	return &claims.UserID
}

// mini helpers (evitar fmt/strconv para mantener compacto)
func atoi(s string) int {
	n := 0
	sign := 1
	ss := strings.TrimSpace(s)
	if strings.HasPrefix(ss, "-") {
		sign = -1
		ss = strings.TrimPrefix(ss, "-")
	}
	for _, ch := range ss {
		if ch < '0' || ch > '9' {
			break
		}
		n = n*10 + int(ch-'0')
	}
	return n * sign
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	var b [32]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + (n % 10))
		n /= 10
	}
	return sign + string(b[i:])
}



