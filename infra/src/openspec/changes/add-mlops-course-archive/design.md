## Context

See proposal.md for scope. On 2026-10-11 upstream HEAD resolved to
74dc45062f43f9b4abfbca738ffa6b56db4d5daa. The downloaded archive contains
the MIT LICENSE, README.md, and lesson directories.

## Goals / Non-Goals

Expose the upstream source tree as Bazel inputs. Installing course dependencies
or running lessons is outside this dependency registration.

## Decisions

Use the existing http_archive pattern with the full commit in the codeload URL,
SRI integrity, and matching strip_prefix. A moving branch would not reproduce
the source; a tag is unavailable. Keep the pin in include.MODULE.bazel.
Provide a filegroup for all archive files and a repository-internal source alias.
Preserve upstream licensing and provenance in the package README.

## Risks / Trade-offs

GitHub-generated archive bytes can change: integrity rejects unexpected bytes.
The checksum is computed locally; upstream supplies no independent checksum.

## Execution

Current candidate: source archive SHA-256
dcd4591cc6e8c74b5c3c3a1688c24bd5de9049e193b71d20f0047226588a81ab.
Acceptance: the source target builds; all 59 upstream source files match the
downloaded archive byte for byte. Module reconciliation produced no lock diff.
The repository quality suite, buildifier, module freshness, OpenSpec validation,
and affected-package semantic lint passed. Next: prepare and validate the exact
delivery candidate, then publish. No live infrastructure changes.
