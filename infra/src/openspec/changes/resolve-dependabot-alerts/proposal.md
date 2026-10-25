## Why

GitHub reports 250 open Dependabot alerts across ten repository dependency
manifests on 2026-10-25. Resolve vulnerable dependency selections while preserving
the affected tools and reproducible Bazel workflows.

## What Changes

- Upgrade affected Python, npm, and Go dependencies and regenerate their locks
  through their owning workflows.
- Preserve consumer compatibility, including Promptfoo telemetry and serialized
  execution patches and the Mermaid rendering patch.
- Remove the unpatched ecdsa dependency from unsigned DFU packaging and
  verify that all captured advisory ranges are absent from the final selections.

## Capabilities

No specification requirements change; this is dependency maintenance.
`skip_specs: true` applies.

## Impact

Root Go dependencies, `tools/pkgsite`, the Python environments in `tools/py`,
`tools/ansible`, `tools/ykman`, and `projects/infinitime`, and npm dependencies in
`tools`, `tools/rules_promptfoo`, `tools/rules_openspec`, and
`projects/mcp_cordis`. Acceptance requires resolving patchable baseline alerts,
checking generated locks against the captured advisories, and validating
representative consumers. Live infrastructure changes and alert dismissal are
outside scope.
