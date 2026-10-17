## Context

See proposal.md for scope. The existing layout and codec packages already own file discovery and serialization; the facade can hold the caller-selected root without changing those implementations.

## Goals / Non-Goals

Make all consumer-facing configuration operations methods of one root-bound Store. Preserve current CLI behavior, file bytes, error policy and validation boundaries.

## Decisions

New(root string) constructs a Store without filesystem access. Load() and Dump(format) accept the existing directory or root-file input through the stored path. Export(configs) writes beneath the stored directory; file input remains a read-only convenience. Keep the root private and return fresh parsed messages on every Load instead of caching or retaining configs.

Remove the package-level operations instead of maintaining duplicate entry points for an unreleased internal API. Keep Version for its existing cross-package serialization contract. Put Store and its methods in one source file, retaining layout and codec as implementation dependencies.

## Risks / Trade-offs

Missed consumers would fail to build. Migrate every master caller and adapt existing E2E workflows before implementation, including reusing a store to export and reload a tree. Verify config dump, native export and state import through their existing workflows. Record the prerequisite for PR #128 instead of including its unmerged implementation.
