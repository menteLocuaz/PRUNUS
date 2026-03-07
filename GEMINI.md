# Prunus - Business Management API

Prunus is a Go-based REST API designed for comprehensive business management, including companies, branches, users, products, and more. It follows a Clean Architecture pattern to ensure scalability and maintainability.

## Project Overview

- **Main Technology:** Go (1.25.4)
- **Architecture:** Layered / Clean Architecture
  - `cmd/`: Application entry point and dependency injection.
  - `pkg/models/`: Domain entities and data models.
  - `pkg/dto/`: Data Transfer Objects for API requests/responses.
  - `pkg/store/`: Data access layer (Repository pattern) using raw SQL.
  - `pkg/services/`: Business logic, validations, and coordination.
  - `pkg/transport/http/`: HTTP handlers and controllers.
  - `pkg/routers/`: Route definitions and middleware setup.
  - `pkg/middleware/`: Authentication, logging, and CORS.
- **Database:** PostgreSQL (v15) with automatic migrations on startup.
- **Authentication:** JWT-based stateless authentication with `bcrypt` for password hashing.

## Key Commands

### Setup & Infrastructure
- **Environment:** `cp .env.example .env` (Edit with local DB and JWT secrets).
- **Database (Docker):** `docker-compose up -d`.
- **Dependencies:** `go mod download` and `go mod tidy`.

### Development
- **Run Server:** `go run cmd/main.go` (Default port: `9090`).
- **Build:** `go build -o prunus cmd/main.go`.
- **Test:** `go test ./...` (Use `-cover` for coverage).
- **Format:** `go fmt ./...`.

## Development Conventions

### Coding Standards
- **Naming:**
  - **Structs:** `PascalCase` (e.g., `Usuario`, `EmpresaHandler`).
  - **Functions:** `PascalCase` for exported, `camelCase` for private.
  - **Interfaces:** Prefix with `Store` for repositories (e.g., `StoreUsuario`).
- **JSON Tags:** Use `snake_case` (e.g., `json:"id_usuario"`). Include `omitempty` for optional fields.
- **Language:** Error messages and comments should be in **Spanish**.

### Implementation Patterns
- **Dependency Injection:** Manual DI in `cmd/main.go` using the constructor pattern (`New...` functions).
- **Database:**
  - Use `pgx/v5` driver.
  - **Soft Deletes:** Always use `deleted_at` field. Update queries should set `deleted_at = NOW()`, and select queries must include `WHERE deleted_at IS NULL`.
  - **Auditory:** All tables should include `created_at`, `updated_at`, and `deleted_at`.
- **Context Usage:** `middleware.RequireAuth` populates the request context with:
  - `user_id`, `user_email`, `user_rol`, `user_sucursal`.

### Error Handling
- Return explicit errors in Spanish from the service layer.
- Use `http.Error` or JSON responses in handlers with appropriate status codes (400 for bad requests, 401/403 for auth issues, 404 for missing resources, 500 for internal errors).

### Security
- Never expose passwords in JSON responses.
- Use parameterized SQL queries to prevent injection.
- Validate all user input in the service layer.
