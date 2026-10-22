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

## Wiring note

The current API still constructs concrete in-memory repositories. A later API
composition change should open one database at startup and inject the SQLite
repositories into the control-plane handler. This task intentionally does not
modify `cmd/api` or authentication code.
