// Package policy adapts Lua configuration and callbacks to the Go cleanup planner.
// It owns the restricted interpreter, callback registry, and fresh Lua projections.
// Inventory, graph traversal, selection precedence, and plans belong to planner;
// credentials and mutation remain outside both packages.
package policy
