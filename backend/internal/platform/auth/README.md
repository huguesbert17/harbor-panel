# Authentication platform package

## Responsibility

This package provides transport-neutral authentication context for the
control-plane modules. It defines the authenticated `Principal`, the
authentication method vocabulary, resource context, and the authorization
contract.

## Boundary

The package does not parse HTTP headers or cookies, validate passwords, store
sessions, issue tokens, or decide hosting-resource ownership. The identity
module owns credential contracts; an API adapter will later translate a valid
session or token into a `Principal`.

Application modules should ask an `Authorizer` about a permission and resource
instead of depending on a particular web framework or credential format. An
API adapter may use `Authenticator` to translate credentials into a principal.

## Security rules

- Do not place plaintext passwords or token secrets in `Principal`.
- Do not make application modules inspect cookies, headers, or bearer-token
  syntax.
- Keep resource ownership checks in the authorization implementation for the
  resource-owning module.
- Return `ErrUnauthenticated` when no valid principal exists and
  `ErrForbidden` when a principal lacks permission.

## Future work

- API authentication middleware;
- role-to-permission policy evaluation;
- request context helpers for principals;
- MFA and recent-authentication requirements;
- audit hooks for authorization decisions.
