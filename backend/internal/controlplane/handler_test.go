package controlplane

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harbor-panel/harbor-panel/internal/domains"
	"github.com/harbor-panel/harbor-panel/internal/hostingaccounts"
	"github.com/harbor-panel/harbor-panel/internal/servers"
)

func TestRegisterServer(t *testing.T) {
	handler := NewHandler(hostingaccounts.NewMemoryRepository(), servers.NewMemoryRepository(), domains.NewService(newTestDomainRepository()))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/servers/registrations", strings.NewReader(`{"name":"web-01","hostname":"web-01.example.com","platform":"debian","architecture":"arm64"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request.WithContext(context.Background()))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestRegisterServerRejectsInvalidPayload(t *testing.T) {
	handler := NewHandler(hostingaccounts.NewMemoryRepository(), servers.NewMemoryRepository(), domains.NewService(newTestDomainRepository()))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/servers/registrations", strings.NewReader(`{"name":"web-01","hostname":"web-01.example.com; rm -rf /","platform":"debian","architecture":"arm64"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
}

type testDomainRepository struct {
	domains map[string]domains.Domain
	zones   map[string]domains.DNSZone
}

func newTestDomainRepository() *testDomainRepository {
	return &testDomainRepository{domains: map[string]domains.Domain{}, zones: map[string]domains.DNSZone{}}
}
func (r *testDomainRepository) CreateDomain(_ context.Context, item domains.Domain) error {
	r.domains[item.ID] = item
	return nil
}
func (r *testDomainRepository) FindDomainByID(_ context.Context, id string) (domains.Domain, error) {
	item, ok := r.domains[id]
	if !ok {
		return domains.Domain{}, domains.ErrNotFound
	}
	return item, nil
}
func (r *testDomainRepository) ListDomains(_ context.Context, _ string) ([]domains.Domain, error) {
	items := make([]domains.Domain, 0, len(r.domains))
	for _, item := range r.domains {
		items = append(items, item)
	}
	return items, nil
}
func (r *testDomainRepository) CreateZone(_ context.Context, item domains.DNSZone) error {
	r.zones[item.ID] = item
	return nil
}
func (r *testDomainRepository) FindZoneByID(_ context.Context, id string) (domains.DNSZone, error) {
	item, ok := r.zones[id]
	if !ok {
		return domains.DNSZone{}, domains.ErrNotFound
	}
	return item, nil
}
func (r *testDomainRepository) ListZones(_ context.Context, _ string) ([]domains.DNSZone, error) {
	items := make([]domains.DNSZone, 0, len(r.zones))
	for _, item := range r.zones {
		items = append(items, item)
	}
	return items, nil
}
func (r *testDomainRepository) ListRecords(_ context.Context, _ string) ([]domains.DNSRecord, error) {
	return nil, nil
}
func (r *testDomainRepository) CreateRecord(_ context.Context, _ string, _ domains.CreateDNSRecordCommand) (domains.DNSRecord, error) {
	return domains.DNSRecord{}, nil
}
