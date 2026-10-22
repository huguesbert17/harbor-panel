package identity

import (
	"context"
	"time"
)

// PasswordHasher abstracts the password-hashing algorithm and its work factor.
// Implementations must use a password-specific, adaptive algorithm and must
// compare hashes in constant time where the algorithm requires it.
type PasswordHasher interface {
	Hash(ctx context.Context, password string) (PasswordHash, error)
	Compare(ctx context.Context, password string, encoded PasswordHash) (bool, error)
}

// PasswordPolicy validates a plaintext password before it reaches a hasher.
// It must not persist or log the supplied password.
type PasswordPolicy interface {
	Validate(password string) error
}

type SessionCreateCommand struct {
	UserID    UserID
	ExpiresAt time.Time
	UserAgent string
	IPAddress string
}

type SessionStore interface {
	Create(ctx context.Context, command SessionCreateCommand) (Session, error)
	FindByID(ctx context.Context, id SessionID) (Session, error)
	Revoke(ctx context.Context, id SessionID, revokedAt time.Time) error
	RevokeAllForUser(ctx context.Context, userID UserID, revokedAt time.Time) error
}

type TokenCreateCommand struct {
	UserID    UserID
	Name      string
	Scopes    []Permission
	ExpiresAt *time.Time
}

// TokenDigestor creates a non-reversible lookup value from a presented token
// secret. It must never be implemented as reversible encryption.
type TokenDigestor interface {
	Digest(secret TokenSecret) TokenDigest
}

type TokenSecretGenerator interface {
	Generate() (TokenSecret, error)
}

type TokenStore interface {
	Create(ctx context.Context, command TokenCreateCommand, digest TokenDigest) (Token, error)
	FindByDigest(ctx context.Context, digest TokenDigest) (Token, error)
	Revoke(ctx context.Context, id TokenID, revokedAt time.Time) error
	RevokeAllForUser(ctx context.Context, userID UserID, revokedAt time.Time) error
	MarkUsed(ctx context.Context, id TokenID, usedAt time.Time) error
}
