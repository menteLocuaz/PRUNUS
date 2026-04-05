# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server (port defaults to 9090)
go run ./cmd/ serve --port 9090

# Run migrations
go run ./cmd/ migrate

# Seed initial data (empresa, rol, usuario root)
go run ./cmd/ register

# Build binary
go build -o prunus ./cmd/

# Tests
go test ./...
go test ./pkg/services -run TestCreateUsuario   # single test

# Format / vet
go fmt ./...
go vet ./...
go mod tidy
```

## Architecture

Clean Architecture with strict layering: **Store → Service → Handler → Router**

```
cmd/           CLI entry (Cobra): serve, migrate, register subcommands
cmd/app.go     Dependency injection — all stores/services/handlers wired here
pkg/
  models/      Domain entities (JSON tags, soft delete fields)
  dto/          Request structs with `validate:` tags; response shapes
  store/        Repository pattern — StoreXXX interfaces + pgx/v5 SQL impl
  services/     Business logic, caching orchestration, logging
  transport/http/  HTTP handlers (parse DTO → call service → response.XXX)
  routers/      Chi route registration; main_router.go mounts all sub-routers
  middleware/   Auth (RequireAuth/OptionalAuth), CORS, rate limit, logger, etc.
  helper/       JWT generation/validation, bcrypt
  utils/        response/ helpers, validator/FormatErrors, pagination, performance tracing
  config/       Viper env loading, DB/Redis init
  config/database/migrations/  Numbered SQL migration files (001–N)
```

### Adding a new resource (checklist)

1. Migration in `pkg/config/database/migrations/0XX_name.go`, register in `migrations.go`
2. Model in `pkg/models/`
3. DTOs in `pkg/dto/` with `validate:` tags
4. `StoreXXX` interface + impl in `pkg/store/`
5. Service in `pkg/services/`
6. Handler in `pkg/transport/http/`
7. Router in `pkg/routers/` (mount with `r.Mount(...)` inside `/api/v1`)
8. Wire everything in `cmd/app.go` `RegisterHandlers()`

### Key patterns

**Response helpers** — always use `pkg/utils/response`:
```go
response.Success(w, "mensaje", data)
response.Created(w, "mensaje", data)
response.BadRequest(w, "mensaje")
response.NotFound(w, "mensaje")
response.ValidationError(w, validator.FormatErrors(err))
response.InternalServerError(w, err.Error())
```

**Validation** — DTOs use `go-playground/validator/v10`; handler pattern:
```go
if err := validator.Validate.Struct(req); err != nil {
    response.ValidationError(w, validator.FormatErrors(err))
    return
}
```

**Pagination** — cursor-based; call `utils.ParsePaginationParams(r)` → `dto.PaginationParams` with `Limit`, `LastID`, `LastDate`. Query params: `limit`, `last_id`, `last_date` (RFC3339).

**JWT context keys** injected by `RequireAuth()`:
- `"user_id"` → `uuid.UUID`
- `"user_email"` → `string`
- `"user_rol"` → `string`
- `"user_sucursal"` → `uuid.UUID`

**Caching** — Services use the `CacheStore` interface (Redis). The system degrades gracefully if Redis is unavailable; never make cache availability a hard requirement.

**Performance tracing** — wrap expensive operations:
```go
defer performance.Trace(ctx, "operation-name", performance.DBThreshold)()
```

**Soft delete** — all queries must filter `WHERE deleted_at IS NULL`. Delete operations set `deleted_at = CURRENT_TIMESTAMP`. Never physically delete rows.

### Environment variables (via Viper / `.env`)

Required: `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET` (min 32 chars)  
Optional: `DB_PORT` (5432), `DB_SSLMODE`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `PORT` (9090), `JWT_EXPIRATION_HOURS` (24)

## Conventions (see also AGENTS.md)

- Error messages and code comments in **Spanish**
- JSON tags in `snake_case`; optional fields use `omitempty`
- `StoreXXX` naming for repository interfaces
- All tables include `created_at`, `updated_at`, `deleted_at`
- Passwords must be stripped from structs before returning in responses
- SQL written as inline strings with named parameters (pgx `@param` syntax); never string-concatenated queries
