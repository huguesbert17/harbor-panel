// Package servers defines the control-plane contracts for managed servers.
package servers

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type ID string

type RegistrationStatus string

const (
	RegistrationPending   RegistrationStatus = "pending"
	RegistrationEnrolling RegistrationStatus = "enrolling"
	RegistrationActive    RegistrationStatus = "active"
	RegistrationDegraded  RegistrationStatus = "degraded"
	RegistrationOffline   RegistrationStatus = "offline"
	RegistrationRevoked   RegistrationStatus = "revoked"
)

type ConnectionStatus string

const (
	ConnectionUnknown      ConnectionStatus = "unknown"
	ConnectionConnected    ConnectionStatus = "connected"
	ConnectionDisconnected ConnectionStatus = "disconnected"
)

type CapabilityName string

const (
	CapabilityWebServer  CapabilityName = "web_server"
	CapabilityPHPRuntime CapabilityName = "php_runtime"
	CapabilityMySQL      CapabilityName = "mysql"
	CapabilityPostgreSQL CapabilityName = "postgresql"
	CapabilityDNS        CapabilityName = "dns"
	CapabilityMail       CapabilityName = "mail"
	CapabilityTLS        CapabilityName = "tls"
	CapabilityBackups    CapabilityName = "backups"
)

type CapabilityStatus string

const (
	CapabilityAvailable   CapabilityStatus = "available"
	CapabilityUnavailable CapabilityStatus = "unavailable"
	CapabilityUnknown     CapabilityStatus = "unknown"
)

type Capability struct {
	Name    CapabilityName   `json:"name"`
	Version string           `json:"version,omitempty"`
	Status  CapabilityStatus `json:"status"`
}

