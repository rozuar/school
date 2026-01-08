# Plan de Desarrollo - Plataforma de Gestion Escolar

## Fase 1: Fundamentos (MVP Demo)

### 1.1 Setup del Proyecto
- [ ] Inicializar monorepo con estructura de carpetas
- [ ] Configurar backend Go con Gin/Fiber
- [ ] Configurar frontend Next.js (profesor/admin)
- [ ] Configurar backoffice Next.js
- [ ] Setup PostgreSQL en Railway
- [ ] Configurar variables de entorno y secretos

### 1.2 Modelo de Datos Base
- [ ] Tabla `establecimientos`
- [ ] Tabla `usuarios` (con roles: profesor, admin, asistente_social, backoffice)
- [ ] Tabla `cursos` (1ro basico a 4to medio)
- [ ] Tabla `alumnos`
- [ ] Tabla `asignaturas`
- [ ] Tabla `profesores_asignaturas` (relacion profesor-asignatura-curso)
- [ ] Tabla `bloques_horarios`
- [ ] Tabla `horarios` (asignatura-curso-bloque-dia)

### 1.3 Autenticacion y Autorizacion
- [ ] Implementar JWT auth en backend
- [ ] Middleware de autenticacion
- [ ] Sistema RBAC (roles y permisos)
- [ ] Endpoints: login, logout, refresh token
- [ ] Proteccion de rutas en frontend

### 1.4 Gestion de Asistencia
- [ ] Tabla `asistencias` (por bloque, no por dia)
- [ ] Tabla `estados_alumno` (presente, ausente, bano, enfermeria, sos)
- [ ] Endpoint: obtener lista de alumnos por curso
- [ ] Endpoint: registrar asistencia de bloque
- [ ] Endpoint: cambiar estado temporal (bano/enfermeria)
- [ ] UI profesor: lista de alumnos con iconos expandibles
- [ ] UI profesor: semaforo visual (negro=ausente, amarillo=bano, rojo=sos)

### 1.5 Sistema de Conceptos, Acciones y Reglas
- [ ] Tabla `conceptos` (eventos estandarizados)
- [ ] Tabla `acciones` (respuestas configurables)
- [ ] Tabla `reglas` (condiciones temporales)
- [ ] Tabla `eventos` (ocurrencias concretas con timestamp)
- [ ] Tabla `auditoria` (trazabilidad completa)
- [ ] Motor de reglas basico (evaluar condiciones)
- [ ] Backoffice: CRUD conceptos
- [ ] Backoffice: CRUD acciones
- [ ] Backoffice: CRUD reglas

### 1.6 Notificaciones Basicas
- [ ] Servicio de notificaciones
- [ ] Integracion email (SendGrid/SMTP)
- [ ] Notificacion automatica por inasistencia
- [ ] Cola de notificaciones pendientes

### 1.7 Dashboard Administrativo
- [ ] Vista general de asistencia del dia
- [ ] Listado de alertas activas
- [ ] Resumen por curso
- [ ] Filtros por fecha y curso

### 1.8 Deploy Demo
- [ ] Configurar Railway para backend
- [ ] Configurar Railway para frontend
- [ ] Configurar Railway para backoffice
- [ ] CI/CD basico (GitHub Actions)
- [ ] Documentacion de deploy

---

## Fase 2: Funcionalidades Completas (Piloto)

### 2.1 Motor de Reglas Completo
- [ ] Evaluacion de reglas temporales complejas
- [ ] Regla: 2 inasistencias sin justificativo en 7 dias
- [ ] Regla: 3 inasistencias no consecutivas en 30 dias
- [ ] Regla: casos especiales inhiben alertas
- [ ] Procesamiento asincrono de reglas (Pub/Sub)

### 2.2 Perfil Asistente Social
- [ ] Dashboard de casos asignados
- [ ] Alertas de seguimiento
- [ ] Registro de intervenciones
- [ ] Historial del alumno
- [ ] Marcado de casos especiales

### 2.3 Gestion de Casos Especiales
- [ ] Tabla `casos_especiales`
- [ ] Flujo de marcado/desmarcado
- [ ] Inhibicion de alertas automaticas
- [ ] Monitoreo continuo sin notificaciones

### 2.4 Justificacion de Inasistencias
- [ ] Endpoint para subir justificativos
- [ ] Tipos: certificado medico, justificacion libre
- [ ] Cambio de estado de inasistencia
- [ ] Notificacion de justificacion recibida

