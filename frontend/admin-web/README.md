# HarborPanel administrator web

This Next.js application is the administrator-facing control plane. Its information architecture is intentionally familiar to hosting operators: dashboard, accounts, packages, domains, DNS, servers, services, backups, security, and activity.

## UI boundaries

- Keep administrator navigation and page-level composition in `app/` and `components/` here.
- Keep shared primitives and API contracts in `frontend/packages/`; do not duplicate them here when they are genuinely shared with the customer application.
- Use original HarborPanel branding, icons, copy, and styling. Familiar workflows are allowed; proprietary assets and source are not.
- Keep this interface presentation-only until versioned backend contracts are available. Do not place privileged server operations in browser code.

## Layout patterns

- `components/admin-shell.tsx` owns the administrator frame, navigation groups, context switcher, and account menu.
- `components/ui.tsx` contains the local `Panel`, `MetricCard`, `StatusBadge`, `ProgressBar`, and link-action patterns used by the first dashboard slice.
- `components/icons.tsx` contains the small original line-icon set used by the shell and dashboard.

## Validation

From `frontend/`:

```sh
npm run build
```
