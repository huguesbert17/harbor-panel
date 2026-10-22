package domains

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type service struct{ repository PersistentRepository }

func NewService(repository PersistentRepository) Service { return &service{repository: repository} }
func (s *service) CreateDomain(ctx context.Context, c CreateDomainCommand) (Domain, error) {
	if strings.TrimSpace(c.AccountID) == "" || !validDomain(c.Name) {
		return Domain{}, fmt.Errorf("%w: account_id and valid name are required", ErrInvalidRequest)
	}
	id, err := newID("dom")
	if err != nil {
		return Domain{}, err
	}
	now := time.Now().UTC()
	item := Domain{ID: id, AccountID: c.AccountID, Name: strings.ToLower(c.Name), Kind: KindPrimary, Status: StatusPending, CreatedAt: now, UpdatedAt: now}
	if err = s.repository.CreateDomain(ctx, item); err != nil {
		return Domain{}, err
	}
	return item, nil
}
func (s *service) ListDomains(ctx context.Context, accountID string) ([]Domain, error) {
	return s.repository.ListDomains(ctx, accountID)
}
func (s *service) CreateZone(ctx context.Context, c CreateZoneCommand) (DNSZone, error) {
	if strings.TrimSpace(c.AccountID) == "" || strings.TrimSpace(c.DomainID) == "" {
		return DNSZone{}, fmt.Errorf("%w: account_id and domain_id are required", ErrInvalidRequest)
	}
	domain, err := s.repository.FindDomainByID(ctx, c.DomainID)
	if err != nil {
		return DNSZone{}, err
	}
	if domain.AccountID != c.AccountID {
		return DNSZone{}, ErrNotFound
	}
	id, err := newID("zone")
	if err != nil {
		return DNSZone{}, err
	}
	now := time.Now().UTC()
	item := DNSZone{ID: id, AccountID: c.AccountID, DomainID: domain.ID, Name: domain.Name, Status: StatusPending, CreatedAt: now, UpdatedAt: now}
	if err = s.repository.CreateZone(ctx, item); err != nil {
		return DNSZone{}, err
	}
	return item, nil
}
func (s *service) ListZones(ctx context.Context, a string) ([]DNSZone, error) {
	return s.repository.ListZones(ctx, a)
}
func (s *service) ListRecords(ctx context.Context, z string) ([]DNSRecord, error) {
	return s.repository.ListRecords(ctx, z)
}
func (s *service) CreateRecord(ctx context.Context, z string, c CreateDNSRecordCommand) (DNSRecord, error) {
	return s.repository.CreateRecord(ctx, z, c)
}
func validDomain(v string) bool {
	if v == "" || len(v) > 253 || strings.TrimSpace(v) != v || strings.HasSuffix(v, ".") || netIP(v) {
		return false
	}
	for _, l := range strings.Split(strings.ToLower(v), ".") {
		if l == "" || len(l) > 63 || l[0] == '-' || l[len(l)-1] == '-' {
			return false
		}
		for _, c := range l {
			if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
				return false
			}
		}
	}
	return true
}
func netIP(v string) bool { return strings.Contains(v, ":") || strings.Count(v, ".") == 3 }
func newID(p string) (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return p + "_" + hex.EncodeToString(b), nil
}
