## Why

Dependency diagrams should follow the actual Bazel graph instead of a maintained
list of edges. The site needs a repeatable Mermaid and WebP view of its build.

## What Changes

- Add `projects/bazel_deps_analyzer`, a Go converter and Bazel rules consuming
  the existing `native.genquery` streamed JSON format.
- Provide deterministic target diagrams with depth, file, and external filters.
- Integrate companion graph and WebP targets into alwaldend.com and register
  the project landing and reference documentation.

## Capabilities

### New Capabilities

- `dependency-diagrams`: Convert a target's JSON dependency graph to Mermaid.

### Modified Capabilities

None.

## Impact

New project, project registry, and site build targets. Existing Cobra and
Mermaid dependencies suffice. Site integration has a separate owner-local
change record. No live deployment or new dependency is needed.
