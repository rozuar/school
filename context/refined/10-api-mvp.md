# API (MVP)

Listado mínimo de endpoints para operar el MVP.

## Base

- Base REST: `/api/v1`
- WebSocket: `/ws`

## Auth

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/refresh`
- `GET /api/v1/auth/me`

## Cursos y alumnos

- `GET /api/v1/cursos`
- `GET /api/v1/cursos/{id}`
- `GET /api/v1/cursos/{id}/alumnos`
- `GET /api/v1/cursos/{id}/horario`

## Horarios (profesor)

- `GET /api/v1/horarios/mis`

## Asistencia

- `POST /api/v1/asistencia/bloque`
- `GET /api/v1/asistencia/curso/{id}/fecha/{fecha}`
- `GET /api/v1/asistencia/horario/{id}/fecha/{fecha}`

## Estados temporales

- `PUT /api/v1/alumnos/{id}/estado-temporal`
- `DELETE /api/v1/alumnos/{id}/estado-temporal`
- `GET /api/v1/estados-temporales`

## Configuración (backoffice)

- `GET /api/v1/conceptos`
- `POST /api/v1/conceptos`
- `PUT /api/v1/conceptos/{id}`
- `DELETE /api/v1/conceptos/{id}`

- `GET /api/v1/acciones`
- `POST /api/v1/acciones`
- `PUT /api/v1/acciones/{id}`
- `DELETE /api/v1/acciones/{id}`

- `GET /api/v1/reglas`
- `POST /api/v1/reglas`
- `PUT /api/v1/reglas/{id}`
- `DELETE /api/v1/reglas/{id}`

## Eventos

- `GET /api/v1/eventos`
- `GET /api/v1/eventos/activos`
- `GET /api/v1/eventos/alumno/{id}`
- `POST /api/v1/eventos`
- `PUT /api/v1/eventos/{id}/cerrar`

## Monitor / Operación

- `GET /api/v1/monitor/snapshot`
- `GET /api/v1/alertas`
