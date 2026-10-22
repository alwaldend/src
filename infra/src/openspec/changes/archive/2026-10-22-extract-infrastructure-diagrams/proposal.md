## Why

Continue PR #141 by converting every remaining infrastructure Drawio page to project-owned Mermaid documentation, then remove `infra/arch`.

## What Changes

- Preserve labels, containers, edge directions, protocols, and historical classification in Mermaid sources.
- Keep DC1 with Mikrotik, service diagrams with their infrastructure owners, and www with the website project.
- Keep archived Proxmox and Harbor diagrams with their existing owners. Put archived Harvester, Opencode, and Hermes diagrams without projects in `infra/archive`.
- Continue the 3x-ui source from PR #141 as one diagram with separate overview, control-plane, and data-plane sections. Render documentation assets through maintained site Mermaid targets.
- Remove the obsolete Drawio project, its validation registration, and stale operational references.

## Capabilities

Documentation relocation and format conversion only; `skip_specs: true`. No infrastructure behavior changes.

## Impact

Infrastructure documentation, the website project's historical architecture documentation, and the repository OpenSpec validation map. No live operations or external dependencies.
