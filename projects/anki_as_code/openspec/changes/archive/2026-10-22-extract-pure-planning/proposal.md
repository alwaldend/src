## Why

Planning is coupled to a native mutation engine even though it only compares supplied values. Extract a pure planning API so callers can load once and plan without starting or retaining a backend.

## What Changes

- Add `plan.Build(configs, snapshot, state)`, returning a protobuf execution plan and path-based report without I/O or input mutation.
- Keep resolution models private and compare managed resources directly, preserving untracked resources and filtered cards.
- Move config moves and proposed bindings into private planning logic; remove the state preview API and duplicate Anki name/structure validation.
- Keep native application and CLI integration in #128; publish this library independently against master.

## Capabilities

### New Capabilities

- `collection-planning`: Construct a deterministic read-only plan from desired configs, a native resource snapshot, and stored bindings.

### Modified Capabilities

None.

## Impact

Affects internal/plan, state preview, plan protobuf documentation, native-server fixture visibility, and native-backed planning E2E coverage. Uses existing Anki protobufs and resource snapshots. Adds no dependencies, storage writes, lifecycle ownership, or CLI commands.
