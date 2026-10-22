# HTTP server package

## Responsibility

Provides framework-free HTTP server defaults shared by the control-plane API and future internal HTTP processes.

## Current behavior

- bounded read, write, header, and idle timeouts;
- graceful shutdown;
- JSON response helper;
- structured request logging.

## Rules

- Keep request handlers thin; business logic belongs in module services.
- Do not put authentication, authorization, or domain rules in this package.
- Preserve safe timeout defaults when adding endpoints.
- Add tests for response shape and middleware behavior.
