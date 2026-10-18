// Package resource supplies generated Anki clients and fresh resource snapshots.
// Store borrows a prepared connection without owning worker or collection lifetime.
// Snapshot performs native reads and returns caller-owned upstream protobufs with
// no writes or persistent cache. Sequential reads require exclusive offline access
// and are not transactional. Consumers own mutations through the generated clients.
// This package does not implement state bindings, configuration projection,
// reconciliation policy, native startup/shutdown or archive publication.
package resource
