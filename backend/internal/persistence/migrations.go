package persistence

import (
	"context"
	"fmt"
	"time"
)

type migration struct {
	version    int
	statements []string
}

var migrations = []migration{
	{
		version: 1,
		statements: []string{
			`CREATE TABLE IF NOT EXISTS hosting_accounts (
				id TEXT PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				status TEXT NOT NULL CHECK (status IN ('active', 'pending', 'suspended')),
				created_at TEXT NOT NULL,
				updated_at TEXT NOT NULL
			)`,
			`CREATE TABLE IF NOT EXISTS managed_servers (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				hostname TEXT NOT NULL,
				platform TEXT NOT NULL,
				architecture TEXT NOT NULL,
				registration_status TEXT NOT NULL CHECK (registration_status IN ('pending', 'enrolling', 'active', 'degraded', 'offline', 'revoked')),
				connection_status TEXT NOT NULL CHECK (connection_status IN ('unknown', 'connected', 'disconnected')),
				agent_version TEXT NOT NULL DEFAULT '',
				last_seen_at TEXT,
				observed_at TEXT NOT NULL,
				capabilities_json TEXT NOT NULL DEFAULT '[]',
				created_at TEXT NOT NULL,
				updated_at TEXT NOT NULL
			)`,
			`CREATE INDEX IF NOT EXISTS idx_hosting_accounts_status ON hosting_accounts (status)`,
			`CREATE INDEX IF NOT EXISTS idx_managed_servers_registration_status ON managed_servers (registration_status)`,
			`CREATE TABLE IF NOT EXISTS schema_migrations (
				version INTEGER PRIMARY KEY,
				applied_at TEXT NOT NULL
			)`,
		},
	},
}

// Migrate applies all pending schema migrations transactionally and is safe to
// call repeatedly when processes start.
func (database *Database) Migrate(ctx context.Context) error {
	if _, err := database.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL
	)`); err != nil {
		return fmt.Errorf("create schema migrations table: %w", err)
	}

	for _, item := range migrations {
		var applied bool
		err := database.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)`, item.version).Scan(&applied)
		if err != nil {
			return fmt.Errorf("check migration %d: %w", item.version, err)
		}
		if applied {
			continue
		}

		transaction, err := database.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %d: %w", item.version, err)
		}
		for _, statement := range item.statements {
			if _, err := transaction.ExecContext(ctx, statement); err != nil {
				_ = transaction.Rollback()
				return fmt.Errorf("apply migration %d: %w", item.version, err)
			}
		}
		if _, err := transaction.ExecContext(ctx, `INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)`, item.version, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
			_ = transaction.Rollback()
			return fmt.Errorf("record migration %d: %w", item.version, err)
		}
		if err := transaction.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", item.version, err)
		}
	}
	return nil
}
