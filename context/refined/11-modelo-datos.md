# Modelo de Datos (borrador SQL)

Este es un **borrador de referencia** (basado en el proyecto descrito en `raw/idea.txt`).

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
    rut VARCHAR(64) UNIQUE,
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
