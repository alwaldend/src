## Purpose

Provide reproducible, readable package dependency overviews directly from
captured Bazel target records without manual graph construction.

## ADDED Requirements

### Requirement: Explicit dependency attribute filtering

Both analysis modes SHALL offer `--build-deps-only`, default false, selecting
explicit `deps` and `srcs` dependency label values intersected with captured
rule inputs before traversal, depth, or aggregation. Nodep, aspect, selector
condition, package-group, and environment-group metadata edges SHALL be
excluded. Generated-file producer edges SHALL remain. All-record package node
selection SHALL remain unchanged. Capture and JSON-input Bazel rules SHALL
expose `build_deps_only = False`. Filtered Mermaid SHALL identify
`build_deps_only=true`; filtered package JSON SHALL add `meta.buildDepsOnly`,
omitting that field by default.

#### Scenario: Complete capture metadata

- **WHEN** filtering is requested
- **THEN** every RULE must contain a `name` attribute marker and selected
  dependency values must use supported, well-formed label encodings.
- **AND** missing metadata or conflicting records fail before output replacement
  with recapture guidance including the required query flags.
- **AND** the capture macro requests name/deps/srcs, default values, unflattened
  selectors and rule inputs/outputs, and disables implicit, tool and nodep edges.

#### Scenario: Structural scope and ambiguity

- **WHEN** a capture contains ordinary tools, data, visibility, alias actual,
  or custom dependency attributes
- **THEN** those attributes do not supply selected dependencies, while a label
  also present in valid explicit deps/srcs and ruleInput remains selected.
- **AND** documentation explains that upstream flags perform tool exclusion,
  collapsed rule inputs cannot establish per-attribute tool provenance, the
  name marker cannot prove capture flags, and the mode is not a complete
  semantic build classifier.

#### Scenario: Compatibility

- **WHEN** filtering is disabled
- **THEN** existing attribute-less captures and default output remain unchanged.

### Requirement: Separate package graph command

The analyzer SHALL expose `analyze package`, accepting streamed JSON
through `--input` and Mermaid output through `--output`, with standard streams
as defaults. It SHALL process all captured targets without a required root.
The existing target conversion SHALL move to `analyze target` with its
existing flags, conversion behavior, and defaults. The bare command and bare
`analyze` SHALL print help and return success without reading input. Root-level
conversion flags SHALL NOT remain available. Existing Bazel rules SHALL invoke
`analyze target`,
and a dedicated JSON-input rule SHALL expose package graph construction.

#### Scenario: Disconnected captured targets

- **WHEN** a capture contains several disconnected components
- **THEN** the package view includes the owners of every captured target
  without adding a synthetic root or requiring `--root`.

#### Scenario: Command selection

- **WHEN** the command is invoked bare or as bare `analyze`
- **THEN** it prints help and exits successfully even if stdin is malformed.
- **WHEN** target conversion is requested
- **THEN** `analyze target --root <label>` preserves its selection behavior, while
  root-level `--root` and package-level `--depth` are rejected.

### Requirement: Deterministic ownership aggregation

The package view SHALL aggregate local targets under the first N package path
segments selected by `--package-depth=N`, which MUST be a positive integer and
default to one. The main root SHALL have caption `root`; local paths SHALL
retain the `//` prefix, and external repositories SHALL remain one `External`
node at every depth. Files and external targets SHALL be
included by default and independently suppressible. Node captions SHALL report
selected target and internal-edge counts. Ordering and duplicate input records
SHALL NOT change the output.

#### Scenario: Shared and internal dependencies

- **WHEN** several targets in one top-level package depend on another package
  and on each other
- **THEN** one node represents each owner, one directed relationship represents
  each distinct owner pair, and internal edges are counted without self-loops.

#### Scenario: Root and external names

- **WHEN** the main root, a local `external` package, and external repositories
  appear in one capture
- **THEN** all three ownership identities remain distinct.

#### Scenario: Configurable package hierarchy

- **WHEN** targets occur in `//tools/compiler/frontend` and
  `//tools/compiler/backend`
- **THEN** depth one produces `//tools`, depth two produces `//tools/compiler`,
  and depth three separates the frontend and backend packages.
- **AND** shorter package paths retain their full name, `root` remains
  distinct from `//root`, and `External` remains distinct from `//External`.
- **AND** zero or negative package depth fails without changing output.

### Requirement: Readable connections without lost direction

The output SHALL omit edge captions and use automatic Mermaid layout. By
default reciprocal package relationships SHALL share a double-headed arrow;
`--merge-reciprocal=false` SHALL retain separate directed arrows. Both forms
SHALL represent exactly the selected inter-package relationships.

#### Scenario: Mutual dependency

- **WHEN** package A depends on B and B depends on A
- **THEN** default output uses a double-headed connection, and disabling
  merging emits both original directions without inventing another relation.

### Requirement: Honest capture boundaries and safe failures

The output SHALL identify itself as an aggregate of captured data, report
unresolved dependency references, and state that aggregate paths do not imply
target-level reachability. Labels SHALL remain escaped content. Invalid input,
invalid ownership labels, and invalid command options SHALL fail before
creating or overwriting output.

#### Scenario: Missing dependency record

- **WHEN** an input target refers to a record absent from the capture
- **THEN** the report counts the unresolved reference and emits no invented
  target or dangling connection.

#### Scenario: Malformed capture

- **WHEN** malformed or conflicting target records are supplied
- **THEN** the command fails and preserves any existing output file.

### Requirement: Shared JSON aggregation contract

The package command SHALL support `--format=mermaid` by default and
`--format=json`. Both serializers SHALL consume the same constructed graph.
JSON SHALL contain sorted `groups` with `caption`, `targetCount`,
`internalEdges`, and sorted normalized `members`; sorted directed `edges`
with `source`, `dependency` group-array indices and `targetEdges` counts;
and `meta` with `packageDepth`, `selectedTargets`, `selectedEdges`,
`unresolvedReferences`, `includeFiles`, and `includeExternal`.
JSON SHALL retain both directions of reciprocal relationships regardless of
Mermaid's arrow-merging option. Invalid formats SHALL fail before output changes.

#### Scenario: Viewer consumes analyzer membership

- **WHEN** JSON is requested at a supported positive package depth
- **THEN** every selected target belongs to exactly one emitted group, and
  groups and directed edges equal the same capture's Mermaid representation
  after expanding double-headed arrows.
- **AND** changing only `--merge-reciprocal` leaves JSON unchanged.
