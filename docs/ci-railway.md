# CI/CD hacia Railway (arquitectura)

## Objetivo

Tener un **pipeline de CI** que valide el monorepo (backend + frontend + backoffice) en cada PR/push, y un **CD hacia Railway** que despliegue automáticamente desde `main`.

## Componentes

- **GitHub**: repositorio + PRs + Actions.
- **GitHub Actions (CI)**: corre tests/build (no despliega por defecto).
- **Railway (CD)**: compila y despliega cada servicio desde GitHub o gatillado por CLI.

## Servicios en Railway (monorepo)

Crear **3 servicios** apuntando al mismo repo, cada uno con su `Root Directory`:

- **api**: `source/api`
- **web**: `source/web`
- **backoffice**: `source/backoffice`

## Variables críticas (Front/Backoffice -> Backend)

Este repo asume que **frontend/backoffice NO hablan a `localhost` en producción**.

- Configura en Railway (en `source/web` y `source/backoffice`):
  - `BACKEND_URL=https://<tu-backend>.up.railway.app` (o URL interna)
  - `NEXT_PUBLIC_BACKEND_URL` (solo si algo necesita usarla en runtime del navegador)

En producción (Railway) se recomienda usar **URL interna**:

- `BACKEND_URL=http://<backend-service>.railway.internal`

Así el proxy de `/api/*` y `/ws` funciona dentro del VPC de Railway y reduces latencia.

## Build/Start (recomendado)

### Backend (Go)

- **Build**: `go build -o server ./cmd/server`
- **Start**: `./server`
- **Env**:
  - `DATABASE_URL` (Railway Postgres)
  - `JWT_SECRET`
  - `PORT` (Railway lo inyecta)
  - Retención (opcionales): `AUDIT_RETENTION_DAYS`, `EXEC_RETENTION_DAYS`, `OUTBOX_RETENTION_DAYS`, `ALERT_RETENTION_DAYS`

### Frontend (Next.js - profesor)

- **Build**: `npm ci && npm run build`
- **Start**: `npm run start -- -p $PORT`
- **Env**:
  - `BACKEND_URL` y/o `NEXT_PUBLIC_BACKEND_URL`

### Backoffice (Next.js - inspectoría/admin)

- **Build**: `npm ci && npm run build`
- **Start**: `npm run start -- -p $PORT`
- **Env**:
  - `BACKEND_URL` y/o `NEXT_PUBLIC_BACKEND_URL`

## Proxy de API y WebSocket (Next.js)

Para evitar CORS y para que WebSocket funcione desde el dominio del front/backoffice, ambos apps usan `rewrites()`:

- `/api/:path*` -> `${BACKEND_URL}/api/:path*`
- `/ws` -> `${BACKEND_URL}/ws`

Y el backoffice abre WS contra `wss://<tu-app>/ws` (no `:8080`), que luego es proxied al backend.

## Conectividad entre servicios (Railway)

Opción A (simple): usar **URL pública** del backend en `BACKEND_URL`.

Opción B (mejor): usar **dominio interno** (private networking) del backend, por ejemplo:

- `BACKEND_URL=http://<backend-service>.railway.internal`

Así el front/backoffice consumen API/WS internamente, y solo expones HTTP público donde lo necesitas.

## Flujo CI (GitHub Actions)

En cada PR/push:

- **backend**: `go test ./...`
- **frontend/backoffice**: `npm ci` + `npm run build`

Esto asegura que Railway reciba commits “verdes”.

## Flujo CD (Railway)

### Opción 1 (recomendada): Railway GitHub Deploys

- Railway escucha pushes a `main` y despliega automáticamente por servicio.
- GitHub Actions queda como **gating** (checks requeridos en PR).

### Opción 2: GitHub Actions despliega con Railway CLI

Requiere secrets en GitHub:

- `RAILWAY_TOKEN`
- `RAILWAY_PROJECT_ID`
- `RAILWAY_ENVIRONMENT_ID` (opcional)
- IDs por servicio (si se separa el deploy por servicio)

El job de deploy corre solo en `push` a `main` (o manual) y ejecuta `railway` para disparar despliegues.

## Diagrama (alto nivel)

```text
PR -> GitHub Actions (CI: tests/build) -> merge main
                                   |
                                   v
                            Railway CD (deploy)
                 +----------------+----------------+
                 |                |                |
            backend service   frontend service  backoffice service
```


