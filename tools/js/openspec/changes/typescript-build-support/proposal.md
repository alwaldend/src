## Why

The diagram renderer and snapshot commands will use TypeScript. The repository
currently has JavaScript runtime/package rules but no compiler or Node/Three.js
declarations. Establish that shared build support in a separate prerequisite PR.

## What Changes

- Add `aspect_rules_ts` 3.10.1 alongside the existing `rules_js` 2.8.3.
- Pin TypeScript 7.0.2, Node types 24.19.2, and Three.js types 0.186.0 in the
  existing tools pnpm workspace; derive the compiler toolchain from that pin.
- Exercise strict compilation, declaration dependencies, Node execution,
  browser execution, and rejection of a real type error end to end.
- Document direct use of the maintained `ts_project` rule.

## Capabilities

### New Capabilities

None; this is shared build tooling with no product behavior delta.

### Modified Capabilities

None. Diagram migrations remain separate deliveries.

## Impact

The user explicitly approved these four pinned dependencies and their locked
upstream transitive dependencies on 2026-11-05. Dependency declarations and
generated locks stay with their existing owners; build guidance and acceptance
belong to tools/js. No runtime diagram code, host configuration, or deployment
is included.
