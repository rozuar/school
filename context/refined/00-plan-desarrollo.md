# Plan de Desarrollo - Plataforma de Gestion Escolar

## Fase 1: Fundamentos (MVP Demo)

### 1.1 Setup del Proyecto
- [ ] Inicializar monorepo con estructura de carpetas
- [ ] Configurar API Go (carpeta `source/api`)
- [ ] Configurar Web Next.js (carpeta `source/web`)
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
- [ ] Implementar JWT auth en la API
- [ ] Middleware de autenticacion
- [ ] Sistema RBAC (roles y permisos)
- [ ] Endpoints: login, logout, refresh token
- [ ] Proteccion de rutas en web/backoffice

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
- [ ] Configurar Railway para API
- [ ] Configurar Railway para web
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
- [ ] Autenticacion en mobile
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
- [ ] Cloud Run para API Go
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

## Criterios de Exito

- [ ] Sistema funciona sin hardcodear reglas
- [ ] Conceptos modificables desde backoffice
- [ ] Acciones se ejecutan correctamente
- [ ] Alertas se disparan solo cuando corresponde
- [ ] Demo presentable a colegio o sostenedor
- [ ] Asistencia por bloque funcionando
- [ ] Semaforo visual implementado
- [ ] Notificaciones llegando a apoderados

---

## Documentos de Referencia (MVP)

- **Estructura repo**: ver `09-estructura-repositorio.md`
- **API (MVP)**: ver `10-api-mvp.md`
- **Modelo de datos (borrador SQL)**: ver `11-modelo-datos.md`
