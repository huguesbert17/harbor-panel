package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harbor-panel/harbor-panel/internal/identity"
)

func TestAuthenticateResolvesBearerToken(t *testing.T) {
	authenticator := &recordingAuthenticator{principal: Principal{
		UserID: "user-1",
		Method: AuthenticationMethodToken,
	}}
	middleware := NewMiddleware(authenticator, nil)
	called := false
	next := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		called = true
		principal, ok := PrincipalFromContext(request.Context())
		if !ok || principal.UserID != "user-1" {
			t.Fatalf("expected authenticated principal, got %#v, present=%t", principal, ok)
		}
		writer.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/accounts", nil)
	request.Header.Set("Authorization", "Bearer token-secret")
	recorder := httptest.NewRecorder()
	middleware.Authenticate(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent || !called {
		t.Fatalf("expected authenticated handler call, status=%d called=%t", recorder.Code, called)
	}
	if authenticator.credentials != (Credentials{Method: AuthenticationMethodToken, Value: "token-secret"}) {
		t.Fatalf("unexpected credentials: %#v", authenticator.credentials)
	}
}

func TestAuthenticateResolvesSessionCookie(t *testing.T) {
	authenticator := &recordingAuthenticator{principal: Principal{
		UserID: "user-1",
		Method: AuthenticationMethodSession,
	}}
	middleware := NewMiddleware(authenticator, nil)
	next := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/accounts", nil)
	request.AddCookie(&http.Cookie{Name: DefaultSessionCookieName, Value: "session-secret"})
	recorder := httptest.NewRecorder()
	middleware.Authenticate(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected authenticated handler call, got %d", recorder.Code)
	}
	if authenticator.credentials != (Credentials{Method: AuthenticationMethodSession, Value: "session-secret"}) {
		t.Fatalf("unexpected credentials: %#v", authenticator.credentials)
	}
}

func TestAuthenticateUsesConfiguredSessionCookie(t *testing.T) {
	authenticator := &recordingAuthenticator{principal: Principal{UserID: "user-1"}}
	middleware := NewMiddleware(authenticator, nil)
	middleware.SessionCookieName = "custom_session"

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: "custom_session", Value: "session-secret"})
	recorder := httptest.NewRecorder()
	middleware.Authenticate(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected custom session cookie to authenticate, got %d", recorder.Code)
	}
}

func TestAuthenticateRejectsMissingAndMalformedCredentials(t *testing.T) {
	tests := []struct {
		name       string
		authorize  string
		withCookie bool
	}{
		{name: "missing"},
		{name: "wrong scheme", authorize: "Basic credentials", withCookie: true},
		{name: "missing bearer value", authorize: "Bearer", withCookie: true},
		{name: "duplicate authorization headers", authorize: "Bearer one", withCookie: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authenticator := &recordingAuthenticator{principal: Principal{UserID: "user-1"}}
			middleware := NewMiddleware(authenticator, nil)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.authorize != "" {
				request.Header.Add("Authorization", test.authorize)
			}
			if test.name == "duplicate authorization headers" {
				request.Header.Add("Authorization", "Bearer two")
			}
			if test.withCookie {
				request.AddCookie(&http.Cookie{Name: DefaultSessionCookieName, Value: "session-secret"})
			}

			recorder := httptest.NewRecorder()
			middleware.Authenticate(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				t.Fatal("handler must not run")
			})).ServeHTTP(recorder, request)

			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", recorder.Code)
			}
			if authenticator.called {
				t.Fatal("authenticator must not receive malformed or missing credentials")
			}
		})
	}
}

func TestAuthenticateDoesNotExposeAuthenticationErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "invalid credentials", err: ErrUnauthenticated, wantStatus: http.StatusUnauthorized},
		{name: "dependency failure", err: errors.New("database password=do-not-expose"), wantStatus: http.StatusServiceUnavailable},
	} {
		t.Run(test.name, func(t *testing.T) {
			authenticator := &recordingAuthenticator{err: test.err}
			middleware := NewMiddleware(authenticator, nil)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Authorization", "Bearer token-secret")
			recorder := httptest.NewRecorder()
			middleware.Authenticate(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				t.Fatal("handler must not run")
			})).ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("expected %d, got %d", test.wantStatus, recorder.Code)
			}
			if body := recorder.Body.String(); strings.Contains(body, "do-not-expose") || strings.Contains(body, "token-secret") {
				t.Fatalf("response exposed authentication detail: %q", body)
			}
		})
	}
}

