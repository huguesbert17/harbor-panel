# Persistence boundary

## Responsibility

This package owns database lifecycle, SQLite configuration, and transactional
schema migrations. Domain-specific repositories are implemented by adapters
under `persistence/` and depend on domain contracts rather than the reverse.

## Current adapter

`persistence/sqlite` is the local persistent adapter. It uses the pure-Go
`modernc.org/sqlite` driver, so the backend and future server agent do not
require CGO or a system SQLite development package.

The adapter supports a file path for local development and `:memory:` for
isolated tests. It currently uses one connection to keep SQLite behavior
deterministic. A future PostgreSQL adapter will provide the production
multi-node deployment path.

## Migrations

`OpenSQLite` initializes the schema and applies pending migrations before
returning. Migrations are versioned in `migrations.go`, run in transactions,
and recorded in `schema_migrations`. New migrations must be additive and have
adapter tests for existing data and rollback-sensitive behavior.

## API wiring

`cmd/api` opens one database during process startup, applies migrations, and
constructs both SQLite repositories before registering the control-plane
routes. The database is closed during process shutdown. The control-plane
handler still depends only on the domain persistence interfaces, so this
composition can later select another adapter without changing HTTP handlers.

The API currently uses `HARBOR_DB_PATH`, defaulting to `./harbor.db`. This is
appropriate for local development and single-process operation; it is not yet
a production multi-node storage strategy.
