## 1. Server-free WebAssembly build

- [x] 1.1 Add the project-local WASM flag and public alias selection; verify the selected Go variant is WebAssembly with neither bundled adapter in its configured dependencies.
- [x] 1.2 Make executable lookup conditional on bundled engines and verify WASI help, config export/dump, malformed input, and unavailable server hosting through the actual module.
- [x] 1.3 Verify mixed native/WASI builds and enabled engine flags preserve sibling target configuration while the selected CLI stays server-free.

## 2. Documentation and validation

- [x] 2.1 Document the project-local flag, runtime use, and networking limits; check the documented command against the produced artifact.
- [x] 2.2 Run existing native CLI workflows and repository quality/lint checks; retain the WASI round-trip output as repeatable evidence.

## Failure scenarios recorded before implementation

- Host ELF output is mistaken for a WebAssembly module.
- A mixed build cross-compiles sibling targets or overrides their requested server engines.
- The Rust descriptor generator is incorrectly analyzed for the Go WASI platform.
- Either Rust server adapter remains in the WASI dependency graph.
- Startup fails while resolving the executable or registering signal handling.
- Config export/dump loses fields or cannot access preopened directories.
- Malformed input writes configuration files before reporting an error.
- Server hosting accidentally starts an engine or reports an ambiguous runtime error.
- Native same-executable worker selection changes as a side effect.
