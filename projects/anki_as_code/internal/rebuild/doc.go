// Package rebuild creates collection archives by applying changes to a private copy.
// Builder owns input checks, staged copying, atomic publication and cleanup.
// Callers supply the apply operation and own collection/configuration semantics,
// home-path expansion and exclusive access to source and destination paths.
package rebuild
