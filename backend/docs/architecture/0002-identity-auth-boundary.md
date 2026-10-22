# ADR 0002: Identity and HTTP authorization boundary

## Status

Accepted

## Decision

Identity contracts live in `internal/identity`. Transport-neutral principals,
authentication credentials, and authorization decisions live in
`internal/platform/auth`. HTTP middleware translates only bearer headers and a
configured session cookie into those contracts.

Routes in `cmd/api` may later be protected with:

```go
middleware := auth.NewMiddleware(authenticator, authorizer)
protected := middleware.Protect(permission, resolver, handler)
mux.Handle("GET /api/v1/resource", protected)
```

Authentication and authorization implementations are supplied by application
wiring. This foundation does not add password hashing, token issuance,
sessions, persistence, cookies beyond credential extraction, or frontend code.

## Security constraints

- Invalid or missing credentials produce `401`; denied permissions produce `403`.
- Unexpected auth dependencies produce `503` without exposing internal errors.
- Empty principals never reach handlers.
- Resource ownership remains with the resource-owning module.
- Principal data contains no password, session secret, or token secret.
- Middleware does not provide arbitrary command execution or bypass audit rules.

## Consequences

The API can adopt authentication route-by-route without coupling domain modules
to HTTP credential formats. A future implementation must provide concrete
credential verification, session/token persistence, role-to-permission policy,
resource ownership checks, MFA, and audit events before production use.
