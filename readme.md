
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
- GORM + SQLite
- golang-jwt/jwt
- bcrypt for password hashing
- structured logging via slog

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
│   └── db.go
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
├── middleware/
│   └── requireAuth.go
├── shared/
│   └── chef.go
├── api.http
├── go.mod
├── go.sum
├── readme.md
└── .env.example (if configured in your environment)
```

## Runtime Configuration

The service is configured through environment variables:

- HOST: default `0.0.0.0`
- PORT: default `8080`
- JWT_SECRET: default `change-me-in-production`
- DATABASE_PATH: default `database/database.db`
- GIN_MODE: default `debug`
- LOG_LEVEL: default `info`

Example:

```bash
export HOST=0.0.0.0
export PORT=8080
export JWT_SECRET=your-secret
export DATABASE_PATH=database/database.db
export GIN_MODE=debug
export LOG_LEVEL=info
```

## Running the Service

```bash
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
