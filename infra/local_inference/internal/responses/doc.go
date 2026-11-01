// Package responses adapts namespaced Responses function tools to the flat
// function-only interface supported by the pinned llama.cpp server. It preserves
// tool identity on streamed replies and replayed calls; it never executes tools.
package responses
