# Platform packages

## Responsibility

Contains narrow, reusable application infrastructure that is not owned by a business domain. Platform code supports modules but must not contain hosting-account, domain, billing, or provider-specific business rules.

## Boundary

Platform packages may depend on the Go standard library and carefully selected cross-cutting infrastructure. They must expose stable primitives and avoid becoming a dumping ground for unrelated helpers.

## Current packages

- `httpserver`: HTTP server lifecycle, JSON responses, and request logging

## Future work

- structured error responses;
- request IDs and tracing;
- authentication middleware;
- job and event abstractions;
- transaction and clock interfaces where tests require them.
