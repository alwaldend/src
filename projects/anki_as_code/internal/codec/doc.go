// Package codec converts protobuf configurations and state to and from TOML and JSON.
// It owns schema-backed serialization, field values and optional presence, and
// rejects TOML datetime coercion. It performs no file I/O, layout discovery,
// version-policy checks, native validation or Anki operations. Config, state and
// application packages own resource semantics and how encoded values are used.
package codec
