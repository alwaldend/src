// Package resources implements Vast.ai Terraform resource schemas and lifecycles.
// It owns planning, state, imports, and cleanup; the shared SDK owns API transport.
// Provider registration consumes factories while concrete implementations stay private.
package resources
