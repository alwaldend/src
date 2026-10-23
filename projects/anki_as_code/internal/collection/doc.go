// Package collection manages reserved offline collection scopes and publication.
// Store owns source alias reservations and session-context cancellation.
// Source inspection uses its injected Afero filesystem. Backend libraries and
// servers retain their own storage contracts and must be able to access those paths.
// It delegates address-selected client creation
// to the backend library. Open starts the factory-created ankiclient.Client and
// opens its collection; that same client handles cleanup and publication. WithCollection
// uses Connection and Stop around a callback.
// Optional engines are injected through a constructor;
// this package stays pure Go regardless of the executable's compiled support.
// Native storage closes before transport shutdown; reservations remain until Client.Done.
// Reservation records observe completion without wrapping the client.
// Enabled package write-back runs during cleanup independently of operation
// errors or cancellation. Database writes persist in place.
// Explicit addresses select network clients and leave the external server running.
// Network servers must share filesystem paths and be reserved for each scope.
// Callers own services, configuration, planning, reconciliation and external
// writer coordination. This package does not discover executables, parse CLI
// input or keep borrowed connections usable after their owner closes.
package collection
