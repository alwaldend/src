// Package ankiservercgo hosts Anki's native CLI through one synchronous C ABI call.
// Run copies arguments, invokes linked Rust and reports nonzero exit statuses.
// Rust owns subcommand parsing, transports, protocol streams, listeners and its
// runtime; Anki owns storage semantics. Serving owns the caller's process until
// return or termination. Offline callers host it through their worker process owner.
// This package does not compose application workflows, spawn workers, change
// process signal handling or implement collection/configuration policy. Callers
// supply the full native argument vector and any required process adapters.
package ankiservercgo
