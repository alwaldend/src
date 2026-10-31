## Context

The generated classic MEMFS factory exposes `FS` and `ERRNO_CODES`. Go uses
error-first callback methods; Rust uses synchronous Node-style calls. The
caller independently passes the same provider to each consumer.

## Goals / Non-Goals

Translate those APIs without implementing storage, inode allocation or a
descriptor table. Toolchain pinning and actual Go/Rust Chromium acceptance
are owned by sibling changes. IndexedDB persistence is outside this provider.

## Decisions

- `createEmscriptenFilesystem(options?)` awaits the generated default factory
  and closes over its module. Each invocation copies the options object and
  lifecycle arrays, retaining function identities and module receivers. This
  prevents upstream module mutation or hook consumption from sharing storage.
- Use FS stream and node operations for descriptor-based metadata and mutation,
  including after rename/unlink. Resolving descriptor paths again would target
  a replacement file and violate open-file identity.
- Translate numeric flags, timestamps, Stats shapes and errors at the boundary.
  Callback wrappers use microtasks and the same synchronous operations.
- Hard links require an upstream primitive; do not emulate them by copying bytes
  or maintaining shadow inodes. Unsupported behavior is reported explicitly.
- Exclusive creation also sets upstream no-follow, so the upstream existence
  check sees the final symlink before it can create a dangling target.
- Writes validate the effective sequential, positioned or append extent with
  safe integer arithmetic. Native array-length or allocation failures during
  regular-file growth become `EFBIG` or `ENOMEM`; unidentified exceptions pass
  through. Rejected writes retain the upstream stream position. There is no
  additional storage quota; JavaScript integer and engine allocation limits apply.
- Read/write error callbacks retain the Node tuple `(error, 0, originalBuffer)`.
  Caller callback exceptions are not treated as filesystem failures.

## Risks / Trade-offs

- MEMFS lacks some POSIX semantics: preserve and document actual upstream
  metadata and reject unavailable operations without fabricating support.
- In-memory fsync does not imply durability: document that completion only
  confirms writes to MEMFS; snapshots remain caller-owned.
- Build integration is concurrent: provider tests must wait for the real
  generated module, and the parent composes shared BUILD changes.
