package hostingaccounts

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("hosting account not found")

type MemoryRepository struct {
	mu       sync.RWMutex
	accounts map[string]Account
}

func NewMemoryRepository(accounts ...Account) *MemoryRepository {
	items := make(map[string]Account, len(accounts))
	for _, account := range accounts {
		items[account.ID] = account
	}
	return &MemoryRepository{accounts: items}
}

func (repository *MemoryRepository) Create(_ context.Context, account Account) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if _, exists := repository.accounts[account.ID]; exists {
		return errors.New("hosting account already exists")
	}
	repository.accounts[account.ID] = account
	return nil
}

func (repository *MemoryRepository) FindByID(_ context.Context, id string) (Account, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	account, exists := repository.accounts[id]
	if !exists {
		return Account{}, ErrNotFound
	}
	return account, nil
}

func (repository *MemoryRepository) List(ctx context.Context) []Account {
	accounts, _ := repository.ListAll(ctx)
	return accounts
}

func (repository *MemoryRepository) ListAll(_ context.Context) ([]Account, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	accounts := make([]Account, 0, len(repository.accounts))
	for _, account := range repository.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}
