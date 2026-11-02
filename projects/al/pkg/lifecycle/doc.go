// Package lifecycle starts owned resources and releases them in dependency order.
// Cancellation and startup failure retain cleanup errors and await resource shutdown.
package lifecycle
