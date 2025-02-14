include .example.envrc

MIGRATION_PATH = ./cmd/migrate/migrations

.PHONY: generate
generate:
	@templ generate

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: up
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "All Docker images are started"

.PHONY: up_dbs
up_dbs:
	@echo "Starting only DB-related Docker containers..."
	docker-compose up -d postgres redis
	@echo "Database services are running"

.PHONY: up_build
up_build:
	@echo "Building and starting Docker images..."
	docker-compose up --build -d
	@echo "All Docker images are started"

.PHONY: down
down:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo "Done!"

.PHONY: build
build:
	@echo "Building..."
	go build -o bin/api-server ./cmd/api-server
	@echo "Built!"

.PHONY: run
run: build
	@echo "Starting the backend server..."
	ADDR=${ADDR} DB_ADDR=${DB_ADDR} MAX_OPEN_CONNS=${MAX_OPEN_CONNS} DB_MAX_IDLE_CONNS=${DB_MAX_IDLE_CONNS} MAX_OPEN_CONNS=${MAX_OPEN_CONNS} ./bin/api-server
	@echo "Server is running!"

.PHONY: setup
setup: up migrate-up seed
	@echo "Database migrations applied, seed data inserted, and Docker containers started!"
