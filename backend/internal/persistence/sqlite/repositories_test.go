package sqlite

import (
	"context"
	"errors"
	"testing"

	"github.com/harbor-panel/harbor-panel/internal/hostingaccounts"
	"github.com/harbor-panel/harbor-panel/internal/persistence"
	"github.com/harbor-panel/harbor-panel/internal/servers"
)

func TestAccountRepositoryPersistsAcrossDatabaseReopen(t *testing.T) {
	path := t.TempDir() + "/harbor.db"
	ctx := context.Background()

	database, err := persistence.OpenSQLite(ctx, path)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	accounts, err := NewAccountRepository(database)
	if err != nil {
		t.Fatalf("create account repository: %v", err)
	}
	wanted := hostingaccounts.Account{ID: "acct_001", Username: "harbor", Status: hostingaccounts.StatusActive}
	if err := accounts.Create(ctx, wanted); err != nil {
		t.Fatalf("create account: %v", err)
	}
	if err := accounts.Create(ctx, wanted); !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}

	database, err = persistence.OpenSQLite(ctx, path)
	if err != nil {
		t.Fatalf("reopen database: %v", err)
	}
	defer database.Close()
	accounts, err = NewAccountRepository(database)
	if err != nil {
		t.Fatalf("recreate account repository: %v", err)
	}
	actual, err := accounts.FindByID(ctx, wanted.ID)
	if err != nil {
		t.Fatalf("find account after reopen: %v", err)
	}
	if actual != wanted {
		t.Fatalf("expected %#v, got %#v", wanted, actual)
	}
}

func TestServerRepositoryPersistsRegistration(t *testing.T) {
	database, err := persistence.OpenSQLite(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	serversRepository, err := NewServerRepository(database)
	if err != nil {
		t.Fatalf("create server repository: %v", err)
	}

	created, err := serversRepository.Register(context.Background(), servers.RegistrationRequest{
		Name:         "web-01",
		Hostname:     "web-01.example.com",
		Platform:     "debian",
		Architecture: "arm64",
	})
	if err != nil {
		t.Fatalf("register server: %v", err)
	}
	if created.Registration != servers.RegistrationPending {
		t.Fatalf("expected pending registration, got %s", created.Registration)
	}
	items, err := serversRepository.ListAll(context.Background())
	if err != nil {
		t.Fatalf("list servers: %v", err)
	}
	if len(items) != 1 || items[0].ServerID != created.ServerID {
		t.Fatalf("expected registered server in list, got %#v", items)
	}
	found, err := serversRepository.FindByID(context.Background(), created.ServerID)
	if err != nil {
		t.Fatalf("find server: %v", err)
	}
	if found.ServerID != created.ServerID || found.ObservedAt.IsZero() {
		t.Fatalf("unexpected server record: %#v", found)
	}
}

func TestRepositoriesRespectCanceledContext(t *testing.T) {
	database, err := persistence.OpenSQLite(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	accounts, err := NewAccountRepository(database)
	if err != nil {
		t.Fatalf("create account repository: %v", err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := accounts.FindByID(canceled, "missing"); err == nil {
		t.Fatal("expected canceled context error")
	}
}
