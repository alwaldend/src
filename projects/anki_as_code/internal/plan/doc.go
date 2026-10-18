// Package plan builds executable collection plans from supplied configuration,
// native snapshots and stored bindings. It owns reference resolution, resource
// comparison, ownership checks and operation preparation using private indexes.
// The protobuf plan is the executor boundary; inputs remain caller-owned.
// This package performs no filesystem or network I/O, native calls, resource
// mutation or state persistence. Callers acquire inputs and ensure freshness;
// reconcile executes plans, and Anki owns native validation, normalization and
// generated-card migrations during execution.
package plan
