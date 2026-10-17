// Package export writes a prepared native collection as editable configs.
// Exporter borrows resource and state stores, reads native snapshots and existing
// bindings for stable filenames, projects native values and paths, and delegates
// binding generation to state.Export and config file output to config.Store.
// It overwrites declared configs in place and leaves other entries alone,
// including state files other than the explicit state destination. Callers own
// cleanup, synchronization, state operations and native lifetime; filesystem
// errors may leave partial config output. Export optionally publishes snapshot
// bindings as TOML after writing configs; it never installs state.
package export
