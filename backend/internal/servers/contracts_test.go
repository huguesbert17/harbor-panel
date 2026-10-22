package servers

import (
	"testing"
	"time"
)

func TestRegistrationRequestValidate(t *testing.T) {
	request := RegistrationRequest{
		Name:         "edge-1",
		Hostname:     "panel.example.com",
		Platform:     "debian",
		Architecture: "amd64",
	}
	if err := request.Validate(); err != nil {
		t.Fatalf("expected registration to be valid, got %v", err)
	}
}

func TestEnrollmentRequestValidate(t *testing.T) {
	request := EnrollmentRequest{
		ServerID:     "server-1",
		EnrollmentID: "enrollment-1",
		PublicKey:    "opaque-public-key",
		Proof:        "opaque-proof",
		AgentVersion: "0.1.0",
		Capabilities: []Capability{{
			Name:    CapabilityWebServer,
			Version: "nginx-1.27",
			Status:  CapabilityAvailable,
		}},
	}
	if err := request.Validate(); err != nil {
		t.Fatalf("expected enrollment to be valid, got %v", err)
	}
}

func TestEnrollmentRequestRejectsDuplicateCapabilities(t *testing.T) {
	request := EnrollmentRequest{
		ServerID:     "server-1",
		EnrollmentID: "enrollment-1",
		PublicKey:    "key",
		Proof:        "proof",
		AgentVersion: "0.1.0",
		Capabilities: []Capability{
			{Name: CapabilityTLS, Status: CapabilityAvailable},
			{Name: CapabilityTLS, Status: CapabilityAvailable},
		},
	}
	if err := request.Validate(); err == nil {
		t.Fatal("expected duplicate capabilities to be rejected")
	}
}

func TestServerStatusSupportsAvailableCapabilitiesOnly(t *testing.T) {
	status := ServerStatus{
		ServerID:     "server-1",
		Registration: RegistrationActive,
		Connection:   ConnectionConnected,
		ObservedAt:   time.Now(),
		Capabilities: []Capability{
			{Name: CapabilityWebServer, Status: CapabilityAvailable},
			{Name: CapabilityMail, Status: CapabilityUnavailable},
		},
	}
	if !status.Supports(CapabilityWebServer) {
		t.Fatal("expected available capability to be supported")
	}
	if status.Supports(CapabilityMail) {
		t.Fatal("did not expect unavailable capability to be supported")
	}
	if err := status.Validate(); err != nil {
		t.Fatalf("expected status to be valid, got %v", err)
	}
}

func TestHeartbeatRequestValidateAcceptsBoundedObservation(t *testing.T) {
	request := HeartbeatRequest{
		ServerID:     "server-1",
		AgentVersion: "0.1.0",
		ObservedAt:   time.Now().UTC(),
		Capabilities: []Capability{{Name: CapabilityWebServer, Status: CapabilityAvailable}},
	}
	if err := request.Validate(); err != nil {
		t.Fatalf("expected bounded heartbeat to be valid, got %v", err)
	}
}
