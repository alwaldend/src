// Package client provides the X API client generated from the pinned schema mirror.
// Bazel supplies generated source at build time. The handwritten media aliases
// preserve raw JSON values for binary/base64 unions unsupported by the pinned
// generator; authentication and HTTP transport are configured by callers.
package client
