// Package health contains the small contracts used by process and orchestration
// health endpoints. It deliberately has no dependency on a web framework.
package health

import "context"

// Checker reports whether a component can currently serve traffic.
type Checker interface {
	Check(ctx context.Context) error
}

// Status is the stable response shape shared by liveness and readiness checks.
type Status struct {
	Status string `json:"status"`
}

const (
	StatusOK       = "ok"
	StatusReady    = "ready"
	StatusNotReady = "not_ready"
)
