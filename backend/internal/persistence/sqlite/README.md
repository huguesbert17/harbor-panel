# SQLite persistence adapter

## Responsibility

This adapter implements the hosting-account and managed-server repository
contracts using the shared database lifecycle from `internal/persistence`.
Domain modules remain independent of SQL and can later receive PostgreSQL or
another adapter without changing their entities and commands.

## Usage

```go
database, err := persistence.OpenSQLite(ctx, "./harbor.db")
if err != nil {
	return err
}
defer database.Close()

accounts, err := sqlite.NewAccountRepository(database)
if err != nil {
	return err
}
servers, err := sqlite.NewServerRepository(database)
if err != nil {
	return err
}
```

`OpenSQLite` runs migrations before returning. Repositories expose the
original collection methods for compatibility and `ListAll` methods for
callers that need to handle database failures explicitly.

## Operational notes

- The adapter uses the pure-Go `modernc.org/sqlite` driver.
- Do not expose the SQLite file to untrusted users.
- Use a file on durable storage for local development; use the future
  PostgreSQL adapter for multi-process or multi-node production deployments.
- Server registration stores the full registration request even though the
  current `ServerStatus` response intentionally exposes only status fields.
