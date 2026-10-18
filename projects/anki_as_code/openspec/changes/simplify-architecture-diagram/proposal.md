## Why

The main architecture diagram mixes package dependencies with server hosting
and transport. A smaller runtime overview should explain App, the linked server,
and Anki storage before introducing package details.

## What Changes

- Replace the main diagram with App, its linked Rust server, self-launch, and
  gRPC connections to linked and external servers.
- Connect both server nodes to Anki storage containing Resources and State.
- Show Anki API beside Server in the linked library, with Server delegating
  storage access through the API.
- Connect App to an external Configs node with a Read/Write edge.
- Show the command entry point as `CLI` pointing to App inside the executable.
- Label the library call Run (C ABI) without adding an ABI node.
- Round the self-launch loop while preserving its endpoints and label; keep
  the external-server gRPC edge solid.
- Label the external node `anki_grpc / External server` and distinguish its
  planned App connection from current behavior.
- Retain the existing package graph as a separate detail diagram.
- Render a dark WebP preview through the existing site-themed Bazel rule.

## Capabilities

### New Capabilities

None. This documentation change sets `skip_specs: true`.

### Modified Capabilities

None. External-server support will be implemented in another PR.

## Impact

Only project documentation, Mermaid sources, and their render targets change.
There are no runtime, dependency, or deployment changes.
