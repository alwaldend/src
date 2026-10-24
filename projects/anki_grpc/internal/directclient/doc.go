// Package directclient adapts generated unary Anki clients to an owned command
// backend without a network transport. Platform adapters supply backend
// creation, command execution and destruction; this package owns protobuf calls,
// native error decoding and lifecycle. Startup runs synchronously and checks its
// context before and after initialization. Command and shutdown cancellation stop
// waiting, while Done observes completed startup, commands and cleanup after Stop.
// Collection scopes must close their collection before stopping the backend.
package directclient
