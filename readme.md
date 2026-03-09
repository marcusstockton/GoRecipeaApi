
# Recipea API

A recipe-sharing platform API built with Go, allowing users (chefs) to create, share, and discover recipes with a community-driven rating system.

## Overview

Recipea is a learning project to develop proficiency with Go while building a practical recipe-sharing application. The API enables users to manage their own recipes, rate others' creations, and share culinary knowledge with the community.

## Features

- **Chef Management**: User registration and authentication with chef profiles
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
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── readme.md            # This file
└── Chef/
    ├── models.go        # Chef data model and GORM definitions
    └── routes.go        # Chef API routes
```

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
See `Chef/routes.go` for detailed endpoint documentation

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