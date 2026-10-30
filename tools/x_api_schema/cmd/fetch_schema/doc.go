// Package main downloads an X OpenAPI schema snapshot using injected credentials.
// It performs one signed GET, validates the bounded response, and atomically
// replaces the destination only after a complete successful download. Downloads
// default to ignored workspace scratch for later mirror publication.
package main
