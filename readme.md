
# Recipea API

A Go-based recipe microservice built with a domain-driven design (DDD) style architecture. The service supports chef registration, JWT authentication, recipe management, likes, and threaded comments.

## Overview

Recipea is a small backend service for managing recipes and chef profiles. The codebase is intentionally structured to demonstrate clean separation between:

- domain logic
- application use cases
- infrastructure adapters
- HTTP transport concerns

The project is designed as a learning-oriented reference for DDD-style microservice architecture in Go.

## Current Architecture

This project follows a pragmatic DDD layering pattern:

- domain: business entities, validation rules, repository contracts, domain errors
- application: orchestration and use cases
- infrastructure: persistence and JWT providers
- interfaces: Gin HTTP handlers, middleware, and request/response binding

## Technology Stack

- Go
- Gin web framework
- GORM + PostgreSQL
- golang-jwt/jwt
- bcrypt for password hashing
- structured logging via slog
- Docker Compose for local service orchestration

## Project Structure

```text
GoRecipeaApi/
├── application/
│   ├── auth/
│   │   └── service.go
│   ├── chef/
│   │   ├── service.go
│   │   └── update.go
│   └── recipe/
│       └── service.go
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   └── config.go
├── database/
│   ├── db.go
│   └── db_test.go
├── domain/
│   ├── chef/
│   │   ├── chef.go
│   │   ├── errors.go
│   │   └── repository.go
│   └── recipe/
│       ├── errors.go
│       ├── recipe.go
│       ├── recipe_test.go
│       └── repository.go
├── infrastructure/
│   ├── auth/
│   │   └── jwt_provider.go
│   ├── logging/
│   │   └── logger.go
│   └── persistence/
│       ├── chef_repository.go
│       └── recipe_repository.go
├── interfaces/
│   └── http/
│       ├── auth_middleware.go
│       ├── chef_handler.go
│       ├── chef_handler_test.go
│       ├── errors.go
│       ├── health.go
│       ├── middleware.go
│       ├── recipe_handler.go
│       ├── recipe_handler_test.go
│       └── auth_middleware_test.go
├── shared/
│   └── chef.go
├── .env
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── api.http
├── go.mod
├── go.sum
├── readme.md
└── .gitignore
```

## Runtime Configuration

This service uses Postgres as the runtime database for all non-test environments. SQLite is still used in local test scenarios only when needed, but the application itself is configured around a Postgres DSN.

Recommended environment variables:

- APP_ENV: `dev`, `uat`, or `prod`
- HOST: default `0.0.0.0`
- PORT: default `8080`
- JWT_SECRET: secret used to sign JWTs
- DATABASE_DSN: Postgres connection string
- GIN_MODE: defaults to `debug` for `dev`, `release` for `uat` and `prod`
- LOG_LEVEL: defaults to `debug` for `dev`, `info` for `uat`, `warn` for `prod`

Example `.env.example`:

```env
APP_ENV=dev
HOST=0.0.0.0
PORT=8080
JWT_SECRET=change-me-in-dev
DATABASE_DSN=host=localhost user=recipea password=recipea dbname=recipea_dev port=5432 sslmode=disable
GIN_MODE=debug
LOG_LEVEL=debug
```

Use a real `.env` file locally, but do not commit secrets. In UAT and production, prefer your deployment platform's environment variables or a secrets manager.

## Running the Service

### Local development with Docker Compose

Start the dev stack:

```bash
docker compose --profile dev up --build
```

This starts:

- Postgres on `localhost:5432`
- API on `http://localhost:8080`
- database name: `recipea_dev`

Stop the dev stack:

```bash
docker compose --profile dev down
```

Stop and remove the dev database volume:

```bash
docker compose --profile dev down -v
```

Check the running containers:

```bash
docker compose --profile dev ps
```

Connect to the database from the host:

```bash
docker exec -it recipea-postgres-dev psql -U recipea -d recipea_dev
```

Query the seeded data:

```bash
docker exec -it recipea-postgres-dev psql -U recipea -d recipea_dev -c "SELECT * FROM recipes;"
```

### UAT

```bash
docker compose --profile uat up --build
```

- Postgres: port `5433`
- API: port `8081`
- DB name: `recipea_uat`

### Production

```bash
docker compose --profile prod up --build
```

