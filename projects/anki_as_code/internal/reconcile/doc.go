// Package reconcile executes prepared protobuf collection plans through Anki.
// It owns mutation order, returned identities and state persistence, borrowing
// resource/state stores for the same prepared collection. It checks expected
// bindings before writes; later failures do not roll back committed database work.
// Callers own input freshness, configuration loading, planning, native lifetime
// and archive publication. This package does not discover collections or rebuild
// plans. Anki owns native validation, structural migrations and card generation.
package reconcile
