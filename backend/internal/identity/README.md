# Identity module

## Responsibility

The identity module owns HarborPanel's framework-independent identity and
credential contracts. It defines users, roles, permissions, password-hashing
boundaries, sessions, and personal access-token metadata.

## Boundary

This module contains domain types and interfaces only. It does not implement
HTTP handlers, cookies, database repositories, password hashing, token
generation, persistence, or login flows. Those concerns must be added behind
these contracts and tested separately.

The module never stores or exposes a plaintext password. A `PasswordHash` is an
encoded verifier, and a `Token` stores only a `TokenDigest`. An `IssuedToken`
contains its secret only at issuance so the API layer can show it once.

## Contracts

- `PasswordHasher` abstracts an adaptive password-specific hash algorithm.
- `PasswordPolicy` validates a plaintext password without persisting it.
- `SessionStore` persists and revokes authenticated sessions.
- `TokenSecretGenerator` creates a secret shown only at token issuance.
- `TokenDigestor` converts presented token secrets into one-way lookup values.
- `TokenStore` persists, finds, revokes, and marks personal access tokens used.

Implementations must remain auditable, avoid logging secrets, and preserve
constant-time verification behavior supplied by the selected password or token
algorithm.

## Authorization vocabulary

Initial roles and identity permissions are stable identifiers, not an
authorization policy. A later authorization module may map roles to
permissions and apply resource-level rules without changing this module's
contracts.

## Future work

- user repository and account lifecycle service;
- secure password-hashing implementation selection;
- session and token issuance services;
- MFA and recovery contracts;
- audit events for authentication and credential changes;
- integration with API authentication middleware.
