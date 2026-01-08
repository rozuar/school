# Backend - Sistema de Monitoreo de Salas## Requisitos Previos- Go 1.21 o superior
- PostgreSQL 12 o superior## Configuración1. **Crear archivo `.env`** con la configuración de la base de datos:```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=school_monitoring
DB_PORT=5432
PORT=8080
```2. **Configurar PostgreSQL**:```bash
# Opción 1: Usar el script de configuración
./setup-db.sh# Opción 2: Configuración manual
sudo systemctl start postgresql
sudo -u postgres psql -c "CREATE DATABASE school_monitoring;"
```3. **Instalar dependencias**:```bash
go mod download
go mod tidy
```## Ejecución### Opción 1: Usar el script de inicio
```bash
./start.sh
```### Opción 2: Ejecutar directamente
```bash
go run main.go
# o
go build -o main .
./main
```## VerificaciónEl servidor estará disponible en `http://localhost:8080`Puedes verificar que está funcionando:
```bash
curl http://localhost:8080/api/v1/dashboard
```## Endpoints Disponibles- `GET /api/v1/salas` - Obtener todas las salas
- `GET /api/v1/salas/{id}` - Obtener una sala específica
- `GET /api/v1/dashboard` - Dashboard de monitoreo
- `GET /api/v1/eventos` - Eventos activos
- `POST /api/v1/profesores/{id}/asistencia` - Registrar asistencia
- `POST /api/v1/eventos` - Crear evento
- `PUT /api/v1/eventos/{id}/cerrar` - Cerrar evento
- `WS /ws` - WebSocket para actualizaciones en tiempo real## Solución de Problemas### Error: "connection refused"
- Verifica que PostgreSQL esté corriendo: `sudo systemctl status postgresql`
- Inicia PostgreSQL: `sudo systemctl start postgresql`
- Verifica las credenciales en el archivo `.env`### Error: "database does not exist"
- Crea la base de datos: `sudo -u postgres psql -c "CREATE DATABASE school_monitoring;"`
- O ejecuta el script: `./setup-db.sh`### Error: "permission denied"
- Verifica los permisos del usuario de PostgreSQL
- Asegúrate de que el usuario tenga acceso a la base de datos
