// Package export writes a prepared native collection as editable configs.
// Exporter borrows resource and state stores, reads native snapshots and existing
// bindings for stable filenames, projects native values and paths,
// and delegates file output to config.Store. It overwrites declared configs
// in place and leaves other entries alone, including state files. Callers own
// cleanup, synchronization, state operations and native lifetime; filesystem
// errors may leave partial output. Export never generates or installs state.
package export
