// Package gazelle generates package transitive visibility from validated TOML
// path mappings. It preserves unrelated BUILD attributes and comments, positions
// existing package declarations before rules, and never searches parent workspaces.
// Package-group definitions remain owned by their declaring BUILD files.
package gazelle
