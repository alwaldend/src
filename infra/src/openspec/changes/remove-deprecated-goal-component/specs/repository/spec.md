## REMOVED Requirements

### Requirement: Legacy goal migration and deprecation

**Reason**: The goal tool was deprecated and its replacement, the pinned
OpenSpec CLI, is in active use. The component, its skill sources, its landing
site, and its DNS record have been removed rather than retained for
compatibility.

**Migration**: Continue to use the affected owner's `openspec/` workspace
through `bazel_agent bazel run //tools/openspec -- ...`. The migration record
at [migration.md](../../../../migration.md) and the preserved history under
`projects/agents/openspec/changes/archive/` remain the historical record. The
removed bytes stay retrievable from git history at the removal commit's
parent.

## ADDED Requirements

### Requirement: Removed legacy goal record

The repository MUST NOT carry a goal tool, goal store, or goal skill, and MUST
NOT reintroduce a goal record format as a maintained work channel. Maintained
work SHALL use the owning component's OpenSpec workspace. The change SHALL
identify the removed component.

#### Scenario: An agent selects durable-work tooling

- **WHEN** an agent needs durable intent, tasks, or continuation state
- **THEN** it uses the owning component's `openspec/` workspace through the
  pinned CLI
- **AND** no goal command, skill, or store is available or documented as an
  active channel

#### Scenario: A reader needs the removed implementation

- **WHEN** a reader needs a file from the removed goal component
- **THEN** the change and migration record name the removed tree
- **AND** the reader retrieves the file from git history at the removal
  commit's parent
