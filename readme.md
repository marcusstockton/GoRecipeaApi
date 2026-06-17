
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
go run .\cmd\api\main.go
```

The API will start on `localhost:8080`

## API Endpoints

### Base URL

`http://localhost:8080`

### Core Routes

- `GET /` - Welcome/home page

### Chef Routes

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---|
| GET | `/chef` | List all chefs | No |
| GET | `/chef/:id` | Get a specific chef by ID | No |
| POST | `/chef` | Create a new chef (register) | No |
| PUT | `/chef/:id` | Update chef profile | No |
| POST | `/chef/login` | Login and receive JWT token | No |
| GET | `/chef/validate` | Validate JWT token | Yes |

#### Chef Request/Response Examples

**Create Chef (POST /chef)**
```json
{
  "first_name": "Gordon",
  "last_name": "Ramsay",
  "email": "gordon@example.com",
  "password": "securePassword123"
}
```

**Login (POST /chef/login)**
```json
{
  "email": "gordon@example.com",
  "password": "securePassword123"
}
```

**Update Chef (PUT /chef/:id)**
```json
{
  "first_name": "Gordon",
  "last_name": "Ramsay",
  "email": "gordon@example.com",
  "password": "newPassword123"
}
```

### Recipe Routes

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---|
| GET | `/recipe` | List all recipes | No |
| GET | `/recipe/:id` | Get a specific recipe by ID | No |
| POST | `/recipe` | Create a new recipe | Yes |
| PUT | `/recipe/:id` | Update a recipe | Yes |
| DELETE | `/recipe/:id` | Delete a recipe | Yes |

#### Recipe Request/Response Examples

**Create Recipe (POST /recipe)**
```json
{
  "title": "Chocolate Cake",
  "description": "A delicious homemade chocolate cake",
  "ingredients": [
    {
      "quantity": "2",
      "unit": "cups",
      "name": "flour",
      "preparation": "sifted"
    },
    {
      "quantity": "1",
      "unit": "cup",
      "name": "sugar"
    }
  ],
  "steps": [
    {
      "order": 1,
      "action": "Preheat oven to 350°F"
    },
    {
      "order": 2,
      "action": "Mix dry ingredients in a bowl"
    }
  ]
}
```

**Update Recipe (PUT /recipe/:id)**
```json
{
  "title": "Chocolate Cake",
  "description": "A delicious homemade chocolate cake",
  "ingredients": [...],
  "steps": [...]
}
```

### Authentication

Protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

Obtain a token by logging in via `POST /chef/login`.

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

- [x] Recipe creation, update, and deletion endpoints
- [x] Chef profile management
- [x] User authentication and authorization
- [ ] Recipe rating/voting system
- [ ] Image upload functionality
- [ ] Recipe search and filtering
- [ ] Recipe categories/tags

## Database

The application uses SQLite with automatic migration enabled. The database file (`database.db`) is created automatically on first run.

Currently migrated models:

- `Chef` - User/chef profiles

## License

This project is created for learning purposes.
