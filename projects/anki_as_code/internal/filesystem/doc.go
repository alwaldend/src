// Package filesystem owns application file operations over an injected Afero filesystem.
// FS supplies shared path checks, file/directory I/O, staging and cleanup, and lends
// its Afero backend to collection owners sharing its namespace.
// Its zero value defaults to OS storage, preserving host home
// and working-directory expansion only for that adapter. Other filesystems retain
// their own path namespace.
// Alias resolution and inspection use only the supplied filesystem's capabilities;
// the package never substitutes host storage for an unsupported operation.
// BasePathFs backing-store link targets are translated to virtual paths, including
// nested base directories; parent traversal and targets outside that base are rejected.
// Rooted paths use a canonical relative spelling, and link targets resolve each
// component before applying parent traversal. Temporary files retain virtual names
// independently of their backing file handles.
// Configuration callers can use the same component traversal while rejecting links,
// including links followed by parent components that would disappear under lexical cleanup.
// Callers own file contents, synchronization, native storage and publication policy.
package filesystem
