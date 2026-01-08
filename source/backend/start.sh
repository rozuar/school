#!/bin/bash# Script para iniciar el backendecho "=== Iniciando Backend ==="# Verificar que existe el archivo .env
if [ ! -f .env ]; then
    echo "ERROR: Archivo .env no encontrado"
    echo "Por favor crea el archivo .env con la configuración de la base de datos"
    exit 1
fi# Cargar variables de entorno
export $(cat .env | grep -v '^#' | xargs)# Verificar conexión a la base de datos
echo "Verificando conexión a la base de datos..."
timeout 2 bash -c "echo > /dev/tcp/${DB_HOST:-localhost}/${DB_PORT:-5432}" 2>/dev/null
if [ $? -ne 0 ]; then
    echo "ERROR: No se puede conectar a PostgreSQL en ${DB_HOST:-localhost}:${DB_PORT:-5432}"
    echo "Por favor verifica que PostgreSQL esté corriendo:"
    echo "  sudo systemctl status postgresql"
    echo "  sudo systemctl start postgresql"
    exit 1
fi# Compilar si es necesario
if [ ! -f main ]; then
    echo "Compilando backend..."
    go build -o main .
fi# Iniciar el servidor
echo "Iniciando servidor en puerto ${PORT:-8080}..."
./main
