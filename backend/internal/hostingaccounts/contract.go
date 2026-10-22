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

// PersistentRepository is the production-shaped persistence contract. List
// keeps the original development-facing collection shape, while ListAll
// exposes storage failures so callers can report them instead of silently
// returning incomplete data.
type PersistentRepository interface {
	Repository
	List(ctx context.Context) []Account
	ListAll(ctx context.Context) ([]Account, error)
}

type Service interface {
	Create(ctx context.Context, command CreateCommand) (Account, error)
	Get(ctx context.Context, id string) (Account, error)
}
