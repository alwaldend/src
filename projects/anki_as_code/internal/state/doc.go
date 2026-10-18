// Package state manages application path bindings in one open Anki collection.
// Store borrows a gRPC connection and selects one metadata key; an empty key uses
// DefaultKey. Dump reads bindings, Import binds existing native objects, Update
// merges entries and Replace installs a complete object. Resources separates
// deck, note, note-type and preset bindings, including note-type member identities.
// Export constructs bindings for projected config paths without native calls or writes.
// Only the selected metadata key is written. Native content, study data, presets
// and unrelated settings remain untouched. Callers own worker/collection lifetime,
// configuration projection, reconciliation and archive publication.
package state
