package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harbor-panel/harbor-panel/internal/controlplane"
	"github.com/harbor-panel/harbor-panel/internal/domains"
	"github.com/harbor-panel/harbor-panel/internal/persistence"
	"github.com/harbor-panel/harbor-panel/internal/persistence/sqlite"
)

func TestSQLiteRepositoriesComposeWithControlPlane(t *testing.T) {
	database, err := persistence.OpenSQLite(context.Background(), t.TempDir()+"/harbor.db")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()

	accounts, err := sqlite.NewAccountRepository(database)
	if err != nil {
		t.Fatalf("create account repository: %v", err)
	}
	servers, err := sqlite.NewServerRepository(database)
	if err != nil {
		t.Fatalf("create server repository: %v", err)
	}
	domainRepository, err := sqlite.NewDomainRepository(database)
	if err != nil {
		t.Fatalf("create domain repository: %v", err)
	}

	handler := controlplane.NewHandler(accounts, servers, domains.NewService(domainRepository))
	request := httptest.NewRequest(http.MethodGet, "/api/v1/accounts", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 from composed control-plane handler, got %d: %s", recorder.Code, recorder.Body.String())
	}
}
