#!/bin/bash# Script para configurar la base de datos PostgreSQLecho "=== Configuración de Base de Datos ==="# Verificar si PostgreSQL está instalado
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL no está instalado."
    echo "Instalando PostgreSQL..."
    sudo apt update
    sudo apt install -y postgresql postgresql-contrib
fi# Verificar si PostgreSQL está corriendo
if ! sudo systemctl is-active --quiet postgresql; then
    echo "Iniciando PostgreSQL..."
    sudo systemctl start postgresql
    sudo systemctl enable postgresql
fi# Leer variables de entorno
source .env 2>/dev/null || {
    echo "Archivo .env no encontrado. Usando valores por defecto."
    DB_HOST=${DB_HOST:-localhost}
    DB_USER=${DB_USER:-postgres}
    DB_PASSWORD=${DB_PASSWORD:-postgres}
    DB_NAME=${DB_NAME:-school_monitoring}
    DB_PORT=${DB_PORT:-5432}
}echo "Configurando base de datos: $DB_NAME"
echo "Usuario: $DB_USER"
echo "Host: $DB_HOST"# Crear base de datos si no existe
sudo -u postgres psql -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || echo "Base de datos ya existe o error al crearla"# Crear usuario si no existe (si no es postgres)
if [ "$DB_USER" != "postgres" ]; then
    sudo -u postgres psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || echo "Usuario ya existe"
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;" 2>/dev/null
fiecho "=== Base de datos configurada ==="
echo "Puedes iniciar el backend con: ./main o go run main.go"
