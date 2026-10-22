// Package collection manages reserved offline collection scopes and publication.
// Store owns source alias reservations and session-context cancellation.
// Source inspection and collection-directory operations share its injected Afero
// filesystem. Servers must independently access that same storage and paths.
// It delegates address-selected transport creation to the backend library.
// Open starts an ankiclient.StartableClient and prepares one collection.
// It returns ankiclient.Client with Connection and Stop. WithCollection bounds
// that lifetime around a callback. This package owns temporary archive restores and publication.
// Optional engines are injected through a constructor;
// this package stays pure Go regardless of the executable's compiled support.
// Native storage closes before transport shutdown; a private Done-only capability
// observes concrete transport completion. Private completion tracking
// retains reservations while owned native work can still access collection files.
// Enabled package write-back runs during cleanup independently of operation
// errors or cancellation. Database writes persist in place.
// Explicit addresses select network clients and leave the external server running.
// Network servers must share filesystem paths and be reserved for each scope.
// Callers own services, configuration, planning, reconciliation and external
// writer coordination. This package does not discover executables, parse CLI
// input or keep borrowed connections usable after their owner closes.
package collection
