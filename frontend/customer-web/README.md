# HarborPanel customer web

The customer web app is the account-facing hosting workspace. It keeps the familiar grouping of a traditional hosting panel while using HarborPanel's own brand, visual language, and implementation.

## Current slice

The dashboard currently demonstrates:

- grouped navigation for websites, domains, applications, mail, and support;
- quick-access links for common account tools;
- account summary cards for sites, databases, mailboxes, and certificates;
- accessible resource usage meters;
- recent account activity;
- website and domain status with SSL state;
- responsive navigation and layouts for smaller screens;
- keyboard focus states, a skip link, semantic headings, labelled navigation, and reduced-motion support.

The values are presentation fixtures until the versioned control-plane API is connected. Keep the fixture shape close to the future API contract so the integration remains an explicit boundary.

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
