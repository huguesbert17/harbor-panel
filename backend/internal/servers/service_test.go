package servers

import (
	"context"
	"testing"
	"time"
)

type recordingAgentStore struct {
	enrollRequest    EnrollmentRequest
	heartbeatRequest HeartbeatRequest
	acceptedAt       time.Time
}

func (store *recordingAgentStore) Enroll(_ context.Context, request EnrollmentRequest, acceptedAt time.Time) (EnrollmentResult, error) {
	store.enrollRequest = request
	store.acceptedAt = acceptedAt
	return EnrollmentResult{ServerID: request.ServerID, Status: RegistrationEnrolling, Connection: ConnectionConnected, AcceptedAt: acceptedAt}, nil
}

func (store *recordingAgentStore) Heartbeat(_ context.Context, request HeartbeatRequest, acceptedAt time.Time) (HeartbeatResult, error) {
	store.heartbeatRequest = request
	store.acceptedAt = acceptedAt
	return HeartbeatResult{ServerID: request.ServerID, Registration: RegistrationActive, Connection: ConnectionConnected, AcceptedAt: acceptedAt}, nil
}

func TestAgentServiceDelegatesValidatedEnrollmentAndHeartbeat(t *testing.T) {
	store := &recordingAgentStore{}
	service := NewAgentService(store)
	service.clock = func() time.Time { return time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC) }

	enrollment := EnrollmentRequest{ServerID: "server-1", EnrollmentID: "enrollment-1", PublicKey: "key", Proof: "proof", AgentVersion: "0.1.0"}
	if _, err := service.Enroll(context.Background(), enrollment); err != nil {
		t.Fatalf("expected enrollment to be delegated, got %v", err)
	}
	if store.enrollRequest.ServerID != enrollment.ServerID || store.acceptedAt.IsZero() {
		t.Fatal("expected enrollment request and acceptance time to reach store")
	}

	heartbeat := HeartbeatRequest{ServerID: "server-1", AgentVersion: "0.1.0", ObservedAt: time.Now().UTC()}
	if _, err := service.Heartbeat(context.Background(), heartbeat); err != nil {
		t.Fatalf("expected heartbeat to be delegated, got %v", err)
	}
	if store.heartbeatRequest.ServerID != heartbeat.ServerID {
		t.Fatal("expected heartbeat request to reach store")
	}
}

func TestAgentServiceRejectsInvalidHeartbeatBeforeStore(t *testing.T) {
	store := &recordingAgentStore{}
	service := NewAgentService(store)

	_, err := service.Heartbeat(context.Background(), HeartbeatRequest{ServerID: "server-1", AgentVersion: "0.1.0"})
	if err == nil {
		t.Fatal("expected invalid heartbeat to fail")
	}
	if store.heartbeatRequest.ServerID != "" {
		t.Fatal("expected invalid heartbeat not to reach store")
	}
}
