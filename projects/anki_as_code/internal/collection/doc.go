// Package collection manages reserved offline collection scopes and publication.
// Store owns source alias reservations and composes explicitly configured worker
// and collection clients, lending a prepared gRPC connection to each callback.
// Native storage closes before worker shutdown; reservations outlive cleanup.
// Enabled package write-back runs during cleanup independently of operation
// errors or cancellation. Database writes persist in place.
// Callers own services, configuration, planning, reconciliation and external
// writer coordination. This package does not discover executables, parse CLI
// input or keep borrowed connections usable beyond callback completion.
package collection
