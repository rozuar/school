# Sistema de Monitoreo Escolar (Tiempo Real)

Sistema centralizado para **monitor inspectoría** + **interfaz profesor** con semáforos, eventos, reglas, alertas y trazabilidad. Foco: **detectar, priorizar y asistir**.

## Estructura

```text
/school
  /source
    /api         Go (API + reglas + WS + persistencia)
    /web         Next.js (profesor)
    /backoffice  Next.js (inspectoría/admin)
    /android     App móvil (base)
```

## Quickstart (local)

### 1) DB

```bash
docker-compose up -d db
```

### 2) Backend

```bash
cd source/api
cp env.example .env
go run ./cmd/server
```

Backend: `http://localhost:8080` (WS: `ws://localhost:8080/ws`)

### 3) Frontend profesor

```bash
cd source/web
cp env.local.example .env.local
npm install
npm run dev -- --port 3000
```

Profesor: `http://localhost:3000`

### 4) Backoffice inspectoría/admin

```bash
cd source/backoffice
cp env.local.example .env.local
npm install
npm run dev -p 3001
```

Backoffice: `http://localhost:3001`

## Seed + demo users

```bash
curl -X POST http://localhost:8080/api/v1/seed
```

Usuarios demo:
- `profesor1@escuela.cl` / `profesor123`
- `inspector@escuela.cl` / `profesor123`
- `backoffice@escuela.cl` / `profesor123`

## Endpoints (principales)
- **Auth**: `POST /api/v1/auth/login`, `GET /api/v1/auth/me`, `GET /api/v1/auth/permisos`
- **Monitor**: `GET /api/v1/monitor/snapshot`
- **Alertas**: `GET /api/v1/alertas?estado=abierta`, `PUT /api/v1/alertas/{id}/cerrar`
- **Asistencia por bloque**: `POST /api/v1/asistencia/bloque`, `GET /api/v1/asistencia/horario/{id}/fecha/{fecha}`
- **Eventos**: `POST /api/v1/eventos`, `GET /api/v1/eventos/activos`
- **Trazabilidad**: `GET /api/v1/auditorias`, `GET /api/v1/acciones-ejecuciones`

## Comandos útiles (Makefile)

```bash
make install
make dev
make build
make test
```
- Integración con comunicación a apoderados
- Reportes automáticos para dirección
- Análisis de patrones y tendencias## ContribuciónEste es un sistema en desarrollo. Las mejoras y contribuciones son bienvenidas.## Licencia[Especificar licencia según corresponda]