type RegistrationRequest struct {
	Name         string `json:"name"`
	Hostname     string `json:"hostname"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
}

type EnrollmentRequest struct {
	ServerID     ID           `json:"server_id"`
	EnrollmentID string       `json:"enrollment_id"`
	PublicKey    string       `json:"public_key"`
	Proof        string       `json:"proof"`
	AgentVersion string       `json:"agent_version"`
	Capabilities []Capability `json:"capabilities"`
}

type EnrollmentResult struct {
	ServerID   ID                 `json:"server_id"`
	Status     RegistrationStatus `json:"status"`
	Connection ConnectionStatus   `json:"connection_status"`
	AcceptedAt time.Time          `json:"accepted_at"`
}

type ServerStatus struct {
	ServerID     ID                 `json:"server_id"`
	Registration RegistrationStatus `json:"registration_status"`
	Connection   ConnectionStatus   `json:"connection_status"`
	AgentVersion string             `json:"agent_version,omitempty"`
	LastSeenAt   *time.Time         `json:"last_seen_at,omitempty"`
	ObservedAt   time.Time          `json:"observed_at"`
	Capabilities []Capability       `json:"capabilities"`
}

func (request RegistrationRequest) Validate() error {
	if !isRequiredText(request.Name, 128) {
		return invalidField("name", "must be non-empty and at most 128 characters")
	}
	if !isHostname(request.Hostname) {
		return invalidField("hostname", "must be a valid hostname or IP address")
	}
	if !isSafeToken(request.Platform, 64) {
		return invalidField("platform", "must be a safe identifier")
	}
	if !isSafeToken(request.Architecture, 64) {
		return invalidField("architecture", "must be a safe identifier")
	}
	return nil
}

func (request EnrollmentRequest) Validate() error {
	if !isSafeToken(string(request.ServerID), 128) {
		return invalidField("server_id", "must be a safe identifier")
	}
	if !isSafeToken(request.EnrollmentID, 128) {
		return invalidField("enrollment_id", "must be a safe identifier")
	}
	if !isOpaqueValue(request.PublicKey, 8192) {
		return invalidField("public_key", "must be a non-empty value without control characters")
	}
	if !isOpaqueValue(request.Proof, 8192) {
		return invalidField("proof", "must be a non-empty value without control characters")
	}
	if !isSafeToken(request.AgentVersion, 64) {
		return invalidField("agent_version", "must be a safe identifier")
	}
	return validateCapabilities(request.Capabilities)
}

func (status ServerStatus) Validate() error {
	if !isSafeToken(string(status.ServerID), 128) {
		return invalidField("server_id", "must be a safe identifier")
	}
	if !isRegistrationStatus(status.Registration) {
		return invalidField("registration_status", "is not recognized")
	}
	if !isConnectionStatus(status.Connection) {
		return invalidField("connection_status", "is not recognized")
	}
	return validateCapabilities(status.Capabilities)
}

func (status ServerStatus) Supports(name CapabilityName) bool {
	for _, capability := range status.Capabilities {
		if capability.Name == name && capability.Status == CapabilityAvailable {
			return true
		}
	}
	return false
}

func validateCapabilities(capabilities []Capability) error {
	seen := make(map[CapabilityName]struct{}, len(capabilities))
	for index, capability := range capabilities {
		if !isCapabilityName(capability.Name) {
			return invalidField(fmt.Sprintf("capabilities[%d].name", index), "is not recognized")
		}
		if _, exists := seen[capability.Name]; exists {
			return invalidField("capabilities", "must not contain duplicates")
		}
		seen[capability.Name] = struct{}{}
		if !isCapabilityStatus(capability.Status) {
			return invalidField(fmt.Sprintf("capabilities[%d].status", index), "is not recognized")
		}
		if capability.Version != "" && !isSafeToken(capability.Version, 64) {
			return invalidField(fmt.Sprintf("capabilities[%d].version", index), "must be a safe identifier")
		}
	}
	return nil
}

func invalidField(field, message string) error {
	return fmt.Errorf("%s %s", field, message)
}

func isRequiredText(value string, max int) bool {
	return value != "" && len(value) <= max && strings.TrimSpace(value) == value && !strings.ContainsAny(value, "\r\n\x00")
}

func isOpaqueValue(value string, max int) bool {
	return value != "" && len(value) <= max && !strings.ContainsAny(value, "\r\n\x00")
}

func isSafeToken(value string, max int) bool {
	if value == "" || len(value) > max || strings.TrimSpace(value) != value {
		return false
	}
	for _, character := range value {
		if (character < 'a' || character > 'z') &&
			(character < 'A' || character > 'Z') &&
			(character < '0' || character > '9') &&
			character != '-' && character != '_' && character != '.' {
			return false
		}
	}
	return true
}

func isHostname(value string) bool {
	if len(value) == 0 || len(value) > 253 || strings.HasSuffix(value, ".") {
		return net.ParseIP(value) != nil
	}
	if net.ParseIP(value) != nil {
		return true
	}
	for _, label := range strings.Split(value, ".") {
		if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, character := range label {
			if (character < 'a' || character > 'z') &&
				(character < 'A' || character > 'Z') &&
				(character < '0' || character > '9') && character != '-' {
				return false
			}
		}
	}
	return true
}

func isCapabilityName(name CapabilityName) bool {
	switch name {
	case CapabilityWebServer, CapabilityPHPRuntime, CapabilityMySQL, CapabilityPostgreSQL,
		CapabilityDNS, CapabilityMail, CapabilityTLS, CapabilityBackups:
		return true
	default:
		return false
	}
}

func isCapabilityStatus(status CapabilityStatus) bool {
	return status == CapabilityAvailable || status == CapabilityUnavailable || status == CapabilityUnknown
}

func isRegistrationStatus(status RegistrationStatus) bool {
	switch status {
	case RegistrationPending, RegistrationEnrolling, RegistrationActive, RegistrationDegraded,
		RegistrationOffline, RegistrationRevoked:
		return true
	default:
		return false
	}
}

func isConnectionStatus(status ConnectionStatus) bool {
	return status == ConnectionUnknown || status == ConnectionConnected || status == ConnectionDisconnected
}
