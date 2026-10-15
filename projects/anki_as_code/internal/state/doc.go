// Package state manages application-owned path bindings inside an open Anki collection.
// Stores borrow a gRPC connection and select one metadata key; callers own the
// worker and collection lifetime. Empty keys select the versioned DefaultKey.
// State.Resources keeps separate typed maps for decks, notes, note types, and presets.
// Dump reads the stored protobuf object, Import binds existing native objects in
// batches, Update merges supplied entries, and Replace installs a complete object. Only the namespaced state key
// is written; Anki content, study data, presets, and unrelated settings are untouched.
package state
