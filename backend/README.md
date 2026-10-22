# HarborPanel backend

The backend is a Go modular monolith with three executable boundaries:

- `cmd/api`: control-plane HTTP API
- `cmd/worker`: asynchronous job runner
- `cmd/agent`: managed-server agent

The API uses SQLite for local persistent development storage. Set `HARBOR_DB_PATH` to choose the database file. A PostgreSQL adapter is planned for multi-node production deployments.

Domain modules live under `internal/`. Infrastructure providers must be accessed through typed contracts and adapters. Run backend checks from this directory:

```sh
GOCACHE=/tmp/harbor-go-cache go test ./...
GOCACHE=/tmp/harbor-go-cache go vet ./...
```
