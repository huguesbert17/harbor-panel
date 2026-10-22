package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/harbor-panel/harbor-panel/internal/identity"
)

const DefaultSessionCookieName = "harbor_session"

// ResourceResolver derives the resource being authorized from an HTTP
// request. It is intentionally transport-level glue; resource ownership
// rules remain in the Authorizer implementation.
type ResourceResolver func(*http.Request) Resource

// Middleware adapts HTTP credentials to the transport-neutral authentication
// and authorization contracts. It does not resolve sessions or tokens itself.
// The configured Authenticator owns those lookups and validation rules.
type Middleware struct {
	Authenticator     Authenticator
	Authorizer        Authorizer
	SessionCookieName string
}

func NewMiddleware(authenticator Authenticator, authorizer Authorizer) Middleware {
	return Middleware{
		Authenticator: authenticator,
		Authorizer:    authorizer,
	}
}

// Authenticate resolves a bearer token or session cookie and makes the
// resulting principal available through PrincipalFromContext. Missing,
// malformed, and invalid credentials are rejected with 401. Authenticator
// failures other than ErrUnauthenticated are treated as an unavailable auth
// dependency rather than being allowed through.
func (middleware Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if middleware.Authenticator == nil || next == nil {
			http.Error(writer, "authentication middleware is unavailable", http.StatusInternalServerError)
			return
		}

		credentials, ok := credentialsFromRequest(request, middleware.sessionCookieName())
		if !ok {
			writeUnauthenticated(writer)
			return
		}

		principal, err := middleware.Authenticator.Authenticate(request.Context(), credentials)
		if err != nil {
			if errors.Is(err, ErrUnauthenticated) {
				writeUnauthenticated(writer)
				return
			}
			http.Error(writer, "authentication service unavailable", http.StatusServiceUnavailable)
			return
		}
		if principal.UserID == "" {
			writeUnauthenticated(writer)
			return
		}

		next.ServeHTTP(writer, request.WithContext(WithPrincipal(request.Context(), principal)))
	})
}

// RequirePermission checks an already-authenticated request before invoking
// the wrapped handler. It must be composed inside Authenticate. A missing
// principal is still rejected with 401, never treated as an empty principal.
func (middleware Middleware) RequirePermission(permission identity.Permission, resolve ResourceResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if middleware.Authorizer == nil || next == nil {
				http.Error(writer, "authorization middleware is unavailable", http.StatusInternalServerError)
				return
			}

			principal, ok := PrincipalFromContext(request.Context())
			if !ok || principal.UserID == "" {
				writeUnauthenticated(writer)
				return
			}

			resource := Resource{}
			if resolve != nil {
				resource = resolve(request)
			}
			if err := middleware.Authorizer.Authorize(request.Context(), principal, permission, resource); err != nil {
				switch {
				case errors.Is(err, ErrUnauthenticated):
					writeUnauthenticated(writer)
				case errors.Is(err, ErrForbidden):
					http.Error(writer, "permission denied", http.StatusForbidden)
				default:
					http.Error(writer, "authorization service unavailable", http.StatusServiceUnavailable)
				}
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

// Protect composes authentication and one permission check in the safe order.
// It is intended for route wiring in cmd/api:
//
//	protected := middleware.Protect(permission, resolve, handler)
//	mux.Handle("GET /api/v1/resource", protected)
//
// Authentication always runs before authorization, and neither layer invokes
// the application handler when its dependency is unavailable.
func (middleware Middleware) Protect(permission identity.Permission, resolve ResourceResolver, next http.Handler) http.Handler {
	return middleware.Authenticate(middleware.RequirePermission(permission, resolve)(next))
}

func (middleware Middleware) sessionCookieName() string {
	if middleware.SessionCookieName != "" {
		return middleware.SessionCookieName
	}
	return DefaultSessionCookieName
}

func credentialsFromRequest(request *http.Request, sessionCookieName string) (Credentials, bool) {
	authorizationValues := request.Header.Values("Authorization")
	if len(authorizationValues) > 1 {
		return Credentials{}, false
	}
	if len(authorizationValues) == 1 {
		parts := strings.Fields(authorizationValues[0])
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			return Credentials{}, false
		}
		return Credentials{Method: AuthenticationMethodToken, Value: parts[1]}, true
	}

	cookie, err := request.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return Credentials{}, false
	}
	return Credentials{Method: AuthenticationMethodSession, Value: cookie.Value}, true
}

func writeUnauthenticated(writer http.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", `Bearer realm="harbor-panel"`)
	http.Error(writer, "authentication required", http.StatusUnauthorized)
}
