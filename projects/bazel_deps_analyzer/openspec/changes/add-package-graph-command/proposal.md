## Why

Large dependency captures need a reproducible package-level overview. The
current command emits target-level graphs, so the requested top-package diagram
has required task-local construction outside the analyzer.

## What Changes

- Add `analyze package` to aggregate every captured target
  by a positive `--package-depth` (default one), the main root package, or
  external repositories, without requiring a synthetic selection root.
- **BREAKING**: move existing target conversion flags to `analyze target`.
  The bare command and bare `analyze` print help and return success. Update the
  existing Bazel rules and provide a packages-from-JSON rule.
- Count targets and internal dependencies, deduplicate inter-package
  connections, omit edge captions, and combine reciprocal connections into
  double-headed arrows by default. Retain an option for separate arrows.
- Generate the existing automatic ELK layout directive, preserving all
  selected dependency relationships and target-conversion semantics.
- Export the same aggregate as JSON with sorted member labels so the viewer
  consumes analyzer-owned group construction rather than duplicating it.
- Document capture limits and verify the command through E2E fixtures and
  the existing real capture and dark Mermaid rendering pipeline.
- Add opt-in `--build-deps-only` to both analysis modes and
  `build_deps_only = False` to capture and JSON-input rules. Keep only explicit
  dependency-label values in `deps` and `srcs`, intersected with captured rule
  inputs, plus generated-file producer links. Defaults remain unchanged.

## Capabilities

### New Capabilities

- `package-graphs`: Deterministic construction of aggregated package dependency
  graphs from Bazel streamed JSON captures.

### Modified Capabilities

None in the archived specification tree. Default flowchart and mindmap
conversion semantics remain unchanged under `analyze target`.

## Impact

The analyzer's CLI, internal graph package, rules, E2E coverage, and README.
The implementation reuses the current parser, Cobra, and Mermaid renderer.
No external dependencies or live infrastructure changes are needed.
