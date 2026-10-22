# Health module

## Responsibility

Provides small contracts and response types for process liveness and service readiness.

## Boundary

This package must remain independent of HTTP frameworks, databases, and deployment providers. Components implement `health.Checker` when they need to participate in readiness checks.

## Current contract

- `Checker`: context-aware component health check
- `Status`: stable JSON response shape
- status constants for healthy, ready, and not-ready responses

## Future work

- aggregate readiness checks for database, queue, and agent connectivity;
- distinguish liveness failures from dependency readiness failures;
- add diagnostic details without exposing secrets.
