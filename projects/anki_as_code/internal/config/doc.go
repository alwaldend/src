// Package config provides a root-bound Store for Anki configuration files.
//
// Store.Load reads the fixed hierarchy through protobuf schemas and checks syntax,
// schema types, and filesystem constraints. Store.Dump encodes the parsed configs.
// Store.Export writes supplied path/config objects directly, without validation,
// normalization, staging, or publication. Callers own collection semantics,
// state-dependent checks, and Anki operations.
// The layout child package owns discovery/file I/O; codec owns serialization
// without a filesystem dependency. Stores select a root without opening Anki
// or retaining parsed messages; each read returns fresh caller-owned values.
package config
