# Instrucciones para Levantar los Servicios

## Situación Actual

La API está lista, pero necesita una base de datos PostgreSQL corriendo.

## Opción 1: Usar Docker (Recomendado - Más Fácil)

Si tienes Docker instalado:

```bash
# Levantar PostgreSQL con Docker
docker-compose up -d# Verificar que está corriendo
docker ps# Ahora levantar el backend
cd source/api
./start.sh
```

## Opción 2: Instalar PostgreSQL Localmente

### En Ubuntu/Debian:

```bash
# Instalar PostgreSQL
sudo apt update
sudo apt install -y postgresql postgresql-contrib# Iniciar PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql# Crear base de datos
sudo -u postgres psql -c "CREATE DATABASE school_monitoring;"# (Opcional) Crear usuario específico
sudo -u postgres psql -c "CREATE USER tu_usuario WITH PASSWORD 'tu_password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE school_monitoring TO tu_usuario;"
```

### Configurar el archivo .env

Asegúrate de que tu archivo `source/api/.env` tenga:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=school_monitoring
DB_PORT=5432
PORT=8080
```

## Opción 3: Usar PostgreSQL Remoto

Si tienes acceso a un PostgreSQL remoto, solo configura el `.env` con:

```env
DB_HOST=tu_servidor_remoto
DB_USER=tu_usuario
DB_PASSWORD=tu_password
DB_NAME=school_monitoring
DB_PORT=5432
PORT=8080
```

## Levantar la API

Una vez que PostgreSQL esté disponible:

```bash
cd source/api
./start.sh
```

O manualmente:

```bash
cd source/api
go run ./cmd/server
```

La API estará disponible en: `http://localhost:8080`

## Verificar que Funciona

```bash
# Verificar el dashboard
curl http://localhost:8080/api/v1/dashboard# Debería retornar un JSON con información de salas y eventos
```

## Levantar Frontends (Opcional)

### Web (Profesores)
```bash
cd source/web
npm install  # Si npm está instalado
npm run dev  # Correrá en http://localhost:3000
```

### Backoffice (Inspectoría)
```bash
cd source/backoffice
npm install  # Si npm está instalado
npm run dev  # Correrá en http://localhost:3001
```

## Solución de Problemas

### Error: "connection refused"
- Verifica que PostgreSQL esté corriendo
- Verifica las credenciales en `.env`
- Verifica que el puerto 5432 esté abierto### Error: "database does not exist"
- Crea la base de datos: `CREATE DATABASE school_monitoring;`### Error: "npm not found"
- Instala Node.js y npm: `sudo apt install nodejs npm`
- O usa nvm para instalar Node.js
