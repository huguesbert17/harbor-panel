# HarborPanel agent guide

This repository is an open-source hosting control platform. Preserve the modular architecture and the security boundary between the control plane and managed servers.

## Before changing code

1. Read this file and the `README.md` at the repository root.
2. Read the `README.md` in the module you are changing.
3. Inspect the relevant architecture decision records in `backend/docs/architecture/`.
4. Check the working tree before editing and preserve unrelated user changes.

## Architecture rules

- Keep the backend in Go unless an architecture decision explicitly changes that direction.
- Prefer a modular monolith over new microservices.
- Each module owns its domain rules, contracts, persistence interfaces, handlers, jobs, and tests.
- Modules communicate through typed interfaces and commands; do not reach into another module's internals or tables.
- Infrastructure implementations belong behind provider interfaces in adapters.
- The control plane must never expose arbitrary shell execution.
- Privileged server changes must use typed, auditable, idempotent agent operations.
- Long-running mutations belong in jobs and must support retries and observable failure states.
- Keep desired state separate from observed server state.
- Keep API contracts versioned and backwards-compatible where possible.

## Go conventions

- Use `gofmt` for all Go files.
- Prefer the standard library until a dependency has a clear, documented benefit.
- Return errors with useful context; do not discard operational errors silently.
- Keep packages small and cohesive.
- Add or update tests with behavior changes.
- Run from `backend/`:

```sh
cd backend
GOCACHE=/tmp/harbor-go-cache go test ./...
GOCACHE=/tmp/harbor-go-cache go vet ./...
```

## Repository separation

- `backend/` contains Go control-plane, worker, agent, modules, and backend documentation.
- `frontend/` contains the Next.js administrator and customer applications plus shared frontend packages.
- Keep frontend code dependent on versioned API contracts, never on backend internals.
- Keep backend code independent of a particular frontend application.

## Documentation requirements

- Update the module README when its responsibility, contract, or integration changes.
- Add an architecture decision record for cross-module or security-boundary decisions.
- Document new environment variables in `.env.example` and the root README.
- Do not use cPanel assets, proprietary text, screenshots, branding, or source code.

## Delivery checklist

- `gofmt` has been run.
- Tests and vet pass.
- Documentation is updated.
- No secrets, credentials, generated files, or `.DS_Store` files are committed.
- The change preserves authorization, auditability, and tenant isolation.
