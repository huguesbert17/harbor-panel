package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/harbor-panel/harbor-panel/internal/identity"
)

func TestRolePolicyUsesExplicitPermissionsAndRoleGrants(t *testing.T) {
	policy := NewRolePolicy(map[identity.Role][]identity.Permission{
		identity.RoleCustomer: {identity.PermissionUsersRead},
	})

	principal := Principal{UserID: "user-1", Roles: []identity.Role{identity.RoleCustomer}}
	if err := policy.Authorize(context.Background(), principal, identity.PermissionUsersRead, Resource{Type: "users", ID: "user-1"}); err != nil {
		t.Fatalf("role grant should authorize: %v", err)
	}

	principal.Permissions = []identity.Permission{identity.PermissionTokensManage}
	if err := policy.Authorize(context.Background(), principal, identity.PermissionTokensManage, Resource{}); err != nil {
		t.Fatalf("explicit permission should authorize: %v", err)
	}
}

func TestRolePolicyDeniesUnknownAndUnauthenticatedRequests(t *testing.T) {
	policy := DefaultRolePolicy()
	if err := policy.Authorize(context.Background(), Principal{}, identity.PermissionUsersRead, Resource{}); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expected unauthenticated error, got %v", err)
	}

	principal := Principal{UserID: "user-1", Roles: []identity.Role{identity.RoleCustomer}}
	if err := policy.Authorize(context.Background(), principal, identity.PermissionRolesManage, Resource{}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
	if err := policy.Authorize(context.Background(), principal, "", Resource{}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("empty permission must be forbidden, got %v", err)
	}
}

func TestNewRolePolicyCopiesGrants(t *testing.T) {
	permissions := []identity.Permission{identity.PermissionUsersRead}
	grants := map[identity.Role][]identity.Permission{identity.RoleCustomer: permissions}
	policy := NewRolePolicy(grants)
	permissions[0] = identity.PermissionRolesManage

	principal := Principal{UserID: "user-1", Roles: []identity.Role{identity.RoleCustomer}}
	if err := policy.Authorize(context.Background(), principal, identity.PermissionUsersRead, Resource{}); err != nil {
		t.Fatalf("policy changed after caller mutated grants: %v", err)
	}
}
