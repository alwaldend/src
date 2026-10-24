## Purpose

Allow caller-owned Emscripten classic MEMFS storage to serve both browser
language runtimes through their existing Node-compatible filesystem contract.

## ADDED Requirements

### Requirement: Independent initialized provider

The asynchronous factory SHALL return an initialized filesystem with callback
and synchronous operations over one MEMFS instance, retaining its module for
the provider lifetime and reserving descriptors zero through two.

#### Scenario: Two callers initialize storage

- **WHEN** two factories resolve and one creates a file
- **THEN** that file is absent from the other provider and ordinary descriptors
  are at least three

#### Scenario: Concurrent callers reuse options

- **WHEN** concurrent factories receive the same options and lifecycle hooks
- **THEN** their storage and nodes remain independent, each invocation runs its
  hooks, and initialization does not mutate the supplied options or hook arrays

### Requirement: Translate operations without replacing storage

The adapter SHALL preserve upstream file bytes, inode metadata, paths, stream
identity and directory behavior through direct FS operations. It SHALL expose
Node numeric flags, error codes and numeric or bigint Stats as requested.
Unavailable operations SHALL fail explicitly rather than claim success.

#### Scenario: Consumers interleave operations

- **WHEN** callback and synchronous calls perform positional binary transfers,
  append, truncation, rename replacement, links and directory operations
- **THEN** they observe the same upstream storage and open descriptors, with
  zero-filled file growth and conventional errors for invalid operations

#### Scenario: Exclusive creation encounters a symlink

- **WHEN** exclusive creation targets an existing symlink, including a dangling link
- **THEN** it fails with `EEXIST` without creating, replacing or truncating its target

#### Scenario: A write extent cannot be represented or allocated

- **WHEN** a write exceeds safe integer arithmetic or encounters an identified
  native storage allocation failure
- **THEN** it returns `EFBIG` or `ENOMEM` without altering file bytes or stream
  position, and unrelated programming exceptions remain distinguishable

#### Scenario: A callback read or write fails

- **WHEN** a callback read or write fails on an invalid or incompatible descriptor
- **THEN** its asynchronous callback receives the error, zero transferred bytes
  and the original buffer exactly once, and caller callback exceptions propagate

#### Scenario: An operation is absent upstream

- **WHEN** a caller requests an unsupported storage operation
- **THEN** it receives an explicit conventional error and no second backing
  store or synthetic inode emulates the operation

### Requirement: Reproducible provider evidence

The provider E2E SHALL use the real generated MEMFS and retain a JSON report
and binary output with hashes. It SHALL distinguish provider verification
from Go/Rust Chromium integration acceptance.

#### Scenario: Provider verification completes

- **WHEN** the standalone E2E runs
- **THEN** its artifacts identify exact adapter/module inputs, tested cases,
  observed limitations and output bytes
