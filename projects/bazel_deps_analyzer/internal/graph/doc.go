// Package graph converts Bazel streamed JSON target records to deterministic
// Mermaid dependency flowcharts or ownership inventories. It owns parsing,
// graph selection, grouping, and escaping; callers own input and output lifetimes
// and filesystem operations.
package graph
