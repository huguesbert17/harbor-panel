package auth

import (
	"context"
	"testing"

	"github.com/harbor-panel/harbor-panel/internal/identity"
)

func TestPrincipalHasPermission(t *testing.T) {
	principal := Principal{
		UserID:      "user-1",
		Permissions: []identity.Permission{identity.PermissionUsersRead},
	}

	if !principal.HasPermission(identity.PermissionUsersRead) {
		t.Fatal("expected granted permission")
	}
	if principal.HasPermission(identity.PermissionUsersManage) {
		t.Fatal("did not expect ungranted permission")
	}
}

func TestEmptyPrincipalHasNoPermissions(t *testing.T) {
	if (Principal{}).HasPermission(identity.PermissionUsersRead) {
		t.Fatal("empty principal must not grant permissions")
	}
}

func TestAuthenticationContractCanBeImplementedWithoutAFramework(t *testing.T) {
	var _ Authenticator = fakeAuthenticator{}
}

type fakeAuthenticator struct{}

func (fakeAuthenticator) Authenticate(context.Context, Credentials) (Principal, error) {
	return Principal{}, nil
}
