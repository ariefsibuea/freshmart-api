# Freshmart API

A RESTful product catalog API for a grocery retail. Supports adding and retrieving grocery products with search, filtering, sorting, and pagination.

## Tech Stack

| Component        | Technology                           |
| ---------------- | ------------------------------------ |
| Language         | Go 1.25.9                            |
| HTTP Framework   | Echo v4                              |
| Database         | MySQL 8.4 LTS                        |
| Database Driver  | `go-sql-driver/mysql` (database/sql) |
| Cache            | Redis 8 Alpine                       |
| Migration        | `golang-migrate`                     |
| Config           | `kelseyhightower/envconfig`          |
| Containerization | Docker + Docker Compose v2+          |

## Running Locally

### Prerequisites

- Go 1.25.9+
- `golang-migrate` CLI — required by migration targets ([installation guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- Docker v24+
- Docker Compose v2+

### Steps

1. Clone the repository:

```bash
git clone https://github.com/ariefsibuea/freshmart-api.git
cd freshmart-api
```

2. Copy the environment file:

```bash
cp .env.example .env
```

3. Start the API:

```bash
make api-start
```

This builds the Go image and starts three services: `api`, `mysql`, and `redis`. The API waits for MySQL and Redis to pass their health checks before starting.

4. Run database migrations and seed data:

```bash
make migrate-up
```

This applies all three migrations: table creation, index creation, and seed data.

5. Verify the API is running:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{ "status": "ok" }
```

### Stopping the API

```bash
# Stop containers (preserves volumes)
make api-stop

# Stop and remove all containers, volumes, and images
make api-reset
```

## Documentation

- **API Reference:** [`docs/api.md`](docs/api.md)
