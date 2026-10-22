# Hosting accounts module

## Responsibility

Owns the lifecycle and policy boundary for customer hosting accounts. An account is the unit that will eventually connect a customer, domains, websites, quotas, packages, and server placement.

## Boundary

The module owns account rules and contracts. It must not directly configure Nginx, PHP-FPM, DNS, mail, or databases. Provisioning is performed through typed jobs and provider contracts.

## Current contract

- `Account`: stable account identity, username, and lifecycle status
- `CreateCommand`: input for account creation
- `Repository`: persistence boundary
- `Service`: application boundary for account use cases

## Future work

- account service implementation;
- package and feature-policy assignment;
- quota validation;
- server placement;
- suspend, unsuspend, and terminate workflows;
- audit events and provisioning jobs.
