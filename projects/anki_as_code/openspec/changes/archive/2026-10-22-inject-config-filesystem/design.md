## Context

App composes config, collection, export and rebuild owners. Each currently performs direct OS operations and several duplicate path expansion. Afero v1.11.0 is already pinned transitively. Server and backend libraries consume paths separately and are explicitly outside the filesystem migration.

## Goals / Non-Goals

Supply one filesystem to every application-owned I/O path and prove configuration operations need no host storage. Keep serialization, collection lifetime and native behavior intact. Browser bindings and server/backend filesystem adapters are separate work.

## Decisions

- Add `Options.FS afero.Fs`, pass it into each filesystem owner, and preserve nil as the OS default. Keep zero-value App behavior. Explicit construction avoids process-global overrides and permits independent applications.
- Centralize path handling in `internal/filesystem`, a small adapter around the supplied Afero filesystem. Only the OS adapter expands host home/current-directory notation; custom storage retains its own relative and absolute names. This removes duplicate implementations and prevents hidden host path lookups.
- Use Afero directory/file/temp operations and optional Lstater/LinkReader capabilities for path inspection and symlink resolution. Never substitute OS calls for missing custom capabilities. Preserve host hardlink checks through the existing pure metadata comparison.
- Compare resolved source/output paths directly before relying on OS file identity. Custom filesystem metadata need not support `os.SameFile`, so metadata comparison alone cannot reject an identical path.
- For `BasePathFs`, translate absolute backing-store link targets using its full base directory and Afero's public path utility. Resolve relative link contents from the virtual link directory. Use one canonical relative spelling for equivalent virtual paths, handle directly nested bases and reject targets outside the base. Process symlink target components before applying parent traversal. Applying the base twice is invalid; resolving through host OS calls would bypass injection, and skipping alias resolution would weaken the existing conflict guards.
- Keep existing backend path contracts visible: a custom filesystem does not transfer files to the server. Config-only operations and borrowed-connection configuration workflows can use memory storage; managed native operations require the backend to access the chosen paths.
- Update internal constructors with explicit filesystem arguments and their existing callers, rather than introducing mutable setters or ambient storage.
- Keep Afero behind a private field in `filesystem.FS`, with only the operations used by application owners exposed. The public App constructor accepts Afero storage once; internal constructors receive the shared wrapper, whose zero value selects OS storage. This centralizes path checks and helper semantics without duplicating the full Afero interface or exposing its raw backend.
- Reject parent traversal outside BasePathFs before access, resolve output aliases through the same boundary, and retain temporary filenames in the virtual directory where they were created. Layered file-handle names are not valid inputs to another filesystem operation. Contextual filesystem errors retain their causes and consumers use `errors.Is`.

## Risks / Trade-offs

- Filesystem wrappers can expose different optional link capabilities → inspect only supplied capabilities and retain ordinary OS symlink E2E coverage.
- Incorrect propagation can silently mix storage → exercise full App config export/dump plus borrowed collection export/plan/apply/state-import with memory-only paths.
- Atomicity depends on the adapter's rename behavior → retain staged writes and cleanup, without promising stronger semantics than the filesystem.

## Decision verdict

Proceed: explicit Afero injection meets the requested boundary using an existing pin. Host defaults stay compatible, while shared path handling avoids repeating filesystem-specific logic. Server storage integration remains separate and is not represented as completed browser support.
