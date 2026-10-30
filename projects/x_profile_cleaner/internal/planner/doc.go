// Package planner builds cleanup plans from language-independent rule descriptions.
// It owns lazy inventory, graph traversal, selection precedence, and explanations.
// Supplied callbacks exchange plain values through opaque identifiers; this package
// has no interpreter, authentication, or mutation dependencies.
package planner
