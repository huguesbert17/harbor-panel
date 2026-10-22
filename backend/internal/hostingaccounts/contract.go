// Package hostingaccounts owns hosting-account lifecycle rules. The first
// version exposes contracts only; persistence and provisioning are intentionally
// separate concerns.
package hostingaccounts

import "context"

type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

const (
	StatusActive    = "active"
	StatusPending   = "pending"
	StatusSuspended = "suspended"
)

type CreateCommand struct {
	Username string
}

type Repository interface {
	Create(ctx context.Context, account Account) error
	FindByID(ctx context.Context, id string) (Account, error)
}

type Service interface {
	Create(ctx context.Context, command CreateCommand) (Account, error)
	Get(ctx context.Context, id string) (Account, error)
}
