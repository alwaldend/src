// Package config reads and writes schema-backed Anki configuration files.
//
// Load reads the fixed hierarchy through protobuf schemas and checks syntax,
// schema types, and filesystem constraints. Dump encodes the parsed configs.
// Export writes supplied path/config objects directly, without validation,
// normalization, staging, or publication. Callers own collection semantics,
// state-dependent checks, and Anki operations.
// The layout child package owns discovery/file I/O; codec owns serialization
// without a filesystem dependency. This package exposes the consumer API.
package config
