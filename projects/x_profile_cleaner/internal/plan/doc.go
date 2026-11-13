// Package plan validates typed protobuf X-operation plans without executing them.
// It owns structural bounds and operation payload validity shared by the worker
// and backend client, plus the trusted metadata envelope used by files and HTTP.
// Applications own authorization and supported executor sets;
// cleanup-specific action/decision agreement remains in package cleanup.
package plan
