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

// AuthorizerFunc adapts a function to Authorizer. It is useful for composing
// a policy implementation during application wiring and in focused tests.
type AuthorizerFunc func(context.Context, Principal, identity.Permission, Resource) error

func (authorizer AuthorizerFunc) Authorize(ctx context.Context, principal Principal, permission identity.Permission, resource Resource) error {
	return authorizer(ctx, principal, permission, resource)
}

// PermissionAuthorizer is the baseline policy for permissions already resolved
// into a Principal. Resource ownership checks should be layered on top by the
// owning module; this type intentionally does not grant access by role name.
type PermissionAuthorizer struct{}

func (PermissionAuthorizer) Authorize(_ context.Context, principal Principal, permission identity.Permission, _ Resource) error {
	if principal.UserID == "" {
		return ErrUnauthenticated
	}
	if !principal.HasPermission(permission) {
		return ErrForbidden
	}
	return nil
}
