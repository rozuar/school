.PHONY: dev install build test clean docker-up docker-down

# Development
dev:
	docker-compose up -d db
	cd source/api && go run ./cmd/server &
	cd source/web && npm run dev -- --port 3000 &
	cd source/backoffice && npm run dev -- --port 3001

# Development (sin Docker: asume DATABASE_URL/DB_* ya configuradas)
dev-nodocker:
	cd source/api && go run ./cmd/server &
	cd source/web && npm run dev -- --port 3000 &
	cd source/backoffice && npm run dev -- --port 3001

# Install dependencies
install:
	cd source/api && go mod download
	cd source/web && npm install
	cd source/backoffice && npm install

# Build all
build:
	cd source/api && go build -o ../../dist/server ./cmd/server
	cd source/web && npm run build
	cd source/backoffice && npm run build

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

# Database
db-init:
	psql $(DATABASE_URL) -f scripts/init-db.sql

db-reset:
	psql $(DATABASE_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	psql $(DATABASE_URL) -f scripts/init-db.sql

# Test
test:
	cd source/api && go test ./...

# Clean
clean:
	rm -rf dist
	rm -rf source/web/.next
	rm -rf source/backoffice/.next
