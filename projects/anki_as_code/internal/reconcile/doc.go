// Package reconcile executes prepared protobuf collection plans through Anki.
// It owns native mutation order, returned identities and state persistence, and
// borrows caller-owned resource/state stores. It performs no config resolution,
// collection discovery or planning. Anki owns validation, structural migrations
// and card generation; callers own input freshness and native lifetime.
package reconcile
