package agent

import (
	"errors"
	"testing"
)

func TestValidateRequestAcceptsKnownOperation(t *testing.T) {
	err := ValidateRequest(Request{
		ID:             "request-1",
		Operation:      OperationCreateWebsite,
		IdempotencyKey: "website-1",
		Parameters: map[string]string{
			"account_id":    "account-1",
			"domain":        "example.com",
			"document_root": "/srv/harbor/accounts/account-1/public",
			"php_version":   "8.3",
		},
	})
	if err != nil {
		t.Fatalf("expected request to be valid, got %v", err)
	}
}

func TestValidateRequestRejectsUnknownParameter(t *testing.T) {
	err := ValidateRequest(Request{
		ID:             "request-1",
		Operation:      OperationHealthCheck,
		IdempotencyKey: "health-1",
		Parameters:     map[string]string{"command": "id"},
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected invalid request error, got %v", err)
	}
}

func TestValidateRequestRejectsUnsafePathAndControlCharacters(t *testing.T) {
	unsafeRequests := []Request{
		{
			ID:             "request-1",
			Operation:      OperationCreateWebsite,
			IdempotencyKey: "website-1",
			Parameters: map[string]string{
				"account_id":    "account-1",
				"domain":        "example.com",
				"document_root": "/srv/harbor/../etc",
			},
		},
		{
			ID:             "request-2",
			Operation:      OperationIssueCertificate,
			IdempotencyKey: "certificate-1",
			Parameters: map[string]string{
				"domain": "example.com\nreload",
			},
		},
	}
	for _, request := range unsafeRequests {
		if err := ValidateRequest(request); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("expected invalid request error, got %v", err)
		}
	}
}

func TestValidateRequestRequiresIdempotencyKey(t *testing.T) {
	err := ValidateRequest(Request{
		ID:        "request-1",
		Operation: OperationHealthCheck,
	})
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected invalid request error, got %v", err)
	}
}

func TestValidateRequestRejectsUnsupportedOperation(t *testing.T) {
	err := ValidateRequest(Request{
		ID:             "request-1",
		Operation:      Operation("execute_shell"),
		IdempotencyKey: "request-1",
	})
	if !errors.Is(err, ErrUnsupportedOperation) {
		t.Fatalf("expected unsupported operation error, got %v", err)
	}
}
