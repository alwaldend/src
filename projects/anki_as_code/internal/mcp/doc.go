// Package mcp exposes application workflows and declared native RPCs as MCP tools.
// Discovery omits implicit forwarded aliases while retaining declared Backend methods.
// Server receives App at construction and implements the separate MCP hosting contract.
// It owns descriptor-derived JSON schemas, protocol envelopes and tool results.
// App retains configuration, collection lifetime, publication and native execution;
// the SDK owns stdio transport and the CLI supplies process cancellation.
// Explicit handles retain App-owned collection scopes until close or disconnect;
// path-based calls retain per-operation lifetime. Mutations are not transactional.
package mcp
