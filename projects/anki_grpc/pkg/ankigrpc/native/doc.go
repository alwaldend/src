// Package native provides optional native Anki integration through the Rust C ABI.
// RunServer hosts the native CLI synchronously; Rust owns argument validation,
// transports and storage semantics. The executable owner registers its worker
// schemes and supplies their entrypoint and arguments.
// Importing this package registers inprocess and inprocess+cgo. Client owns its connection and
// drains active calls before releasing native state. Cancellation stops waiting
// for startup, commands and shutdown without interrupting Rust. Done signals final
// cleanup, including late initialization. Equal labels never share state.
// Importing this package links the native dependency but initializes no backend.
// Generic client and application imports remain pure Go. Callers own collection
// workflows, executable packaging and process signals; native crashes share the
// hosting process lifetime.
package native
