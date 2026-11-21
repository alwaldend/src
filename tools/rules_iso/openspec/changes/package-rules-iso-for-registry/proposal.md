## Why

A registry can expose the ISO ruleset independently without maintaining
a nested Bazel workspace for development.

## What Changes

- Assemble the ruleset's existing public packages into a source archive.
- Supply the module declaration for registry publication.

## Capabilities

### Modified Capabilities

- `project-rules-iso` gains an independently consumable Bzlmod archive.

## Impact

The existing extension label and flash target remain available in the archive.
