# Domains and DNS module

Owns domain identity and DNS zone/record contracts for customer and administrator experiences. It does not configure DNS providers or managed servers directly.

## First bounded slice

- `GET /api/v1/domains?account_id=...`
- `GET /api/v1/dns/zones?account_id=...`
- `GET /api/v1/dns/zones/{zoneID}/records`
- `POST /api/v1/dns/zones/{zoneID}/records`

The local adapter is SQLite-backed. Record writes validate `A`, `AAAA`, `CNAME`, `NS`, `MX`, and `TXT`, TTL bounds, and duplicates.

## Future work

Durable repositories, account authorization, lifecycle jobs, DNS-provider adapters, observed-state reconciliation, audit events, and idempotency keys.
