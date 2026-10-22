# Control-plane API module

## Responsibility

Exposes the first versioned HTTP routes over backend module contracts. It owns transport concerns only; account and server rules remain in their domain modules.

## Current development routes

- `GET /api/v1/accounts`
- `GET /api/v1/servers`
- `POST /api/v1/servers/registrations`

The current repositories are in-memory development implementations. Authentication and persistent storage must be added before production use.
