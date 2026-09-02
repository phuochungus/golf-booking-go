# GORM + Goose Persistence Migration Design

**Date:** 2026-09-02  
**Status:** Approved in chat; pending written-spec review

## Goal

Move the application's runtime database access from generated sqlc code to GORM v2 while keeping Goose as the sole schema-migration tool. Existing MySQL schema, connection-pool settings, and Facility CRUD behavior must remain compatible.

## Scope

In scope:

- GORM v2 MySQL initialization and connection-pool configuration.
- GORM model definitions for every table currently represented by sqlc.
- A GORM-backed Facility repository exposing create, filtered list, get-by-ID, and partial update operations.
- Soft-delete behavior using GORM's `DeletedAt` field while preserving the existing nullable `deleted_at` columns.
- Repository tests that run without an external MySQL or Redis service.
- Removing obsolete sqlc source/config/generated files once no production code references them.

Out of scope:

- Replacing Goose migrations with `AutoMigrate`.
- Adding new HTTP endpoints, authentication, authorization flows, or business rules.
- Changing the SQL schema or migration history.
- Refactoring Redis, logging, configuration loading, or router setup beyond database type wiring.

## Alternatives Considered

1. **GORM runtime + Goose migrations (recommended):** GORM owns models and queries; Goose remains the versioned, reviewable schema authority. This matches the requested toolchain and avoids destructive or implicit production schema changes.
2. **GORM runtime + AutoMigrate:** simpler local setup, but schema changes become implicit and can diverge from the existing Goose history; rejected for this migration.
3. **Compatibility adapter around sqlc:** keeps generated query signatures while delegating some calls to GORM, but leaves sqlc concepts and generated artifacts in the architecture; rejected because the goal is to move runtime persistence to GORM.

## Architecture

### Database initialization

`internal/initialize/mysql.go` builds the same MySQL DSN, calls `gorm.Open(mysql.Open(dsn), ...)`, and stores the resulting `*gorm.DB` in `global.DB`. It retrieves the underlying `*sql.DB` via `db.DB()` to apply the existing idle/open/lifetime settings and to ping the server. Initialization errors are logged through the existing logger and cause the same panic-based startup failure behavior.

`global.DB` is the only process-wide database handle. Goose continues to run through the existing Makefile commands and is not invoked automatically by application startup.

### Models

`internal/database/models.go` defines GORM structs for `organizations`, `facilities`, `admins`, `authz_objects`, `authz_actions`, `authz_permissions`, `authz_roles`, `authz_role_permissions`, `authz_admin_roles`, and `authz_admin_permissions`.

Each model uses an integer primary key mapped to `id`, explicit `gorm:"column:..."` tags where needed, and `CreatedAt`, `UpdatedAt`, and `DeletedAt gorm.DeletedAt` fields mapped to the existing timestamp columns. Nullable foreign keys such as `facility_id` use pointer integer fields. Table names are explicit through `TableName()` methods or tags so pluralization cannot change the existing schema.

### Facility repository

`internal/database/facility.go` owns the `FacilityRepository` type and accepts a `*gorm.DB` in its constructor. It provides:

- `Create(ctx, facility) error`
- `List(ctx, filter) ([]Facility, error)` where optional upper/lower rating and organization filters are applied only when present.
- `GetByID(ctx, id) (Facility, error)` and returns GORM's `ErrRecordNotFound` for a missing or soft-deleted row.
- `Update(ctx, id, changes) error` where only non-nil change fields are written, preserving existing values for omitted fields.

All methods use `WithContext(ctx)`, return the original GORM error, and rely on GORM's default soft-delete scope. No raw SQL is required for these operations.

### Transaction boundary

The repository operates on any `*gorm.DB`, including a transaction clone returned by `db.Begin()`. No global transaction manager is introduced in this migration.

## Data flow

Startup loads configuration → GORM opens MySQL → underlying `*sql.DB` receives pool settings and ping → `global.DB` is available to repositories. A request-level caller constructs a repository from `global.DB`, passes a context and typed input, and receives a model or error. Goose migration commands remain an independent deployment step before startup.

## Error handling

- Startup open, pool, or ping errors use the existing logger and panic behavior.
- Repository methods return GORM errors unchanged so callers can distinguish `gorm.ErrRecordNotFound` from driver errors.
- Empty update input is rejected with a stable package error rather than issuing an unbounded update.
- No method silently ignores database errors.

## Testing strategy

Use GORM's SQLite dialector with an in-memory database for repository tests. The test setup creates only the `facilities` table shape needed by the repository, including `deleted_at`, and enables GORM logger silence. Tests cover:

1. create and generated ID;
2. list with no filters and each optional filter combination;
3. get-by-ID success and not-found behavior;
4. partial update preserving omitted fields;
5. soft-delete exclusion from normal list/get queries.

Initialization tests use a small seam around the GORM opener or inspect pool settings through a real in-memory/sqlmock-compatible handle; they must not require a running MySQL server.

## Compatibility and cleanup

Goose files under `migrations/` and Makefile targets remain unchanged. `sqlc/query.sql`, `sqlc/schema.sql`, `sqlc.yaml`, and generated `internal/database/db.go`/`query.sql.go` artifacts are removed or replaced only after repository call sites are updated. The final tree must contain no production import of `database/sql` for the global DB handle and no sqlc-generated package comments.

## Acceptance criteria

- `go test ./...` passes without external services.
- `go vet ./...` passes.
- Application database initialization compiles with `global.DB` as `*gorm.DB` and preserves configured pool values.
- Facility CRUD behavior and soft-delete semantics are covered by tests.
- Goose remains the documented and executable migration path; no AutoMigrate call exists.
