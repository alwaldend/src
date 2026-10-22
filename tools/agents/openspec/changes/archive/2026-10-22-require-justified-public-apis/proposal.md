## Why

The Go skill does not require a reason for exporting declarations, and the Python skill limits its private-by-default rule to attributes. Helpers can consequently become accidental public APIs.

## What Changes

- Require a concrete consumer or contract before exposing Go or Python declarations.
- Keep implementation details private and Python package re-exports explicit.
- Add evaluation cases covering unnecessary exports and required public contracts.

## Capabilities

Agent guidance only; no product behavior or specification changes. This change uses `skip_specs: true`.

## Impact

Update the canonical Go and Python skills and their evaluation fixtures. Maintain this shared agent-guidance change under `tools/agents`, following the existing code-documentation change. Application source, dependencies, and generated code remain unchanged.
