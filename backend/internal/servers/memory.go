package servers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrNotFound = errors.New("server not found")

type MemoryRepository struct {
	mu      sync.RWMutex
	servers map[ID]ServerStatus
}

func NewMemoryRepository(servers ...ServerStatus) *MemoryRepository {
	items := make(map[ID]ServerStatus, len(servers))
	for _, server := range servers {
		items[server.ServerID] = server
	}
	return &MemoryRepository{servers: items}
}

func (repository *MemoryRepository) List(ctx context.Context) []ServerStatus {
	servers, _ := repository.ListAll(ctx)
	return servers
}

func (repository *MemoryRepository) ListAll(_ context.Context) ([]ServerStatus, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	servers := make([]ServerStatus, 0, len(repository.servers))
	for _, server := range repository.servers {
		servers = append(servers, server)
	}
	return servers, nil
}

func (repository *MemoryRepository) FindByID(_ context.Context, id ID) (ServerStatus, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	server, exists := repository.servers[id]
	if !exists {
		return ServerStatus{}, ErrNotFound
	}
	return server, nil
}

func (repository *MemoryRepository) Register(_ context.Context, request RegistrationRequest) (ServerStatus, error) {
	if err := request.Validate(); err != nil {
		return ServerStatus{}, err
	}
	repository.mu.Lock()
	defer repository.mu.Unlock()
	serverID := ID(fmt.Sprintf("srv_%d", len(repository.servers)+1))
	now := time.Now().UTC()
	status := ServerStatus{ServerID: serverID, Registration: RegistrationPending, Connection: ConnectionUnknown, ObservedAt: now}
	repository.servers[serverID] = status
	return status, nil
}
