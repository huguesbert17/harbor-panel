# Worker process

The worker process owns asynchronous and long-running work such as provisioning, certificate issuance, backups, restores, and reconciliation.

Jobs must be idempotent, retryable, observable, and auditable.
