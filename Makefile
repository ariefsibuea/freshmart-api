.PHONY: sync-deps
sync-deps:
	@if [ ! -f go.mod ]; then \
		echo "Error: go.mod not found. Run 'go mod init' first."; \
		exit 1; \
	fi
	go mod tidy

.PHONY: api-start
api-start: sync-deps
	docker-compose up --build -d api mysql redis

.PHONY: api-stop
api-stop:
	docker-compose down

.PHONY: api-clean
api-clean:
	docker-compose down -v
	docker image prune -f
	docker volume prune -f
