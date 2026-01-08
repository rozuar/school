.PHONY: dev install build test clean docker-up docker-down

# Development
dev:
	docker-compose up -d db
	cd source/backend && go run ./cmd/server &
	cd source/frontend && npm run dev -- --port 3000 &
	cd source/backoffice && npm run dev -- --port 3001

# Install dependencies
install:
	cd source/backend && go mod download
	cd source/frontend && npm install
	cd source/backoffice && npm install

# Build all
build:
	cd source/backend && go build -o ../../dist/server ./cmd/server
	cd source/frontend && npm run build
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
	cd source/backend && go test ./...

# Clean
clean:
	rm -rf dist
	rm -rf source/frontend/.next
	rm -rf source/backoffice/.next
