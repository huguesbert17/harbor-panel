# HarborPanel customer web

The customer web app is the account-facing hosting workspace. It keeps the familiar grouping of a traditional hosting panel while using HarborPanel's own brand, visual language, and implementation.

## Current slice

The dashboard currently demonstrates:

- grouped navigation for websites, domains, applications, mail, and support;
- quick-access links for common account tools;
- account summary cards backed by the control-plane account and server collections;
- accessible resource usage panel with API loading, error, and empty states;
- recent account activity;
- explicit loading, error, and empty states for API-backed and not-yet-connected areas;
- responsive navigation and layouts for smaller screens;
- keyboard focus states, a skip link, semantic headings, labelled navigation, and reduced-motion support.

The dashboard uses the versioned control-plane API endpoints `GET /api/v1/accounts` and `GET /api/v1/servers`. Domain, usage, activity, and service-specific data remain empty states until their dedicated customer APIs exist; the UI does not present fixture values as live data.

## Configuration

Copy `.env.example` to `.env.local` and set the control-plane origin:

```sh
NEXT_PUBLIC_HARBOR_API_BASE_URL=http://localhost:8080
```

Do not include `/api/v1` in the value. The API client validates the response envelopes and the snake_case backend fields locally before exposing camelCase data to the UI.

## Boundaries

- Keep customer-specific presentation and route composition in this folder.
- Use `@harbor-panel/ui` and versioned API packages without modifying them from a customer-web change.
- Do not put privileged server operations in browser code.
- Use links and typed API actions that will map to customer permissions when the control plane is connected.
- Do not use cPanel assets, screenshots, proprietary copy, or branding.

## Development

From `frontend/`:

```sh
npm run dev:customer
npm run build
```

The app is intended to remain a focused customer experience; administrator workflows belong in `frontend/admin-web`.
