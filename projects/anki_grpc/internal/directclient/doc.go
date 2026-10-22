// Package directclient adapts generated unary Anki clients to an owned command
// backend without a network transport. Native and browser adapters supply backend
// creation, command execution and destruction; this package owns protobuf calls,
// native error decoding and lifecycle. Cancellation stops waiting, while Done
// observes completed startup, commands and cleanup after Stop. Collection scopes
// must close their collection before stopping the backend.
package directclient