### 2.5 Mobile React Native
- [ ] Setup proyecto React Native
- [ ] Autenticacion mobile
- [ ] Vista de casos asignados (inspector/asistente social)
- [ ] Recepcion de alertas push
- [ ] Confirmacion de acciones en terreno
- [ ] Integracion Firebase Cloud Messaging

### 2.6 Reportes Basicos
- [ ] Reporte de asistencia por curso
- [ ] Reporte de asistencia por alumno
- [ ] Reporte de inasistencias justificadas/no justificadas
- [ ] Exportacion a PDF/Excel

---

## Fase 3: Produccion GCP

### 3.1 Migracion de Infraestructura
- [ ] Cloud Run para backend Go
- [ ] Cloud Run para frontends Next.js
- [ ] Cloud SQL PostgreSQL
- [ ] Configurar Pub/Sub para eventos
- [ ] Cloud Tasks para acciones diferidas
- [ ] Secret Manager para credenciales

### 3.2 Escalabilidad
- [ ] Auto-scaling en Cloud Run
- [ ] Connection pooling para DB
- [ ] Cache con Redis/Memorystore
- [ ] CDN para assets estaticos

### 3.3 Observabilidad
- [ ] Cloud Logging estructurado
- [ ] Cloud Monitoring dashboards
- [ ] Alertas de errores y latencia
- [ ] Tracing distribuido

### 3.4 Seguridad Avanzada
- [ ] Cloud Load Balancer con SSL
- [ ] IAM roles granulares
- [ ] Identity Platform para auth
- [ ] Auditoria de accesos
- [ ] Backups automaticos de DB

### 3.5 Alta Disponibilidad
- [ ] Multi-zona para Cloud SQL
- [ ] Health checks en Cloud Run
- [ ] Failover automatico
- [ ] Plan de disaster recovery

---

## Fase 4: Funcionalidades Extendidas (Futuro)

### 4.1 Portal Apoderados
- [ ] Registro de apoderados
- [ ] Vista de asistencia del alumno
- [ ] Subida de justificativos
- [ ] Comunicacion con profesores
- [ ] Notificaciones de eventos

### 4.2 Gestion Academica
- [ ] Registro de calificaciones
- [ ] Boletines de notas
- [ ] Informes academicos
- [ ] Planificacion de clases

### 4.3 Multi-establecimiento
- [ ] Soporte para multiples colegios
- [ ] Configuracion por establecimiento
- [ ] Reportes consolidados
- [ ] Admin de plataforma

---

## Estructura de Repositorio

```
school/
├── context/
│   ├── raw/           # Documentos originales
│   └── refined/       # Documentacion procesada
└── source/
    ├── backend/       # Go API
    │   ├── cmd/
    │   ├── internal/
    │   │   ├── api/
    │   │   ├── auth/
    │   │   ├── models/
    │   │   ├── services/
    │   │   └── rules/
    │   └── migrations/
    ├── frontend/      # Next.js Profesor/Admin
    │   └── src/
    │       ├── app/
    │       ├── components/
    │       └── services/
    ├── backoffice/    # Next.js Backoffice
    │   └── src/
    │       ├── app/
    │       ├── components/
    │       └── services/
    └── mobile/        # React Native
        └── src/
```

---

## Endpoints API (MVP)

### Auth
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `POST /api/auth/refresh`

### Cursos
- `GET /api/cursos`
- `GET /api/cursos/:id/alumnos`

### Asistencia
- `POST /api/asistencia/bloque`
- `GET /api/asistencia/curso/:id/fecha/:fecha`
- `PUT /api/alumnos/:id/estado`

### Conceptos
- `GET /api/conceptos`
- `POST /api/conceptos`
- `PUT /api/conceptos/:id`
- `DELETE /api/conceptos/:id`

### Acciones
- `GET /api/acciones`
- `POST /api/acciones`
- `PUT /api/acciones/:id`

### Reglas
- `GET /api/reglas`
- `POST /api/reglas`
- `PUT /api/reglas/:id`

### Eventos
- `POST /api/eventos`
- `GET /api/eventos/alumno/:id`

### Notificaciones
- `GET /api/notificaciones`
- `PUT /api/notificaciones/:id/leida`

---

## Modelo de Datos Inicial

