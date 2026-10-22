// Package rebuild creates collection archives by applying to a private copy.
// Builder owns source/destination checks, staged copying, atomic publication and
// cleanup. Callers supply the apply operation and own collection/configuration
// semantics and exclusive access to source and destination. All file operations
// and path resolution use the injected Afero filesystem.
// This package does not load configs, plan changes, call Anki or own native
// process and collection lifetimes; the supplied operation owns that work.
package rebuild
