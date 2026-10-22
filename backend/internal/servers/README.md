# Servers module

## Responsibility

Defines the control-plane contracts for registering managed servers,
establishing agent enrollment, reporting server status, and discovering typed
capabilities.

This module owns lifecycle vocabulary and contract validation only. It does
not authenticate users, issue credentials, persist server records, open
network connections, or execute privileged operations.

## Lifecycle

The intended lifecycle is:

```text
pending -> enrolling -> active -> degraded/offline
                             -> revoked
```

The enrollment request carries an opaque public key and proof. Verification,
one-time enrollment handling, credential issuance, and replay protection must
be implemented by a future enrollment transport and identity/security layer.
These contracts intentionally do not define user identity or authentication.

## Capabilities and status

Agents report capabilities as an allowlisted name, optional provider version,
and status. A capability is usable only when its status is `available`.
`ServerStatus.Supports` is a convenience check for that rule; it does not
replace authorization or provider health checks.

Status separates registration lifecycle from connection state. This lets the
control plane distinguish a revoked server from a temporarily disconnected
agent without inferring one from the other.

## Extension rules

New capabilities and lifecycle states require contract tests and documentation.
Do not place provider-specific implementation details, database access, or
arbitrary command data in this package.
