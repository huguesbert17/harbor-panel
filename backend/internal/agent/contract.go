// Package agent defines the typed control-plane to managed-server boundary.
package agent

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidRequest       = errors.New("invalid agent request")
	ErrUnsupportedOperation = errors.New("unsupported agent operation")
)

type Operation string

const (
	OperationHealthCheck      Operation = "health_check"
	OperationCreateAccount    Operation = "create_hosting_account"
	OperationCreateWebsite    Operation = "create_website"
	OperationIssueCertificate Operation = "issue_certificate"
	OperationRemoveWebsite    Operation = "remove_website"
)

type OperationStatus string

const (
	OperationAccepted  OperationStatus = "accepted"
	OperationSucceeded OperationStatus = "succeeded"
	OperationFailed    OperationStatus = "failed"
	OperationRejected  OperationStatus = "rejected"
)

type Request struct {
	ID             string            `json:"id"`
	Operation      Operation         `json:"operation"`
	IdempotencyKey string            `json:"idempotency_key"`
	Parameters     map[string]string `json:"parameters"`
}

type Result struct {
	ID        string            `json:"id"`
	Operation Operation         `json:"operation"`
	Status    OperationStatus   `json:"status"`
	Message   string            `json:"message,omitempty"`
	Observed  map[string]string `json:"observed,omitempty"`
}

type Client interface {
	Execute(ctx context.Context, request Request) (Result, error)
}

// ValidateRequest checks the operation envelope and its allowlisted parameters.
// It deliberately does not execute, authorize, or persist anything.
func ValidateRequest(request Request) error {
	if !isSafeToken(request.ID, 128) {
		return invalidField("id", "must be a non-empty safe identifier")
	}
	if !isSafeToken(request.IdempotencyKey, 128) {
		return invalidField("idempotency_key", "must be a non-empty safe identifier")
	}

	spec, ok := operationSpecs[request.Operation]
	if !ok {
		return fmt.Errorf("%w: %q", ErrUnsupportedOperation, request.Operation)
	}

	for key, value := range request.Parameters {
		if !spec.allows(key) {
			return invalidField("parameters."+key, "parameter is not allowed for this operation")
		}
		if err := validateParameterValue(key, value); err != nil {
			return err
		}
	}
	for _, key := range spec.required {
		value, present := request.Parameters[key]
		if !present || strings.TrimSpace(value) == "" {
			return invalidField("parameters."+key, "parameter is required")
		}
	}

	return validateOperationParameters(request.Operation, request.Parameters)
}

type operationSpec struct {
	required []string
	optional []string
}

func (s operationSpec) allows(parameter string) bool {
	for _, allowed := range s.required {
		if parameter == allowed {
			return true
		}
	}
	for _, allowed := range s.optional {
		if parameter == allowed {
			return true
		}
	}
	return false
}

var operationSpecs = map[Operation]operationSpec{
	OperationHealthCheck: {},
	OperationCreateAccount: {
		required: []string{"account_id", "username"},
		optional: []string{"quota_bytes"},
	},
	OperationCreateWebsite: {
		required: []string{"account_id", "domain", "document_root"},
		optional: []string{"php_version"},
	},
	OperationIssueCertificate: {
		required: []string{"domain"},
		optional: []string{"contact_email"},
	},
	OperationRemoveWebsite: {
		required: []string{"account_id", "domain"},
	},
}

var versionPattern = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+){1,2}$`)

func validateOperationParameters(operation Operation, parameters map[string]string) error {
	for _, key := range []string{"account_id", "username"} {
		if value, ok := parameters[key]; ok && !isSafeToken(value, 128) {
			return invalidField("parameters."+key, "must be a safe identifier")
		}
	}
	if value, ok := parameters["domain"]; ok && !isHostname(value) {
		return invalidField("parameters.domain", "must be a valid hostname")
	}
	if value, ok := parameters["document_root"]; ok && !isSafeAbsolutePath(value) {
		return invalidField("parameters.document_root", "must be a safe absolute path")
	}
	if value, ok := parameters["quota_bytes"]; ok {
		quota, err := strconv.ParseUint(value, 10, 64)
		if err != nil || quota == 0 {
			return invalidField("parameters.quota_bytes", "must be a positive unsigned integer")
		}
	}
	if value, ok := parameters["php_version"]; ok && !versionPattern.MatchString(value) {
		return invalidField("parameters.php_version", "must be a version number")
	}
	if value, ok := parameters["contact_email"]; ok && !isContactEmail(value) {
		return invalidField("parameters.contact_email", "must be a valid contact email")
	}

	// Keep the switch explicit so adding an operation requires an intentional
	// validation decision instead of silently accepting arbitrary parameters.
	switch operation {
	case OperationHealthCheck, OperationCreateAccount, OperationCreateWebsite,
		OperationIssueCertificate, OperationRemoveWebsite:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedOperation, operation)
	}
}

func validateParameterValue(field, value string) error {
	if value == "" || len(value) > 4096 || strings.IndexFunc(value, func(r rune) bool {
		return r == '\x00' || r == '\r' || r == '\n'
	}) >= 0 {
		return invalidField("parameters."+field, "must be a non-empty scalar without control characters")
	}
	return nil
}

func invalidField(field, message string) error {
	return fmt.Errorf("%w: %s %s", ErrInvalidRequest, field, message)
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

func isSafeAbsolutePath(value string) bool {
	if value == "" || !strings.HasPrefix(value, "/") || value == "/" || strings.ContainsRune(value, '\x00') {
		return false
	}
	cleaned := path.Clean(value)
	return cleaned == value && !strings.HasPrefix(strings.TrimPrefix(value, "/"), "../")
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

func isContactEmail(value string) bool {
	if strings.Count(value, "@") != 1 || strings.ContainsAny(value, " \t") {
		return false
	}
	parts := strings.SplitN(value, "@", 2)
	return isSafeToken(parts[0], 128) && isHostname(parts[1])
}
