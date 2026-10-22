// Package domains owns hosted domains and DNS desired state.
package domains

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	ErrNotFound       = errors.New("domain resource not found")
	ErrAlreadyExists  = errors.New("domain resource already exists")
	ErrInvalidRequest = errors.New("invalid domain request")
)

type Domain struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	Name         string    `json:"name"`
	Kind         string    `json:"kind"`
	Status       string    `json:"status"`
	DocumentRoot string    `json:"document_root,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const (
	KindPrimary   = "primary"
	KindAddon     = "addon"
	KindAlias     = "alias"
	StatusActive  = "active"
	StatusPending = "pending"
	StatusPaused  = "paused"
)

type DNSZone struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	DomainID  string    `json:"domain_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type DNSRecord struct {
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	TTL    int    `json:"ttl"`
	Value  string `json:"value"`
}
type CreateDomainCommand struct {
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
}
type CreateZoneCommand struct {
	AccountID string `json:"account_id"`
	DomainID  string `json:"domain_id"`
}
type CreateDNSRecordCommand struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	TTL   int    `json:"ttl"`
	Value string `json:"value"`
}
type Repository interface {
	ListDomains(context.Context, string) ([]Domain, error)
	ListZones(context.Context, string) ([]DNSZone, error)
	ListRecords(context.Context, string) ([]DNSRecord, error)
	CreateRecord(context.Context, string, CreateDNSRecordCommand) (DNSRecord, error)
}
type PersistentRepository interface {
	Repository
	CreateDomain(context.Context, Domain) error
	FindDomainByID(context.Context, string) (Domain, error)
	CreateZone(context.Context, DNSZone) error
	FindZoneByID(context.Context, string) (DNSZone, error)
}
type Service interface {
	CreateDomain(context.Context, CreateDomainCommand) (Domain, error)
	ListDomains(context.Context, string) ([]Domain, error)
	CreateZone(context.Context, CreateZoneCommand) (DNSZone, error)
	ListZones(context.Context, string) ([]DNSZone, error)
	ListRecords(context.Context, string) ([]DNSRecord, error)
	CreateRecord(context.Context, string, CreateDNSRecordCommand) (DNSRecord, error)
}

func (c CreateDNSRecordCommand) Validate() error {
	if strings.TrimSpace(c.Name) == "" || strings.TrimSpace(c.Value) == "" {
		return fmt.Errorf("%w: name and value are required", ErrInvalidRequest)
	}
	if c.TTL < 60 || c.TTL > 86400 {
		return fmt.Errorf("%w: ttl must be between 60 and 86400", ErrInvalidRequest)
	}
	switch strings.ToUpper(c.Type) {
	case "A":
		if ip := net.ParseIP(c.Value); ip == nil || strings.Contains(c.Value, ":") {
			return fmt.Errorf("%w: A record value must be IPv4", ErrInvalidRequest)
		}
	case "AAAA":
		if ip := net.ParseIP(c.Value); ip == nil || !strings.Contains(c.Value, ":") {
			return fmt.Errorf("%w: AAAA record value must be IPv6", ErrInvalidRequest)
		}
	case "CNAME", "NS", "MX", "TXT":
	default:
		return fmt.Errorf("%w: unsupported record type %q", ErrInvalidRequest, c.Type)
	}
	return nil
}
