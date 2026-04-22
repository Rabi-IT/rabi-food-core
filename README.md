# Rabi Food Core

Backend monolith for the Rabi Food platform — a multi-tenant food delivery system.

[![Go Version](https://img.shields.io/badge/Go-1.26-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Stack

- **Runtime**: Go 1.26
- **HTTP**: Fiber v2
- **Database**: PostgreSQL via pgx/v5
- **Auth**: GoTrue (Supabase)
- **DI**: samber/do
- **Logs**: zerolog (Wide Event pattern)
- **Validation**: go-playground/validator
- **Tests**: testify suite + httpexpect
- **Observability**: Prometheus + Grafana + Loki + Alloy

## Architecture

Feature-based Modular Monolith. Each domain is self-contained with its own controller, gateway, use cases, and routes. Dependencies between features are unidirectional and intentional — never circular.

## Project Structure

```
rabi-food-core/
├── api/                        # Go application
│   ├── config/
│   ├── app_context/            # Session (UserSession)
│   ├── domain/                 # Shared types (auth/role, payment_status)
│   ├── features/               # Business modules (auth, category, order, product, subscription, tenant)
│   │   └── [feature]/
│   │       ├── controller/     # HTTP handlers (fiber_*.go)
│   │       ├── gateway/        # Data access (interface.go + pgx_*.go)
│   │       ├── model/
│   │       ├── routes/
│   │       └── usecases/       # Business logic + integration_test.go
│   ├── fixtures/               # Integration test helpers
│   └── libs/
│       ├── database/           # PgxAdapter, migrations, pagination
│       ├── di/                 # Dependency injection wiring
│       ├── errs/               # Domain errors per feature
│       ├── http/               # FiberAdapter + middlewares
│       ├── logger/             # Wide Event pattern
│       └── validator/
├── infra/
│   ├── api/                    # Docker Compose for API + GoTrue
│   ├── monitoring/             # Grafana, Prometheus, Loki, Alloy
│   └── terraform/              # AWS infrastructure (ECS, RDS, ECR)
└── Taskfile.yml
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [Docker](https://docs.docker.com/get-docker/) + Docker Compose
- [Task](https://taskfile.dev/installation/)

### Running locally

```bash
# Start Postgres + GoTrue + pgAdmin
task db

# Run pending migrations and seed
task migrate

# Start the API (loads .env.test)
task dev
```

### Running tests

```bash
task test
```

## Local services

| Service | URL |
|---------|-----|
| API | http://localhost:3000 |
| API Docs (Swagger) | http://localhost:3000/docs |
| GoTrue (auth) | http://localhost:9999 |
| pgAdmin | http://localhost:5050 |
| Web | http://app.localhost:5173 |

## Commands

| Command | Description |
|---------|-------------|
| `task dev` | Start development environment |
| `task web` | Start web development server |
| `task db` | Start Postgres + GoTrue + pgAdmin (Docker) |
| `task db-down` | Stop and remove containers |
| `task test` | Run integration tests in Docker |
| `task migrate` | Run pending migrations |
| `task migrate-status` | Show migration status |
| `task migrate-down` | Rollback last migration |
| `task lint` | Run golangci-lint |
| `task docs` | Generate Swagger docs |
| `task mockgen` | Generate mocks |
| `task logs` | Tail API container logs |
| `task clean_docker` | Prune all Docker resources |

## Deployment

Deployed to AWS ECS via GitHub Actions on every GitHub Release. The pipeline runs tests, builds and pushes a Docker image to ECR, runs database migrations, then updates the ECS service.

## License

MIT
