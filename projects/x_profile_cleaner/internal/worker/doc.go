// Package worker implements transient asynchronous archive-script jobs over private gRPC.
// Configured providers fetch exports directly; scripts collect typed operations.
// The composition root supplies trusted runtime helpers. Result semantics and X
// execution belong to applications; this package has no mutation or credential API.
// The caller owns supervision and operating-system resource isolation.
package worker
