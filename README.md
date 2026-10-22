# HarborPanel

HarborPanel is an open-source hosting control platform with a familiar hosting-panel mental model and an independent implementation. It is intended to provide separate administrator and customer experiences without using proprietary cPanel assets, source code, or branding.

The project aims to be behaviorally familiar, visually original, API-first, auditable, and safe to operate on production Linux servers.

## Foundation

The first foundation is a Go modular monolith with three executable boundaries:

- `api`: control-plane HTTP API
- `worker`: asynchronous job runner
- `agent`: managed-server agent boundary

The API, worker, and agent are intentionally small at this stage. Domain modules communicate through contracts, while infrastructure providers are introduced through adapters. We are not starting with microservices; the modules remain independently testable inside one repository.

## Current status

This repository contains the initial Go foundation:

- API process with liveness, readiness, and version endpoints
- Worker process boundary for asynchronous operations
- Managed-server agent process boundary
- Health contracts and framework-free HTTP primitives
- Initial hosting-account module contract
- Architecture decision record for the modular Go foundation

The project is not production-ready yet. Authentication, persistence, server enrollment, authorization, jobs, reconciliation, and provider adapters are planned next.

## Run

```sh
cd backend
go run ./cmd/api
go run ./cmd/worker
go run ./cmd/agent
```

The API exposes:

- `GET /healthz` — liveness
- `GET /readyz` — readiness
- `GET /api/v1` — API metadata

## Test

```sh
cd backend
GOCACHE=/tmp/harbor-go-cache go test ./...
GOCACHE=/tmp/harbor-go-cache go vet ./...
```

## Repository layout

```text
backend/              Go control plane, worker, agent, and backend modules
  cmd/                 executable entry points
  internal/            private application modules and platform code
  docs/                backend architecture decisions
frontend/             Next.js applications and shared frontend packages
  admin-web/           WHM-like administration interface
  customer-web/        cPanel-like customer interface
  packages/            shared UI, API client, types, and validation
docs/                 repository-wide product and architecture notes
```

Each module has a local README describing its responsibility and boundary. Read the root [AGENTS.md](AGENTS.md) before making changes.

## Design rules

1. Modules own their domain rules and persistence contracts.
2. Modules communicate through interfaces and typed commands, not internal table access.
3. Privileged server changes go through typed agent operations; there is no public arbitrary-shell endpoint.
4. Long-running mutations become jobs and must be idempotent.
5. The control plane records desired state separately from observed server state.

## Development roadmap

The first major vertical slice will be:

```text
Register server
  -> create hosting account
  -> create domain
  -> configure web server
  -> issue TLS certificate
  -> deploy a test site
  -> inspect job and audit history
  -> remove the site cleanly
```

Planned modules, added incrementally, are:

1. Identity and authentication
2. Servers and agent enrollment
3. Hosting accounts and packages
4. Domains and DNS
5. Websites and runtimes
6. Databases
7. Certificates
8. Jobs and reconciliation
9. Audit and notifications
10. Backups and restores
11. Mail
12. Resellers and delegated administration

## License

HarborPanel is licensed under the **GNU Affero General Public License, version 3 or later** (`AGPL-3.0-or-later`). See [LICENSE](LICENSE).

This license is intended to keep improvements to the hosted control platform available to the community, including modifications used to provide the software as a network service. Contributions must be compatible with this license.

This project is independent and is not affiliated with cPanel, LLC or any other proprietary hosting-panel vendor.
