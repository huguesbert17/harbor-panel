// Package auth contains transport-neutral authentication and authorization
// primitives shared by backend modules.
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/harbor-panel/harbor-panel/internal/identity"
)

var (
	ErrUnauthenticated = errors.New("authentication required")
	ErrForbidden       = errors.New("permission denied")
)

type AuthenticationMethod string

const (
	AuthenticationMethodSession AuthenticationMethod = "session"
	AuthenticationMethodToken   AuthenticationMethod = "personal_access_token"
)

// Principal is the authenticated subject made available to application
// modules. It contains authorization data already resolved by the auth layer;
// modules do not inspect cookies, headers, or token formats.
type Principal struct {
	UserID          identity.UserID
	SessionID       identity.SessionID
	Roles           []identity.Role
	Permissions     []identity.Permission
	Method          AuthenticationMethod
	AuthenticatedAt time.Time
}

type Credentials struct {
	Method AuthenticationMethod
	Value  string
}

type Authenticator interface {
	Authenticate(ctx context.Context, credentials Credentials) (Principal, error)
}

type principalContextKey struct{}

// WithPrincipal returns a context carrying an authenticated principal.
// Middleware is responsible for placing only a successfully authenticated
// principal into a request context.
func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

// PrincipalFromContext returns the principal installed by authentication
// middleware. The boolean is false when the request has not been
// authenticated.
func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(Principal)
	return principal, ok
}

func (p Principal) HasPermission(permission identity.Permission) bool {
	for _, granted := range p.Permissions {
		if granted == permission {
			return true
		}
	}
	return false
}

type Resource struct {
	Type string
	ID   string
}

// Authorizer evaluates a permission for a principal and optional resource.
// Resource-level checks remain with the implementation because hosting
// accounts, domains, and servers will define their own ownership rules.
type Authorizer interface {
	Authorize(ctx context.Context, principal Principal, permission identity.Permission, resource Resource) error
}
