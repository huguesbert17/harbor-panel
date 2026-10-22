# HarborPanel frontend

The frontend will contain separate Next.js applications for the administrator and customer experiences:

```text
frontend/
├── admin-web/       # WHM-like administration interface
├── customer-web/    # cPanel-like customer interface
└── packages/        # shared UI, API client, types, and validation
```

The applications share components and API contracts, but do not share page-level navigation or permissions. The administrator interface is operational and server-focused; the customer interface is account-focused and task-oriented.

## Rules

- Use React, Next.js, and TypeScript.
- Use the backend API rather than importing Go internals or database models.
- Keep frontend models camelCase and translate at the API boundary.
- Enforce authorization in the backend; frontend permissions only improve presentation.
- Keep shared components in `frontend/packages/`.
- Add browser tests for completed vertical slices.
