# Go Project Structure Template

A comprehensive template demonstrating best practices for structuring a production-ready Go project. This example includes authentication, user management, email services, database migrations, and more.

## 📁 Project Structure Overview

```
.
├── cmd/                      # Application entry points
│   ├── server/              # Main HTTP server
│   │   └── main.go
│   └── db/                  # Database CLI tools
│       ├── main.go
│       ├── migration_helper.go
│       └── migrations/      # SQL migration files
├── internal/                # Private application code (not importable by other packages)
│   ├── config/             # Configuration management
│   │   ├── env.go
│   │   └── myconstant/
│   ├── handler/            # HTTP handlers and middleware
│   │   ├── middleware/
│   │   └── response/
│   ├── i18n/               # Internationalization support
│   ├── models/             # Data models and enums
│   ├── modules/            # Feature modules (auth, user, etc.)
│   │   ├── auth/
│   │   └── user/
│   ├── pkg/                # Reusable internal packages
│   │   ├── sqlhelper/
│   │   └── validate/
│   ├── server/             # Server setup and initialization
│   └── service/            # Business logic services
├── scripts/                # Utility scripts
├── util/                   # Utility functions
├── go.mod                  # Go module definition
└── README.md              # This file
```

## 🎯 Directory Purpose Guide

### `cmd/` - Command/Application Entry Points

Each subdirectory under `cmd/` represents a standalone executable:

- **`cmd/server/`** - The main HTTP API server. Contains the entry point (`main.go`) that initializes the server, loads configuration, and starts the application.
- **`cmd/db/`** - Database management CLI. Handles migrations and database initialization.

**Why separate?** Allows you to build and run multiple executables from a single codebase.

### `internal/` - Private Application Code

Code in the `internal/` directory cannot be imported by other projects (Go's convention). This is where your core business logic lives.

#### `internal/config/` - Configuration

- Environment variable loading
- Application constants and configuration values

#### `internal/handler/` - HTTP Request Handling

- **`middleware/`** - HTTP middleware (authentication, logging, etc.)
- **`response/`** - Standardized response formatting and error handling

#### `internal/i18n/` - Internationalization

- Multi-language support
- Translation files (en.json, fr.json, vi.json)
- Middleware for language detection

#### `internal/models/` - Data Models

- Common enums and constants
- Shared data structures

#### `internal/modules/` - Feature Modules

Organized by feature/domain:

- **`auth/`** - Authentication logic
  - `dto.go` - Data transfer objects
  - `handler.go` - HTTP handlers
  - `module.go` - Module initialization
  - `repository.go` - Data access layer
  - `util.go` - Helper functions
- **`user/`** - User management
  - Similar structure to auth module

**Best Practice:** Each module is self-contained and can be independently tested.

#### `internal/pkg/` - Reusable Internal Packages

- **`sqlhelper/`** - Database query utilities
- **`validate/`** - Input validation and binding

#### `internal/server/` - Server Setup

- Database initialization
- Environment setup
- HTTP client configuration
- Module registration

#### `internal/service/` - Business Logic Services

- **`emailservice/`** - Email sending and templates
- **`jwt/`** - JWT token management
- **`myredis/`** - Redis client and operations

### `util/` - Global Utilities

Project-wide utility functions (context helpers, array utilities, async operations)

### `scripts/` - Build and Utility Scripts

Scripts for building, deployment, or development automation

## 🏗️ Architectural Patterns

### Module Pattern

Each feature (auth, user) is organized as a module with:

- **Handler** - HTTP request handling
- **Service/Repository** - Business logic and data access
- **DTO** - Request/Response models
- **Module** - Dependency injection and initialization

Example flow:

```
HTTP Request → Handler → Service/Repository → Database → Response
```

### Separation of Concerns

- **Handlers** - Deal with HTTP (requests/responses)
- **Services** - Contain business logic
- **Repositories** - Handle data access
- **Models** - Define data structures

### Middleware Pattern

Middleware functions in `handler/middleware/` wrap HTTP handlers for:

- Authentication validation
- Request logging
- Error handling
- Internationalization

## 🚀 Getting Started

### Prerequisites

- Go 1.19 or higher
- PostgreSQL (or configured database)
- Redis (optional, for session management)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/haidang93/template-golang-project.git
cd example
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

### Running the Application

**Development Mode:**

```bash
go run cmd/server/main.go
```

**Build Release:**

```bash
go build -o server cmd/server/main.go
```

### Database Migrations

```bash
go run cmd/db/main.go up      # Apply migrations
go run cmd/db/main.go down    # Rollback migrations
```

## 📋 Key Features Demonstrated

- ✅ Clean project structure
- ✅ Modular architecture
- ✅ Authentication & Authorization
- ✅ Database migrations
- ✅ Error handling
- ✅ Middleware pattern
- ✅ Internationalization (i18n)
- ✅ Email service integration
- ✅ Redis integration
- ✅ Input validation
- ✅ JWT token management
- ✅ Standardized API responses

## 🔧 Configuration

### Environment Variables

Configuration is managed through environment variables loaded in `internal/config/env.go`.

### Database

- PostgreSQL is used as the primary database
- Migrations are stored in `cmd/db/migrations/`
- Migration helper utilities in `cmd/db/migration_helper.go`

## 📝 Coding Conventions

### File Naming

- Use lowercase with underscores: `file_name.go`
- Suffix files by purpose: `model.go`, `handler.go`, `repository.go`, `interface.go`

### Package Organization

- Keep packages focused on a single responsibility
- Use interfaces for dependency injection
- Group related functionality in the same package

### Error Handling

Standardized error responses defined in `internal/handler/response/err_handling.go`

## 🧪 Testing

Each module can be tested independently:

```bash
go test ./internal/modules/auth/...
go test ./internal/modules/user/...
```

## 📚 Resources

This template demonstrates patterns from:

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## 📄 License

This template is provided as-is for educational purposes.

---

**Note:** This is a template showcasing Go best practices. Adapt it to your project's specific needs.
