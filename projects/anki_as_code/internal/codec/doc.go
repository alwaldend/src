// Package codec converts protobuf configurations and state to and from TOML and JSON.
//
// It preserves field values and optional presence, checks syntax/schema types,
// and rejects TOML datetime coercion. It performs no file I/O, layout discovery,
// version-policy checks, native validation, or Anki operations.
package codec
