// Package export projects prepared native collections into editable config files.
// Exporter borrows resource/state stores, reads snapshots and bindings for stable
// filenames, derives editable values and paths, and delegates binding generation
// to state.Export and config file output to config.Store. Optional snapshot
// bindings are published as TOML after config output; stored state is never installed.
// Declared configs are overwritten and other entries are retained. Filesystem
// failures may leave partial configs; state output uses atomic replacement.
// Every output inspection and file operation uses the injected Afero filesystem.
// Callers own collection/source exclusion checks, native lifetime, synchronization
// and workflow policy. This package does not start workers, open collections,
// reconcile resources or own the supplied resource/state stores.
package export
