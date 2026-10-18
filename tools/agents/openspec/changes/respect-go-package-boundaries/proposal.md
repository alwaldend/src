## Why

Documenting package purposes is insufficient when implementations bypass their stated owners. The Go skill should require reading and respecting package contracts while keeping those contracts accurate.

## What Changes

- Require responsibility, side-effect and dependency ownership contracts in Go package documentation.
- Require implementation and review to honor those contracts and identify drift instead of broadening comments to excuse misplaced logic.
- Extend the existing offline evaluation cases for command/application boundaries and inaccurate package comments.

## Capabilities

### New Capabilities

- `go-package-boundaries`: Treat Go package responsibility comments as implementation and review contracts.

### Modified Capabilities

None.

## Impact

Canonical repo-go skill and its existing Promptfoo evaluation files. No provider, dependency or host configuration changes.
