# Authentication platform package

## Responsibility

This package provides authentication context and framework-free HTTP
middleware for the control-plane modules. It defines the authenticated
`Principal`, the authentication method vocabulary, resource context, and the
authorization contract.

## Boundary

The package does not validate passwords, store sessions, issue tokens, or decide
hosting-resource ownership. The identity module owns credential contracts and
the configured `Authenticator` resolves presented credentials into a
`Principal`. The HTTP adapter only translates an `Authorization: Bearer ...`
header or the configured session cookie into those credentials.

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
- Never fall back from malformed bearer credentials to a session cookie.
- Never pass an empty principal to an application handler.
- Do not include credentials or authenticator errors in HTTP responses.

## HTTP middleware

Create the middleware with an `Authenticator` and an `Authorizer`:

```go
middleware := auth.NewMiddleware(authenticator, authorizer)
protected := middleware.RequirePermission(
    identity.PermissionUsersRead,
    func(request *http.Request) auth.Resource {
        return auth.Resource{Type: "users", ID: request.PathValue("userID")}
    },
)(handler)
protected = middleware.Authenticate(protected)
```

`Authenticate` accepts a bearer token or the `harbor_session` cookie by
default. Set `SessionCookieName` when the deployment uses another cookie name.
Successful requests receive their principal through `PrincipalFromContext`.
Missing or invalid credentials return `401`; denied permissions return `403`;
unexpected authentication or authorization dependency failures return `503`.
Middleware misconfiguration returns `500` and never invokes the wrapped
handler.

## Future work

- role-to-permission policy evaluation;
- MFA and recent-authentication requirements;
- audit hooks for authorization decisions.
