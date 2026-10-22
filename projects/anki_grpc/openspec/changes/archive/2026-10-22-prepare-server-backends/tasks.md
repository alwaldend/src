## 1. Implementation and validation

- [x] 1.1 Extract the existing cgo adapter, platform-independent Rust CLI library and native command targets; verify native CLI and retain public executable names.
- [x] 1.2 Include the existing SQLite native/WASI fixture and built-in Rust platform selection without server code; verify cross-runtime reopen.
- [x] 1.3 Verify native transports, linked consumers, SQLite artifacts and OpenSpec.

## 2. Built-in platform selection

- [x] 2.1 Remove the custom WASI transition package, build native/WASI SQLite probes with rules_rust platform attributes, and verify cross-reopen artifacts; delivery receipts own required quality/lint gates.
