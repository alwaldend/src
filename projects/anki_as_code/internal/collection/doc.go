// Package collection coordinates source reservations and archive publication.
// Store.WithCollection delegates native preparation and cleanup to backend managers
// and lends a prepared gRPC connection. Callers construct and reuse their own
// resource/state stores and exporters within the callback. Source aliases remain
// reserved through native cleanup. Successful changed package scopes publish back
// to their source; database writes persist in place. External writers and work
// beyond callback completion remain caller-owned.
package collection
