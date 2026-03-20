
# Recipea API

A recipe-sharing platform API built with Go, allowing users (chefs) to create, share, and discover recipes with a community-driven rating system.

## Overview

Recipea is a learning project to develop proficiency with Go while building a practical recipe-sharing application. The API enables users to manage their own recipes, rate others' creations, and share culinary knowledge with the community.

## Features

- **Chef Management**: User registration and authentication with chef profiles
- jwt auth/access tokens
- **Recipe Creation**: Create and manage recipes with flexible ingredients and step-by-step instructions
- **Recipe Rating**: Thumbs up/down voting system for community feedback
- **Image Uploads**: Support for recipe images and chef profile pictures
- **Flexible Ingredients**: Support for various measurement types (tablespoons, grams, milliliters, whole items, etc.)
  - Example: "1 tbsp honey", "100g spinach", "1 lemon", "75ml dry white wine", "2 sticks celery"
- **Recipe Structure**:
  - **Ingredients Section**: Flexible format with quantities and units
  - **Method Section**: Numbered step-by-step instructions

## Tech Stack

- **Language**: Go
- **Framework**: Gin Web Framework
- **Database**: SQLite with GORM ORM
- **Module**: recipea.com/m

## Project Structure

```
GoRecipeaApi/
├── main.go
├── go.mod
├── readme.md
├── database/
│   └── db.go
├── domain/
│   └── chef/
│       ├── chef.go
│       ├── errors.go
│       └── repository.go
├── application/
│   └── chef/
│       ├── service.go
│       └── update.go
├── infrastructure/
│   ├── auth/
│   │   └── jwt_provider.go
│   └── persistence/
│       └── chef_repository.go
├── interfaces/
│   └── http/
│       ├── chef_handler.go
│       ├── auth_middleware.go
│       └── chef_handler_test.go
├── middleware/
│   └── requireAuth.go
└── shared/
    └── chef.go
```

## DDD migration note

The project has been refactored from legacy `Chef/` route/service code into a DDD-style layered architecture:

- **Domain**: `domain/chef` contains the Chef entity behavior and typed domain errors.
- **Application**: `application/chef` implements use cases like create, update, list, login with explicit request/response objects.
- **Infrastructure**: `infrastructure/persistence` and `infrastructure/auth` provide concrete DB and JWT adapters.
- **Interface / Adapters**: `interfaces/http` implements Gin HTTP controllers and middleware wiring to application services.

The old `Chef/` package used in earlier versions has been removed and is no longer required.

## Prerequisites

- Go 1.16 or higher
- SQLite3

## Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd GoRecipeaApi
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

## Usage

Run the application:

```bash
go run .\main.go
```

The API will start on `localhost:8080`

## API Endpoints

### Base URL

`http://localhost:8080`

### Core Routes

- `GET /` - Welcome/home page

### Chef Routes

See `interfaces/http/chef_handler.go` for detailed endpoint implementation.

## How to extend (DDD)

To add new domain entities (e.g. Recipe aggregate, Rating domain):

1. Add an aggregate in `domain/recipe/` (entity methods, validation, domain errors).
2. Define repository interface in `domain/recipe/repository.go`.
3. Implement repository adapter in `infrastructure/persistence/` (GORM or other DB adapter).
4. Add use-case services in `application/recipe/` (create/update/get/list) that depend only on domain and repository interfaces.
5. Add HTTP handlers in `interfaces/http/` for request/response DTOs and call application services.
6. Wire in `main.go` with `NewRecipeRepository(...)`, `NewRecipeService(...)`, and route registration.

This keeps behavior and business rules in domain/application layers and isolates infra/HTTP details.

## Features in Development

- [ ] Recipe creation, update, and deletion endpoints
- [ ] Recipe rating/voting system
- [ ] Image upload functionality
- [ ] Recipe search and filtering
- [ ] Chef profile management
- [ ] Recipe categories/tags
- [ ] User authentication and authorization

## Database

The application uses SQLite with automatic migration enabled. The database file (`database.db`) is created automatically on first run.

Currently migrated models:

- `Chef` - User/chef profiles

## License

This project is created for learning purposes.
