package identity

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestUserCanAuthenticateOnlyWhenActive(t *testing.T) {
	for _, test := range []struct {
		status UserStatus
		want   bool
	}{
		{status: UserStatusInvited, want: false},
		{status: UserStatusActive, want: true},
		{status: UserStatusSuspended, want: false},
		{status: UserStatusDisabled, want: false},
	} {
		user := User{Status: test.status}
		if got := user.CanAuthenticate(); got != test.want {
			t.Errorf("status %q: got %t, want %t", test.status, got, test.want)
		}
	}
}

func TestSessionIsActive(t *testing.T) {
	now := time.Date(2026, time.October, 22, 12, 0, 0, 0, time.UTC)

	if !(Session{ExpiresAt: now.Add(time.Hour)}).IsActive(now) {
		t.Fatal("expected unrevoked, unexpired session to be active")
	}
	if (Session{ExpiresAt: now}).IsActive(now) {
		t.Fatal("expected session at its expiry instant to be inactive")
	}
	revokedAt := now.Add(-time.Minute)
	if (Session{ExpiresAt: now.Add(time.Hour), RevokedAt: &revokedAt}).IsActive(now) {
		t.Fatal("expected revoked session to be inactive")
	}
}

func TestTokenIsActive(t *testing.T) {
	now := time.Date(2026, time.October, 22, 12, 0, 0, 0, time.UTC)

	if !(Token{}).IsActive(now) {
		t.Fatal("expected token without expiry to be active")
	}
	expiresAt := now.Add(time.Hour)
	if !(Token{ExpiresAt: &expiresAt}).IsActive(now) {
		t.Fatal("expected unrevoked, unexpired token to be active")
	}
	revokedAt := now.Add(-time.Minute)
	if (Token{ExpiresAt: &expiresAt, RevokedAt: &revokedAt}).IsActive(now) {
		t.Fatal("expected revoked token to be inactive")
	}
}

func TestContractsCanBeImplementedWithoutFrameworkDependencies(t *testing.T) {
	var _ PasswordHasher = fakePasswordHasher{}
	var _ PasswordPolicy = fakePasswordPolicy{}
	var _ UserStore = fakeUserStore{}
	var _ SessionStore = fakeSessionStore{}
	var _ TokenDigestor = fakeTokenDigestor{}
	var _ TokenSecretGenerator = fakeTokenSecretGenerator{}
	var _ TokenStore = fakeTokenStore{}
}

type fakeUserStore struct{}

func (fakeUserStore) FindByID(context.Context, UserID) (User, error) {
	return User{}, errors.New("not found")
}

func (fakeUserStore) FindByEmail(context.Context, string) (User, error) {
	return User{}, errors.New("not found")
}

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(context.Context, string) (PasswordHash, error) {
	return PasswordHash{}, nil
}

func (fakePasswordHasher) Compare(context.Context, string, PasswordHash) (bool, error) {
	return false, nil
}

type fakePasswordPolicy struct{}

func (fakePasswordPolicy) Validate(string) error { return nil }

type fakeSessionStore struct{}

func (fakeSessionStore) Create(context.Context, SessionCreateCommand) (Session, error) {
	return Session{}, nil
}

func (fakeSessionStore) FindByID(context.Context, SessionID) (Session, error) {
	return Session{}, errors.New("not found")
}

func (fakeSessionStore) Revoke(context.Context, SessionID, time.Time) error { return nil }

func (fakeSessionStore) RevokeAllForUser(context.Context, UserID, time.Time) error { return nil }

type fakeTokenDigestor struct{}

func (fakeTokenDigestor) Digest(TokenSecret) TokenDigest { return "digest" }

type fakeTokenSecretGenerator struct{}

func (fakeTokenSecretGenerator) Generate() (TokenSecret, error) { return "secret", nil }

type fakeTokenStore struct{}

func (fakeTokenStore) Create(context.Context, TokenCreateCommand, TokenDigest) (Token, error) {
	return Token{}, nil
}

func (fakeTokenStore) FindByDigest(context.Context, TokenDigest) (Token, error) {
	return Token{}, errors.New("not found")
}

func (fakeTokenStore) Revoke(context.Context, TokenID, time.Time) error { return nil }

func (fakeTokenStore) RevokeAllForUser(context.Context, UserID, time.Time) error { return nil }

func (fakeTokenStore) MarkUsed(context.Context, TokenID, time.Time) error { return nil }
