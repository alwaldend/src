// Package layout reads and writes the fixed Anki configuration hierarchy.
// It owns filesystem discovery, regular-file checks, supported versions,
// path-derived deck titles, note deck defaults and config file I/O. Codec owns
// syntax/schema conversion. Loading leaves source files unchanged; exporting
// writes supplied objects directly and may leave partial output on failure.
// Callers own collection semantics, native identities, synchronization, cleanup,
// staging and publication. This package does not open Anki or manage native work.
package layout