- Postgres: port `5434`
- API: port `8082`
- DB name: `recipea_prod`

### Local Go run outside Docker

```bash
set -a
. ./.env
set +a
go run ./cmd/api/main.go
```

The server binds to the configured host and port, defaulting to:

```text
http://0.0.0.0:8080
```

## Health and Readiness

The service exposes health endpoints for operational monitoring:

- `GET /health`
- `GET /ready`

These are useful for load balancers, container health checks, and service discovery.

## API Endpoints

### Root

- `GET /` - welcome response
- `GET /health` - database health check
- `GET /ready` - readiness check

### Chef Endpoints

| Method | Route | Description | Auth |
|---|---|---|---|
| GET | `/chef/` | list chefs | No |
| GET | `/chef/:id` | fetch a chef | No |
| POST | `/chef/` | create chef | No |
| PUT | `/chef/:id` | update chef profile | No |
| POST | `/chef/login` | login and receive JWT | No |
| GET | `/chef/validate` | validate active token | Yes |

### Recipe Endpoints

| Method | Route | Description | Auth |
|---|---|---|---|
| GET | `/recipe/` | list recipes | No |
| GET | `/recipe/:id` | fetch recipe by id | No |
| GET | `/recipe/:id/comments` | list comments | No |
| POST | `/recipe/` | create recipe | Yes |
| PUT | `/recipe/:id` | update recipe | Yes |
| DELETE | `/recipe/:id` | delete recipe | Yes |
| POST | `/recipe/:id/like` | like a recipe | Yes |
| DELETE | `/recipe/:id/like` | remove a like | Yes |
| POST | `/recipe/:id/comment` | add comment or reply | Yes |
| DELETE | `/recipe/comment/:commentID` | delete comment | Yes |

## Authentication

Protected routes require a JWT bearer token:

```http
Authorization: Bearer <token>
```

### Example login request

```json
{
  "email": "gordon@example.com",
  "password": "securePassword123"
}
```

### Example response

```json
{
  "token": "<jwt>",
  "chef": {
    "id": 1,
    "first_name": "Gordon",
    "last_name": "Ramsay",
    "email": "gordon@example.com"
  }
}
```

## Domain Rules Implemented

The domain currently enforces:

- chef password hashing and validation
- unique chef email
- recipe title validation
- recipe owner authorization rules
- recipe duplicate title prevention per owner
- one like per chef per recipe
- self-like prevention
- comment content validation
- nested comment reply validation

## Development Notes

This service is intentionally structured as a demonstration of DDD principles in Go, but it remains a pragmatic implementation rather than a heavy enterprise framework stack. It is well suited as a learning reference or a foundational microservice skeleton.

## Database Migrations

At the moment, the service uses GORM `AutoMigrate` during startup. That is fine for early-stage local development, but it is not a long-term migration strategy for shared environments because schema changes are applied implicitly and are harder to review, roll back, or audit.

### Recommended future approach

Use versioned SQL migration files managed by a tool such as `golang-migrate/migrate`.

Typical pattern:

1. Create a migration file with a name like:
   - `001_create_chefs_table.up.sql`
   - `001_create_chefs_table.down.sql`
2. Run migrations in CI/CD or during deployment before the app starts.
3. Keep migration history in version control.
4. Never change a past migration; add a new migration instead.

Example flow:

```bash
migrate -database "$DATABASE_DSN" -path ./migrations up
```

Then keep the app startup focused on booting the service rather than changing the schema on the fly.

### Practical recommendation for this repo

- keep `AutoMigrate` only while the project is in active prototyping
- add a `migrations/` folder once the schema stabilizes
- use a migration tool in UAT and prod pipelines
- test migration scripts in a disposable database before release

## Testing

The project includes unit and HTTP-level tests for domain behavior and route authorization flows.

Run tests:

```bash
go test ./...
```

## Known Limitations / Next Improvements

The current codebase is a strong architectural skeleton, but it is not yet a full production deployment. Potential next steps include:

- richer error classification for API responses
- rate limiting and request throttling
- OpenAPI/Swagger spec
- graceful shutdown handling
- CI pipeline and linting
- containerization and deployment manifests
- DB migrations instead of auto-migrate in startup
- structured observability and metrics export

## License

This project is intended for learning, demonstration, and backend architecture exploration.
