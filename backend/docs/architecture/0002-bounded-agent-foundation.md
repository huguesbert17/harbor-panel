# ADR 0002: Bounded server-agent foundation

## Status

Accepted

## Decision

The first server-agent service exposes only typed enrollment and heartbeat
contracts. `internal/servers.AgentService` validates these envelopes and
delegates state changes through `AgentStore`.

## Boundaries

- Enrollment proof verification, authentication, credential issuance, and
  replay protection belong to a future security/transport layer.
- Persistence belongs in an adapter implementing `AgentStore`.
- The agent protocol contains no arbitrary command, shell, script, or argument
  execution path.
- Capability reporting describes observed availability; it does not grant
  permission to perform an operation.

## Consequences

The service is testable without a database or network and makes the security
boundary explicit. A later transport must reject unknown/revoked agents and
must preserve auditability and idempotency before adding privileged typed
operations.
