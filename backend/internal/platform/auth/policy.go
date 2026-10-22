package auth

import (
	"context"

	"github.com/harbor-panel/harbor-panel/internal/identity"
)

// RolePolicy is a deny-by-default authorization policy. It evaluates explicit
// principal permissions first, then the configured role grants. Resource
// ownership and tenant checks remain outside this policy and should be
// implemented by the resource-owning module or a composed Authorizer.
type RolePolicy struct {
	grants map[identity.Role]map[identity.Permission]struct{}
}

// NewRolePolicy creates an immutable policy snapshot from role grants. The
// input is copied so callers cannot mutate authorization behavior after
// construction. Empty roles and permissions are ignored.
func NewRolePolicy(grants map[identity.Role][]identity.Permission) *RolePolicy {
	policy := &RolePolicy{grants: make(map[identity.Role]map[identity.Permission]struct{}, len(grants))}
	for role, permissions := range grants {
		if role == "" {
			continue
		}
		roleGrants := make(map[identity.Permission]struct{}, len(permissions))
		for _, permission := range permissions {
			if permission != "" {
				roleGrants[permission] = struct{}{}
			}
		}
		policy.grants[role] = roleGrants
	}
	return policy
}

// Authorize implements Authorizer. It never grants access for an unknown
// permission, role, or resource. Resource-specific checks can be composed
// around this policy by a module that owns the resource.
func (policy *RolePolicy) Authorize(_ context.Context, principal Principal, permission identity.Permission, _ Resource) error {
	if policy == nil || principal.UserID == "" {
		return ErrUnauthenticated
	}
	if permission == "" {
		return ErrForbidden
	}
	if principal.HasPermission(permission) {
		return nil
	}
	for _, role := range principal.Roles {
		if _, ok := policy.grants[role][permission]; ok {
			return nil
		}
	}
	return ErrForbidden
}

// DefaultRolePolicy contains only identity-management grants. Feature
// modules should extend this policy at composition time rather than making
// this package aware of their domain permissions.
func DefaultRolePolicy() *RolePolicy {
	return NewRolePolicy(map[identity.Role][]identity.Permission{
		identity.RolePlatformAdmin: {
			identity.PermissionUsersRead,
			identity.PermissionUsersManage,
			identity.PermissionSessionsRevoke,
			identity.PermissionTokensManage,
			identity.PermissionRolesManage,
		},
		identity.RoleReadOnly: {
			identity.PermissionUsersRead,
		},
	})
}
