## Why

The repository-specific openspec skill repeats generic create, apply, spec,
and archive steps already owned by the pinned upstream operation skills.
The user requested removing that duplication while keeping repository integration.

## What Changes

- Rename the local skill to repo-openspec and update routing and references.
- Limit the local skill to owner selection, pinned CLI invocation, repository
  work-state conventions, and explicit routing to upstream operation skills.
- Keep policy at AGENTS.md and generic artifact procedures upstream.
- Replace generic specification evaluation coverage with operation-routing cases.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- project-agents: Separate repository OpenSpec integration from upstream operations.

## Impact

The relocated tools/openspec skill and its offline evaluation configuration.
No upstream instructions, CLI pins, or change-record policy are changed.
