# ADR 0001: Modular Go foundation

## Status

Accepted

## Decision

HarborPanel starts as a Go modular monolith with separate API, worker, and agent executable boundaries. Domain modules expose contracts and do not depend on infrastructure providers directly.

## Reasons

- Go produces small, portable binaries for managed servers.
- A modular monolith keeps early development and local testing simple.
- Typed contracts let providers such as Nginx, Apache, PowerDNS, and different database engines evolve independently.
- Separate worker and agent boundaries make privileged and long-running operations explicit.

## Consequences

- We must maintain clear package boundaries and avoid shared database shortcuts.
- The first foundation has no external dependencies beyond the Go toolchain.
- Jobs, persistence, authentication, and provider adapters will be added behind the existing boundaries.
