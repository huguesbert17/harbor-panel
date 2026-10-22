// Package persistence owns database lifecycle and schema migration primitives.
// Domain-specific repositories live in adapter subpackages.
package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// Database wraps the control-plane SQL connection. SQLite is the first local
// adapter; keeping this wrapper small leaves room for a PostgreSQL adapter
// without leaking driver details into domain modules.
type Database struct {
	db *sql.DB
}

// OpenSQLite opens a SQLite database and initializes its schema. A path of
// ":memory:" is supported for isolated tests. The single connection keeps
// SQLite in-memory databases stable and is sufficient for the local adapter;
// production multi-node storage will use a server database adapter later.
func OpenSQLite(ctx context.Context, path string) (*Database, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("sqlite path must not be empty")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	database := &Database{db: db}
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable sqlite foreign keys: %w", err)
	}
	if _, err := db.ExecContext(ctx, "PRAGMA busy_timeout = 5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("configure sqlite busy timeout: %w", err)
	}
	if err := database.Migrate(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return database, nil
}

// SQL exposes the connection to repository adapters inside the persistence
// boundary. It is intentionally not exposed by domain packages.
func (database *Database) SQL() *sql.DB {
	return database.db
}

// Close releases the underlying database connection.
func (database *Database) Close() error {
	if database == nil || database.db == nil {
		return nil
	}
	return database.db.Close()
}
