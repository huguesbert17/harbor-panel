# Agent contract module

## Responsibility

Defines the typed control-plane to managed-server boundary. It owns operation
names, request/result envelopes, and validation for the operations that a
managed-server agent may accept.

This package does not execute commands, authorize callers, persist state, or
implement transport. Those responsibilities belong to future process and
infrastructure layers.

## Safe operation contract

`ValidateRequest` applies an explicit allowlist per operation. It rejects
unknown parameters, missing required values, control characters, unsafe paths,
invalid hostnames, and unsupported operation names. In particular, there is no
generic command, script, shell, or argument parameter.

Every request requires an idempotency key. Implementations of `Client` must
validate requests before dispatch, authorize them in the control plane, record
an audit event, and return observed state in the result where applicable.

## Initial operations

- health check
- create hosting account
- create website
- issue certificate
- remove website

Results use typed operation statuses (`accepted`, `succeeded`, `failed`, or
`rejected`). Adding an operation requires a typed constant, an explicit
parameter schema, validation tests, and documentation of its idempotency and
audit behavior.