func TestAuthenticateRejectsEmptyPrincipal(t *testing.T) {
	middleware := NewMiddleware(&recordingAuthenticator{principal: Principal{}}, nil)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer token-secret")
	recorder := httptest.NewRecorder()
	middleware.Authenticate(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		t.Fatal("handler must not run")
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected empty principal to be rejected with 401, got %d", recorder.Code)
	}
}

func TestRequirePermissionAuthorizesResource(t *testing.T) {
	authorizer := &recordingAuthorizer{}
	middleware := NewMiddleware(nil, authorizer)
	called := false
	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		called = true
		writer.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/api/v1/servers/server-1", nil)
	request = request.WithContext(WithPrincipal(request.Context(), Principal{UserID: "user-1"}))
	recorder := httptest.NewRecorder()
	middleware.RequirePermission(identity.PermissionUsersRead, func(*http.Request) Resource {
		return Resource{Type: "servers", ID: "server-1"}
	})(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent || !called {
		t.Fatalf("expected authorized handler call, status=%d called=%t", recorder.Code, called)
	}
	if authorizer.principal.UserID != "user-1" || authorizer.permission != identity.PermissionUsersRead || authorizer.resource != (Resource{Type: "servers", ID: "server-1"}) {
		t.Fatalf("unexpected authorization request: %#v", authorizer)
	}
}

func TestRequirePermissionRejectsUnauthenticatedAndForbiddenRequests(t *testing.T) {
	for _, test := range []struct {
		name       string
		principal  *Principal
		err        error
		wantStatus int
	}{
		{name: "missing principal", wantStatus: http.StatusUnauthorized},
		{name: "forbidden", principal: &Principal{UserID: "user-1"}, err: ErrForbidden, wantStatus: http.StatusForbidden},
		{name: "authorizer unavailable", principal: &Principal{UserID: "user-1"}, err: errors.New("dependency failure"), wantStatus: http.StatusServiceUnavailable},
	} {
		t.Run(test.name, func(t *testing.T) {
			authorizer := &recordingAuthorizer{err: test.err}
			middleware := NewMiddleware(nil, authorizer)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.principal != nil {
				request = request.WithContext(WithPrincipal(request.Context(), *test.principal))
			}
			recorder := httptest.NewRecorder()
			middleware.RequirePermission(identity.PermissionUsersRead, nil)(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				t.Fatal("handler must not run")
			})).ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("expected %d, got %d", test.wantStatus, recorder.Code)
			}
		})
	}
}

func TestPermissionAuthorizerRequiresResolvedPermission(t *testing.T) {
	authorizer := PermissionAuthorizer{}
	principal := Principal{UserID: "user-1", Permissions: []identity.Permission{identity.PermissionUsersRead}}

	if err := authorizer.Authorize(context.Background(), principal, identity.PermissionUsersRead, Resource{}); err != nil {
		t.Fatalf("expected granted permission, got %v", err)
	}
	if err := authorizer.Authorize(context.Background(), principal, identity.PermissionUsersManage, Resource{}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden permission, got %v", err)
	}
	if err := authorizer.Authorize(context.Background(), Principal{}, identity.PermissionUsersRead, Resource{}); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expected unauthenticated principal, got %v", err)
	}
}

func TestProtectAuthenticatesBeforeAuthorizing(t *testing.T) {
	authenticator := &recordingAuthenticator{principal: Principal{UserID: "user-1"}}
	middleware := NewMiddleware(authenticator, PermissionAuthorizer{})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	middleware.Protect(identity.PermissionUsersRead, nil, http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		t.Fatal("handler must not run without the resolved permission")
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", recorder.Code)
	}
}

type recordingAuthenticator struct {
	credentials Credentials
	principal   Principal
	err         error
	called      bool
}

func (authenticator *recordingAuthenticator) Authenticate(_ context.Context, credentials Credentials) (Principal, error) {
	authenticator.called = true
	authenticator.credentials = credentials
	return authenticator.principal, authenticator.err
}

type recordingAuthorizer struct {
	principal  Principal
	permission identity.Permission
	resource   Resource
	err        error
}

func (authorizer *recordingAuthorizer) Authorize(_ context.Context, principal Principal, permission identity.Permission, resource Resource) error {
	authorizer.principal = principal
	authorizer.permission = permission
	authorizer.resource = resource
	return authorizer.err
}