```sql
-- Establecimientos
CREATE TABLE establecimientos (
    id UUID PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Usuarios
CREATE TABLE usuarios (
    id UUID PRIMARY KEY,
    establecimiento_id UUID REFERENCES establecimientos(id),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    rol VARCHAR(50) NOT NULL, -- profesor, admin, asistente_social, backoffice
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Cursos
CREATE TABLE cursos (
    id UUID PRIMARY KEY,
    establecimiento_id UUID REFERENCES establecimientos(id),
    nombre VARCHAR(50) NOT NULL, -- "1 Basico", "4 Medio"
    nivel VARCHAR(20) NOT NULL, -- basica, media
    created_at TIMESTAMP DEFAULT NOW()
);

-- Alumnos
CREATE TABLE alumnos (
    id UUID PRIMARY KEY,
    curso_id UUID REFERENCES cursos(id),
    rut VARCHAR(12) UNIQUE,
    nombre VARCHAR(255) NOT NULL,
    apellido VARCHAR(255) NOT NULL,
    caso_especial BOOLEAN DEFAULT FALSE,
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Asignaturas
CREATE TABLE asignaturas (
    id UUID PRIMARY KEY,
    establecimiento_id UUID REFERENCES establecimientos(id),
    nombre VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Bloques Horarios
CREATE TABLE bloques_horarios (
    id UUID PRIMARY KEY,
    establecimiento_id UUID REFERENCES establecimientos(id),
    numero INT NOT NULL,
    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL
);

-- Horarios
CREATE TABLE horarios (
    id UUID PRIMARY KEY,
    curso_id UUID REFERENCES cursos(id),
    asignatura_id UUID REFERENCES asignaturas(id),
    profesor_id UUID REFERENCES usuarios(id),
    bloque_id UUID REFERENCES bloques_horarios(id),
    dia_semana INT NOT NULL -- 1=lunes, 5=viernes
);

-- Asistencias
CREATE TABLE asistencias (
    id UUID PRIMARY KEY,
    alumno_id UUID REFERENCES alumnos(id),
    horario_id UUID REFERENCES horarios(id),
    fecha DATE NOT NULL,
    estado VARCHAR(20) NOT NULL, -- presente, ausente, justificado
    registrado_por UUID REFERENCES usuarios(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Estados Temporales
CREATE TABLE estados_temporales (
    id UUID PRIMARY KEY,
    alumno_id UUID REFERENCES alumnos(id),
    tipo VARCHAR(20) NOT NULL, -- bano, enfermeria, sos
    inicio TIMESTAMP NOT NULL,
    fin TIMESTAMP,
    registrado_por UUID REFERENCES usuarios(id)
);

-- Conceptos
CREATE TABLE conceptos (
    id UUID PRIMARY KEY,
    establecimiento_id UUID REFERENCES establecimientos(id),
    codigo VARCHAR(50) UNIQUE NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    activo BOOLEAN DEFAULT TRUE
);

-- Acciones
CREATE TABLE acciones (
    id UUID PRIMARY KEY,
    establecimiento_id UUID REFERENCES establecimientos(id),
    codigo VARCHAR(50) UNIQUE NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    tipo VARCHAR(50) NOT NULL, -- notificacion, alerta, cambio_estado
    parametros JSONB,
    activo BOOLEAN DEFAULT TRUE
);

-- Reglas
CREATE TABLE reglas (
    id UUID PRIMARY KEY,
    establecimiento_id UUID REFERENCES establecimientos(id),
    nombre VARCHAR(100) NOT NULL,
    condicion JSONB NOT NULL,
    accion_id UUID REFERENCES acciones(id),
    activo BOOLEAN DEFAULT TRUE
);

-- Eventos
CREATE TABLE eventos (
    id UUID PRIMARY KEY,
    concepto_id UUID REFERENCES conceptos(id),
    alumno_id UUID REFERENCES alumnos(id),
    curso_id UUID REFERENCES cursos(id),
    origen VARCHAR(50) NOT NULL, -- profesor, sistema
    origen_usuario_id UUID REFERENCES usuarios(id),
    datos JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Auditoria
CREATE TABLE auditoria (
    id UUID PRIMARY KEY,
    tabla VARCHAR(50) NOT NULL,
    registro_id UUID NOT NULL,
    accion VARCHAR(20) NOT NULL, -- INSERT, UPDATE, DELETE
    datos_anteriores JSONB,
    datos_nuevos JSONB,
    usuario_id UUID REFERENCES usuarios(id),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Criterios de Exito

- [ ] Sistema funciona sin hardcodear reglas
- [ ] Conceptos modificables desde backoffice
- [ ] Acciones se ejecutan correctamente
- [ ] Alertas se disparan solo cuando corresponde
- [ ] Demo presentable a colegio o sostenedor
- [ ] Asistencia por bloque funcionando
- [ ] Semaforo visual implementado
- [ ] Notificaciones llegando a apoderados
