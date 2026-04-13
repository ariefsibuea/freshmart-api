COMPOSE := docker compose

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: sync-deps
sync-deps:
	@if ! command -v go >/dev/null 2>&1; then \
		echo "Error: Go is not installed. Please install Go first: https://go.dev/dl/"; \
		exit 1; \
	fi
	@if [ ! -f go.mod ] || [ ! -f go.sum ]; then \
		echo "Error: go.mod or go.sum not found. Run 'go mod init' first."; \
		exit 1; \
	fi
	go mod tidy

.PHONY: api-start
api-start: sync-deps
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy .env.example to .env and configure it!"; \
		exit 1; \
	fi
	$(COMPOSE) up --build -d

.PHONY: api-stop
api-stop:
	$(COMPOSE) down

.PHONY: api-reset
api-reset:
	$(COMPOSE) down -v
	docker image prune -f
	docker volume prune -f

MIGRATE := migrate
MIGRATIONS_PATH := ./internal/migrations
# Use localhost for migrations from host machine (MySQL port is exposed on host)
DB_URL := mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp(localhost:$(MYSQL_PORT))/$(MYSQL_DATABASE)

.PHONY: migrate-up
migrate-up:
	@if ! command -v $(MIGRATE) >/dev/null 2>&1; then \
		echo "Error: golang-migrate is not installed. Install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	@if ! command -v $(MIGRATE) >/dev/null 2>&1; then \
		echo "Error: golang-migrate is not installed. Install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

.PHONY: migrate-down-all
migrate-down-all:
	@if ! command -v $(MIGRATE) >/dev/null 2>&1; then \
		echo "Error: golang-migrate is not installed. Install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@echo "WARNING: This will rollback ALL migrations and drop all tables!"
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down
