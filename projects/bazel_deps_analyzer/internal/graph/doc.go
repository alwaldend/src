// Package graph converts Bazel streamed JSON target records to deterministic
// Mermaid dependency flowcharts, ownership inventories, or shared package
// aggregates serialized as Mermaid or JSON. It owns parsing, graph selection,
// ownership construction, opt-in explicit deps/srcs attribute filtering, and
// escaping. Tool exclusion belongs to upstream query flags; streamed rule
// inputs do not establish per-attribute transition provenance. Callers own input and output lifetimes
// and filesystem operations.
package graph
