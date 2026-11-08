// Package main provides Vault resource types to AL over its generic provider RPC.
// Capability startup is lazy and performs no secret reads. Fetches and shutdown
// retain the invocation's authentication, cancellation and cleanup boundaries.
package main
