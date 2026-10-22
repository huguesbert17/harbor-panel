// Package sqlite implements persistent domain repositories on the shared
// SQLite database owned by internal/persistence.
package sqlite

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/harbor-panel/harbor-panel/internal/hostingaccounts"
	"github.com/harbor-panel/harbor-panel/internal/persistence"
	"github.com/harbor-panel/harbor-panel/internal/servers"
)

var (
	ErrAlreadyExists = errors.New("record already exists")
)

// AccountRepository persists hosting accounts.
type AccountRepository struct {
	database *persistence.Database
}

var _ hostingaccounts.PersistentRepository = (*AccountRepository)(nil)

func NewAccountRepository(database *persistence.Database) (*AccountRepository, error) {
	if database == nil || database.SQL() == nil {
		return nil, errors.New("persistence database is required")
	}
	return &AccountRepository{database: database}, nil
}

func (repository *AccountRepository) Create(ctx context.Context, account hostingaccounts.Account) error {
	if account.ID == "" || account.Username == "" {
		return errors.New("hosting account id and username are required")
	}
	if !validAccountStatus(account.Status) {
		return fmt.Errorf("invalid hosting account status %q", account.Status)
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := repository.database.SQL().ExecContext(ctx, `
		INSERT INTO hosting_accounts (id, username, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`, account.ID, account.Username, account.Status, now, now)
	if err != nil {
		if isConstraintError(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("create hosting account: %w", err)
	}
	return nil
}

func (repository *AccountRepository) FindByID(ctx context.Context, id string) (hostingaccounts.Account, error) {
	var account hostingaccounts.Account
	err := repository.database.SQL().QueryRowContext(ctx, `
		SELECT id, username, status FROM hosting_accounts WHERE id = ?`, id).
		Scan(&account.ID, &account.Username, &account.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return hostingaccounts.Account{}, hostingaccounts.ErrNotFound
	}
	if err != nil {
		return hostingaccounts.Account{}, fmt.Errorf("find hosting account: %w", err)
	}
	return account, nil
}

func (repository *AccountRepository) List(ctx context.Context) []hostingaccounts.Account {
	accounts, _ := repository.ListAll(ctx)
	return accounts
}

func (repository *AccountRepository) ListAll(ctx context.Context) ([]hostingaccounts.Account, error) {
	rows, err := repository.database.SQL().QueryContext(ctx, `
		SELECT id, username, status FROM hosting_accounts ORDER BY username, id`)
	if err != nil {
		return nil, fmt.Errorf("list hosting accounts: %w", err)
	}
	defer rows.Close()

	accounts := make([]hostingaccounts.Account, 0)
	for rows.Next() {
		var account hostingaccounts.Account
		if err := rows.Scan(&account.ID, &account.Username, &account.Status); err != nil {
			return nil, fmt.Errorf("scan hosting account: %w", err)
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate hosting accounts: %w", err)
	}
	return accounts, nil
}

// ServerRepository persists managed-server registration and status records.
type ServerRepository struct {
	database *persistence.Database
}

var _ servers.PersistentRepository = (*ServerRepository)(nil)

func NewServerRepository(database *persistence.Database) (*ServerRepository, error) {
	if database == nil || database.SQL() == nil {
		return nil, errors.New("persistence database is required")
	}
	return &ServerRepository{database: database}, nil
}

func (repository *ServerRepository) Register(ctx context.Context, request servers.RegistrationRequest) (servers.ServerStatus, error) {
	if err := request.Validate(); err != nil {
		return servers.ServerStatus{}, err
	}

	serverID, err := newServerID()
	if err != nil {
		return servers.ServerStatus{}, fmt.Errorf("generate server id: %w", err)
	}
	now := time.Now().UTC()
	status := servers.ServerStatus{
		ServerID:     servers.ID(serverID),
		Registration: servers.RegistrationPending,
		Connection:   servers.ConnectionUnknown,
		ObservedAt:   now,
		Capabilities: []servers.Capability{},
	}
	capabilities, err := json.Marshal(status.Capabilities)
	if err != nil {
		return servers.ServerStatus{}, fmt.Errorf("encode server capabilities: %w", err)
	}
	_, err = repository.database.SQL().ExecContext(ctx, `
		INSERT INTO managed_servers (
			id, name, hostname, platform, architecture, registration_status,
			connection_status, agent_version, last_seen_at, observed_at,
			capabilities_json, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		serverID, request.Name, request.Hostname, request.Platform, request.Architecture,
		string(status.Registration), string(status.Connection), status.AgentVersion,
		nil, now.Format(time.RFC3339Nano), string(capabilities), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		if isConstraintError(err) {
			return servers.ServerStatus{}, ErrAlreadyExists
		}
		return servers.ServerStatus{}, fmt.Errorf("register server: %w", err)
	}
	return status, nil
}

func (repository *ServerRepository) List(ctx context.Context) []servers.ServerStatus {
	items, _ := repository.ListAll(ctx)
	return items
}

func (repository *ServerRepository) ListAll(ctx context.Context) ([]servers.ServerStatus, error) {
	rows, err := repository.database.SQL().QueryContext(ctx, serverSelect+` ORDER BY name, id`)
	if err != nil {
		return nil, fmt.Errorf("list managed servers: %w", err)
	}
	defer rows.Close()

	items := make([]servers.ServerStatus, 0)
	for rows.Next() {
		item, err := scanServerStatus(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate managed servers: %w", err)
	}
	return items, nil
}

func (repository *ServerRepository) FindByID(ctx context.Context, id servers.ID) (servers.ServerStatus, error) {
	item, err := scanServerStatus(repository.database.SQL().QueryRowContext(ctx, serverSelect+` WHERE id = ?`, string(id)))
	if errors.Is(err, sql.ErrNoRows) {
		return servers.ServerStatus{}, servers.ErrNotFound
	}
	if err != nil {
		return servers.ServerStatus{}, err
	}
	return item, nil
}

const serverSelect = `
	SELECT id, registration_status, connection_status, agent_version,
		last_seen_at, observed_at, capabilities_json
	FROM managed_servers`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanServerStatus(row rowScanner) (servers.ServerStatus, error) {
	var (
		id, registration, connection, agentVersion, observedAt, capabilitiesJSON string
		lastSeenAt                                                               sql.NullString
	)
	if err := row.Scan(&id, &registration, &connection, &agentVersion, &lastSeenAt, &observedAt, &capabilitiesJSON); err != nil {
		return servers.ServerStatus{}, fmt.Errorf("scan managed server: %w", err)
	}
	capabilities := make([]servers.Capability, 0)
	if err := json.Unmarshal([]byte(capabilitiesJSON), &capabilities); err != nil {
		return servers.ServerStatus{}, fmt.Errorf("decode managed server capabilities: %w", err)
	}
	observed, err := time.Parse(time.RFC3339Nano, observedAt)
	if err != nil {
		return servers.ServerStatus{}, fmt.Errorf("decode managed server observed time: %w", err)
	}
	var lastSeen *time.Time
	if lastSeenAt.Valid {
		value, err := time.Parse(time.RFC3339Nano, lastSeenAt.String)
		if err != nil {
			return servers.ServerStatus{}, fmt.Errorf("decode managed server last-seen time: %w", err)
		}
		lastSeen = &value
	}
	status := servers.ServerStatus{
		ServerID:     servers.ID(id),
		Registration: servers.RegistrationStatus(registration),
		Connection:   servers.ConnectionStatus(connection),
		AgentVersion: agentVersion,
		LastSeenAt:   lastSeen,
		ObservedAt:   observed,
		Capabilities: capabilities,
	}
	if err := status.Validate(); err != nil {
		return servers.ServerStatus{}, fmt.Errorf("validate managed server: %w", err)
	}
	return status, nil
}

func newServerID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "srv_" + hex.EncodeToString(bytes), nil
}

func validAccountStatus(status string) bool {
	switch status {
	case hostingaccounts.StatusActive, hostingaccounts.StatusPending, hostingaccounts.StatusSuspended:
		return true
	default:
		return false
	}
}

func isConstraintError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "constraint") || strings.Contains(message, "unique")
}
