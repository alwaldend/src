## Why

Several Go tools register command state during import, discard filesystem failures, or run child processes without the caller's context. Package documentation coverage also varies across tools.

## What Changes

Audit first-party Go packages under tools against repo-go and correct the deviations while preserving the flags, interfaces, and build contracts their callers use.

## Impact

The changes affect repository development tools and their Bazel targets.
