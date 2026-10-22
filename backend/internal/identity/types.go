// Package identity owns the platform's identity and credential domain types.
// It deliberately contains no persistence, transport, hashing, or framework
// implementation.
package identity

import "time"

type UserID string
type SessionID string
type TokenID string

type Role string
type Permission string

const (
	RolePlatformAdmin Role = "platform_admin"
	RoleReseller      Role = "reseller"
	RoleCustomer      Role = "customer"
	RoleReadOnly      Role = "read_only"
)

const (
	PermissionUsersRead      Permission = "identity.users.read"
	PermissionUsersManage    Permission = "identity.users.manage"
	PermissionSessionsRevoke Permission = "identity.sessions.revoke"
	PermissionTokensManage   Permission = "identity.tokens.manage"
	PermissionRolesManage    Permission = "identity.roles.manage"
)

type UserStatus string

const (
	UserStatusInvited   UserStatus = "invited"
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDisabled  UserStatus = "disabled"
)

type User struct {
	ID          UserID     `json:"id"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	Status      UserStatus `json:"status"`
	Roles       []Role     `json:"roles"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u User) CanAuthenticate() bool {
	return u.Status == UserStatusActive
}

// PasswordHash is an encoded password verifier. Value must contain an
// algorithm-specific encoded hash, never a plaintext password.
type PasswordHash struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"value"`
}

type Session struct {
	ID         SessionID  `json:"id"`
	UserID     UserID     `json:"user_id"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	LastSeenAt time.Time  `json:"last_seen_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	UserAgent  string     `json:"user_agent,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
}

func (s Session) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && now.Before(s.ExpiresAt)
}

// Token is persisted token metadata. The presented secret is intentionally
// absent; only its digest should be stored by a token implementation.
type Token struct {
	ID         TokenID      `json:"id"`
	UserID     UserID       `json:"user_id"`
	Name       string       `json:"name"`
	Digest     TokenDigest  `json:"-"`
	Scopes     []Permission `json:"scopes"`
	CreatedAt  time.Time    `json:"created_at"`
	ExpiresAt  *time.Time   `json:"expires_at,omitempty"`
	LastUsedAt *time.Time   `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time   `json:"revoked_at,omitempty"`
}

func (t Token) IsActive(now time.Time) bool {
	if t.RevokedAt != nil {
		return false
	}
	return t.ExpiresAt == nil || now.Before(*t.ExpiresAt)
}

// TokenSecret is returned only when a token is issued. Callers must not store
// it after presenting it to the user.
type TokenSecret string

// TokenDigest is the one-way representation used for token lookup.
type TokenDigest string

type IssuedToken struct {
	Token  Token       `json:"token"`
	Secret TokenSecret `json:"secret"`
}
