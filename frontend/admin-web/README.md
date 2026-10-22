# HarborPanel administrator web

This Next.js application is the administrator-facing control plane. Its information architecture is intentionally familiar to hosting operators: dashboard, accounts, packages, domains, DNS, servers, services, backups, security, and activity.

## UI boundaries

- Keep administrator navigation and page-level composition in `app/` and `components/` here.
- Keep shared primitives and API contracts in `frontend/packages/`; do not duplicate them here when they are genuinely shared with the customer application.
- Use original HarborPanel branding, icons, copy, and styling. Familiar workflows are allowed; proprietary assets and source are not.
- Keep data access in typed client and server-proxy layers. Do not place privileged server operations in browser code.

## API integration

The dashboard reads `GET /api/v1/accounts` and `GET /api/v1/servers` through the typed client in `lib/api-client.ts`. The same-origin route handler in `app/api/control-plane/[...path]/route.ts` forwards those reads to the control plane. Configure the server-side control-plane URL with:

```sh
HARBOR_API_BASE_URL=http://localhost:8080
```

The client maps backend `snake_case` response fields into local `camelCase` types and validates the collection envelope before rendering. The dashboard keeps API-reported values authoritative; fields not exposed by the current backend are shown as `Not reported` rather than inferred from placeholder data.

Loading, transport/API errors, empty collections, and unavailable server metrics have explicit UI states. The browser only reads the same-origin admin proxy; it does not receive agent credentials or execute server commands.

## Layout patterns

- `components/admin-shell.tsx` owns the administrator frame, navigation groups, context switcher, and account menu.
- `components/ui.tsx` contains the local `Panel`, `MetricCard`, `StatusBadge`, `ProgressBar`, and link-action patterns used by the first dashboard slice.
- `components/icons.tsx` contains the small original line-icon set used by the shell and dashboard.

## Validation

From `frontend/`:

```sh
npm run build
```
