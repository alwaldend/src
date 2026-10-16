// Package resource provides native Anki clients and fresh resource snapshots.
// Store borrows a prepared connection and owns no backend or collection lifetime.
// It leaves state bindings, configuration projection and reconciliation policy to
// callers. Snapshot reads upstream protobufs without writes or a persistent cache;
// sequential reads require exclusive offline access and are not transactional.
package resource
