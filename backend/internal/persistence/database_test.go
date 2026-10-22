package persistence

import (
	"context"
	"database/sql"
	"testing"
)

func TestOpenSQLiteInitializesSchemaAndMigrationsAreIdempotent(t *testing.T) {
	database, err := OpenSQLite(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(context.Background()); err != nil {
		t.Fatalf("run migrations a second time: %v", err)
	}

	var version int
	if err := database.SQL().QueryRow(`SELECT MAX(version) FROM schema_migrations`).Scan(&version); err != nil {
		t.Fatalf("read schema version: %v", err)
	}
	if version != 1 {
		t.Fatalf("expected schema version 1, got %d", version)
	}

	for _, table := range []string{"hosting_accounts", "managed_servers"} {
		var name string
		if err := database.SQL().QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&name); err != nil {
			if err == sql.ErrNoRows {
				t.Fatalf("expected table %q", table)
			}
			t.Fatalf("inspect table %q: %v", table, err)
		}
	}
}
