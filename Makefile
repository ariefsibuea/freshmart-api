COMPOSE := docker compose

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
