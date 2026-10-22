package servers

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEnrollmentRejected = errors.New("agent enrollment rejected")
	ErrHeartbeatRejected  = errors.New("agent heartbeat rejected")
)

// AgentStore is the persistence boundary required by AgentService. Concrete
// database or remote implementations belong in adapters outside this module.
// Implementations must enforce one-time enrollment, replay protection, and
// server authorization in their own security-aware layer.
type AgentStore interface {
	Enroll(ctx context.Context, request EnrollmentRequest, acceptedAt time.Time) (EnrollmentResult, error)
	Heartbeat(ctx context.Context, request HeartbeatRequest, acceptedAt time.Time) (HeartbeatResult, error)
}

// AgentService validates the bounded agent protocol and delegates state
// changes to an explicitly supplied store. It never executes commands,
// authenticates callers, issues credentials, or selects a transport.
type AgentService struct {
	store AgentStore
	clock func() time.Time
}

func NewAgentService(store AgentStore) *AgentService {
	return &AgentService{store: store, clock: func() time.Time { return time.Now().UTC() }}
}

// Enroll validates an agent enrollment envelope before handing it to the
// store. Proof verification and credential issuance are intentionally outside
// this scaffold and must not be inferred from validation alone.
func (service *AgentService) Enroll(ctx context.Context, request EnrollmentRequest) (EnrollmentResult, error) {
	if service == nil || service.store == nil {
		return EnrollmentResult{}, ErrEnrollmentRejected
	}
	if err := request.Validate(); err != nil {
		return EnrollmentResult{}, err
	}
	result, err := service.store.Enroll(ctx, request, service.now())
	if err != nil {
		return EnrollmentResult{}, errors.Join(ErrEnrollmentRejected, err)
	}
	return result, nil
}

// Heartbeat validates an observed status report before handing it to the
// store. The store remains responsible for rejecting unknown or revoked
// servers and applying freshness/replay policy.
func (service *AgentService) Heartbeat(ctx context.Context, request HeartbeatRequest) (HeartbeatResult, error) {
	if service == nil || service.store == nil {
		return HeartbeatResult{}, ErrHeartbeatRejected
	}
	if err := request.Validate(); err != nil {
		return HeartbeatResult{}, err
	}
	result, err := service.store.Heartbeat(ctx, request, service.now())
	if err != nil {
		return HeartbeatResult{}, errors.Join(ErrHeartbeatRejected, err)
	}
	return result, nil
}

func (service *AgentService) now() time.Time {
	if service.clock == nil {
		return time.Now().UTC()
	}
	return service.clock().UTC()
}
