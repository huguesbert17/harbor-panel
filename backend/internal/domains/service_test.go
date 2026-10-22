package domains

import "testing"

func TestDNSRecordValidation(t *testing.T) {
	if err := (CreateDNSRecordCommand{Name: "@", Type: "A", TTL: 300, Value: "192.0.2.10"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (CreateDNSRecordCommand{Name: "@", Type: "A", TTL: 30, Value: "192.0.2.10"}).Validate(); err == nil {
		t.Fatal("expected TTL validation error")
	}
}

func TestDomainNamesRejectUnsafeInput(t *testing.T) {
	if validDomain("example.com;rm") {
		t.Fatal("expected unsafe domain to be rejected")
	}
	if !validDomain("Example.com") {
		t.Fatal("expected valid domain")
	}
}
