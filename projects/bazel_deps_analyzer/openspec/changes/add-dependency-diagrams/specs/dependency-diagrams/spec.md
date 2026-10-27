## Purpose

Generate repeatable dependency diagrams from Bazel target graph JSON so build
relationships can be inspected and rendered with existing Mermaid tooling.

## ADDED Requirements

### Requirement: Target graph conversion

The converter SHALL consume Bazel `streamed_jsonproto` target records, select
the requested root and its reachable dependencies, and emit a Mermaid
flowchart with dependency arrows, full labels, and rule kinds. It SHALL emit
identical bytes regardless of record or dependency order and safely quote
labels. Generated files SHALL point to their generating rules when included.

#### Scenario: Shared dependencies

- **WHEN** two dependencies share a transitive dependency
- **THEN** the diagram contains one node for that target and both incoming edges

#### Scenario: Invalid input

- **WHEN** JSON is malformed, empty, has conflicting labels, or lacks the root
- **THEN** conversion fails with a contextual error and no partial diagram

### Requirement: Bounded graph selection

The converter SHALL support unlimited or nonnegative dependency depth, all
four flow directions, and optional inclusion of external targets and files.
Defaults SHALL omit files and external targets and use left-to-right flow.
It SHALL omit edges to targets absent from the selected graph.

#### Scenario: Depth boundary

- **WHEN** depth is one
- **THEN** the root and direct visible dependencies appear with no deeper nodes

#### Scenario: Optional external leaves

- **WHEN** the external-leaves option is selected
- **THEN** the first external target on each dependency path appears as a leaf
- **AND** traversal and outgoing edges stop at every external target
- **AND** local targets reached independently remain visible
- **AND** defaults and the existing external-expansion option remain unchanged

#### Scenario: Optional package grouping

- **WHEN** package grouping is selected
- **THEN** full package names label subgraphs and short target names label nodes
- **AND** external packages are nested under their external repository
- **AND** targets with equal short names retain separate identities and edges
- **AND** ungrouped diagrams remain the default

#### Scenario: Identify external dependencies

- **WHEN** external targets are visible
- **THEN** they all belong to one subgraph named `External`
- **AND** local targets remain outside that subgraph
- **AND** repository and package groups are nested inside it when enabled

### Requirement: Selectable ownership overview

The command and both rules SHALL accept `diagram_type` with `flowchart` as the
default and `mindmap` as an explicit alternative. Both SHALL select the same
targets using the existing build-input, depth, file, and external options.
Mindmap branches SHALL describe ownership rather than dependency edges, with
each selected target appearing once. Package grouping SHALL remain optional,
and all external targets SHALL share an `External` branch. Grouped external
targets SHALL retain their repository and full package hierarchy.

#### Scenario: Mindmap overview

- **WHEN** mindmap output with package grouping and external leaves is requested
- **THEN** short target captions belong to full package branches
- **AND** external packages belong to their repositories under `External`
- **AND** the diagram identifies itself as a build-input inventory
- **AND** shared dependencies and cycles do not duplicate or recurse over targets

#### Scenario: Preserve defaults and reject invalid types

- **WHEN** no diagram type is supplied
- **THEN** output is identical to explicit `flowchart` output
- **AND** an unsupported diagram type fails without creating or overwriting output

### Requirement: Build composition

The rules SHALL produce a Mermaid artifact from either a target or an existing
JSON artifact using declared actions and the repository toolchain. The target
macro SHALL capture JSON through `genquery` without recursively invoking Bazel.
It SHALL retain implicit non-tool inputs and exclude execution tools and
non-dependency attributes such as visibility.

#### Scenario: Real build graph

- **WHEN** a fixture target's diagram is built and rendered
- **THEN** the Mermaid and WebP reflect that fixture's actual dependency edges
