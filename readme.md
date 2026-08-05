
# Recipea API

A recipe-sharing platform API built with Go, allowing users (chefs) to create, share, and discover recipes with a community-driven rating system.

## Overview

Recipea is a learning project to develop proficiency with Go while building a practical recipe-sharing application. The API enables users to manage their own recipes, rate others' creations, and share culinary knowledge with the community.

## Features

- **Chef Management**: User registration and authentication with chef profiles
- **JWT auth/access tokens**: Secure access for protected endpoints
- **Recipe Creation**: Create and manage recipes with flexible ingredients and step-by-step instructions
- **Recipe Engagement**: Recipes now support likes and comments from other chefs
  - One chef can like a recipe once
  - Chefs cannot like their own recipes
  - Comments can be added to recipes and replies can be posted to comments via `parent_id`
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
├── go.mod
├── readme.md
├── cmd/
│   └── api/
│       └── main.go
├── database/
│   └── db.go
├── domain/
│   ├── chef/
│   │   ├── chef.go
│   │   ├── errors.go
│   │   └── repository.go
│   └── recipe/
│       ├── recipe.go
│       ├── errors.go
│       └── repository.go
├── application/
│   ├── chef/
│   │   ├── service.go
│   │   └── update.go
│   └── recipe/
│       └── service.go
├── infrastructure/
│   ├── auth/
│   │   └── jwt_provider.go
│   └── persistence/
│       ├── chef_repository.go
│       └── recipe_repository.go
├── interfaces/
│   └── http/
│       ├── auth_middleware.go
│       ├── chef_handler.go
│       ├── chef_handler_test.go
│       └── recipe_handler.go
├── middleware/
│   └── requireAuth.go
└── shared/
    └── chef.go
```

## DDD Architecture

The project follows a Domain-Driven Design (DDD) layered architecture:

- **Domain**: `domain/chef` and `domain/recipe` contain entity behavior, validation rules, and typed domain errors.
- **Application**: `application/chef` and `application/recipe` implement use cases (create, update, list, login, etc.) with explicit request/response DTOs.
- **Infrastructure**: `infrastructure/persistence` provides database adapters using GORM; `infrastructure/auth` provides JWT authentication.
- **Interface / Adapters**: `interfaces/http` implements Gin HTTP handlers and middleware, wiring requests to application services.

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
go run ./cmd/api/main.go
```

Or from the cmd/api directory:

```bash
cd cmd/api && go run main.go
```

The API will start on `localhost:8080`

## API Endpoints

### Base URL

`http://localhost:8080`

### Core Routes

- `GET /` - Welcome/home page
- `POST /chef/login` - Chef authentication
- `GET /chef/validate` - Validate JWT token (requires auth)

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
| POST | `/recipe/:id/like` | Like a recipe | Yes |
| DELETE | `/recipe/:id/like` | Remove your like from a recipe | Yes |
| POST | `/recipe/:id/comment` | Add a comment or reply to a comment | Yes |
| GET | `/recipe/:id/comments` | List comments for a recipe | No |
| DELETE | `/comment/:commentID` | Delete your own comment | Yes |

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

**Recipe response example**
```json
{
  "id": 1,
  "title": "Chocolate Cake",
  "description": "A delicious homemade chocolate cake",
  "owned_by": 1,
  "ingredients": [...],
  "steps": [...],
  "like_count": 3
}
```

**Add a comment (POST /recipe/:id/comment)**
```json
{
  "content": "This recipe looks amazing!"
}
```

**Reply to a comment (POST /recipe/:id/comment)**
```json
{
  "content": "Thanks, I appreciate the feedback!",
  "parent_id": 1
}
```

**List comments (GET /recipe/:id/comments)**
```json
{
  "comments": [
    {
      "id": 1,
      "recipe_id": 1,
      "chef_id": 2,
      "content": "This recipe looks amazing!"
    }
  ]
}
```

### Authentication

Protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

Obtain a token by logging in via `POST /chef/login`.

## How to Extend (DDD Pattern)

To add new domain entities (e.g., Rating, Comment, Category aggregates):

1. **Define the Domain** - Create an aggregate in `domain/{entity}/` with:
   - `{entity}.go` - Entity struct, value objects, and business logic
   - `repository.go` - Repository interface (contracts, not implementation)
   - `errors.go` - Domain-specific errors

2. **Implement Repositories** - Create adapter in `infrastructure/persistence/{entity}_repository.go`:
   - Implements the repository interface using GORM
   - Handles database operations
   - Translates between domain models and database models

3. **Add Use Cases** - Create services in `application/{entity}/`:
   - Implements business operations (create, update, list, delete, etc.)
   - Depends only on domain and repository interfaces
   - Contains application-specific logic and DTOs

4. **Expose via HTTP** - Create handlers in `interfaces/http/{entity}_handler.go`:
   - Defines request/response DTOs
   - Handles HTTP concerns (status codes, headers)
   - Calls application services

5. **Wire It Together** - Update `cmd/api/main.go`:
   - Create repository instance: `NewRecipeRepository(db)`
   - Create service instance: `NewRecipeService(repository)`
   - Register HTTP routes and handlers

This keeps domain behavior isolated, makes testing easier, and maintains a clear separation of concerns.

## Implementation Status

### Implemented
- [x] Chef profile management (create, read, update)
- [x] Chef authentication (registration and JWT-based login)
- [x] Recipe CRUD operations (create, read, update, delete)
- [x] Recipe ingredient and step structure
- [x] Recipe likes and comments
- [x] Comment replies via parent comment IDs
- [x] Authorization middleware for protected endpoints
- [x] DDD layered architecture (domain, application, infrastructure, interfaces)

### In Development
- [ ] Image upload functionality (recipes & chef profiles)
- [ ] Recipe search and filtering
- [ ] Recipe categories/tags
- [ ] Notification system for comments and likes

## Database

The application uses SQLite with automatic migration enabled. The database file (`database.db`) is created automatically on first run.

Currently migrated models:

- `Chef` - User/chef profiles
- `Recipe` - Recipe documents with ingredients and steps
- `RecipeLike` - Tracks which chefs liked which recipes
- `RecipeComment` - Stores recipe comments and reply threads

## License

This project is created for learning purposes.
