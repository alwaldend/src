## Why

First-party components still group source under `main/<language>` or
`private`, obscuring the roles defined by the repository layout policy.

## What Changes

- **BREAKING**: Move legacy source in `projects/` and `tools/` into role-based
  directories and update all maintained consumers without compatibility paths.
- Separate command entry points, internal implementation, schemas, rules,
  assets, and patches according to their responsibilities.
- Update imports, Bazel labels, packaging, Gradle source sets, and documentation.

## Capabilities

No behavioral requirements change. This structural refactor uses
`skip_specs: true`; existing component specifications remain authoritative.

## Impact

Nine product projects and legacy tool components, including standalone rule
modules and their consumers. No new dependencies, deployment, host configuration
changes, or merge. Acceptance requires the old directories to be absent,
maintained references to resolve, representative output inspection, affected
build checks, and the repository delivery gates.
