package controlplane

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harbor-panel/harbor-panel/internal/hostingaccounts"
	"github.com/harbor-panel/harbor-panel/internal/servers"
)

func TestRegisterServer(t *testing.T) {
	handler := NewHandler(hostingaccounts.NewMemoryRepository(), servers.NewMemoryRepository())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/servers/registrations", strings.NewReader(`{"name":"web-01","hostname":"web-01.example.com","platform":"debian","architecture":"arm64"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request.WithContext(context.Background()))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestRegisterServerRejectsInvalidPayload(t *testing.T) {
	handler := NewHandler(hostingaccounts.NewMemoryRepository(), servers.NewMemoryRepository())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/servers/registrations", strings.NewReader(`{"name":"web-01","hostname":"web-01.example.com; rm -rf /","platform":"debian","architecture":"arm64"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
}
