# Hosting accounts module

## Responsibility

Owns the lifecycle and policy boundary for customer hosting accounts. An account is the unit that will eventually connect a customer, domains, websites, quotas, packages, and server placement.

## Boundary

The module owns account rules and contracts. It must not directly configure Nginx, PHP-FPM, DNS, mail, or databases. Provisioning is performed through typed jobs and provider contracts.

## Current contract

- `Account`: stable account identity, username, and lifecycle status
- `CreateCommand`: input for account creation
- `Repository`: persistence boundary
- `PersistentRepository`: error-aware collection persistence boundary
- `Service`: application boundary for account use cases

The domain package does not import SQL or a database driver. The local
implementation is `internal/persistence/sqlite.AccountRepository`, which
implements both `Repository` and `PersistentRepository`. `List` is retained
for compatibility with the first API slice; new callers should use `ListAll`
when they need to distinguish an empty result from a storage failure.

## Future work

- account service implementation;
- package and feature-policy assignment;
- quota validation;
- server placement;
- suspend, unsuspend, and terminate workflows;
- audit events and provisioning jobs.
