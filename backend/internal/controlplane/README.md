# Control-plane API module

## Responsibility

Exposes the first versioned HTTP routes over backend module contracts. It owns transport concerns only; account and server rules remain in their domain modules.

## Current development routes

- `GET /api/v1/accounts`
- `GET /api/v1/servers`
- `POST /api/v1/servers/registrations`

The current repositories are in-memory development implementations. The
handler is intentionally transport-agnostic and should be wrapped by
`platform/auth.Middleware` at the API composition boundary:

```go
middleware := auth.NewMiddleware(authenticator, authorizer)
protected := middleware.Authenticate(apiHandler)
```

Use `RequirePermission` inside that composition when a route needs a specific
permission or resource check. Authentication and persistent storage must be
added before production use.
