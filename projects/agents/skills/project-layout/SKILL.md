---
name: project-layout
description: >-
  Choose and review repository paths by project ownership and source role.
  Use before creating or moving source, adding a project, or deciding a
  directory layout anywhere in the repository.
---

# Organize source by role

Apply this policy across the repository, not only under `projects/`.

## Choose the owning project

Identify the narrowest component that owns the source. A project root can live
under any repository boundary, such as `projects/<project>`,
`tools/<project>`, or `infra/<project>`; do not treat a top-level boundary as
the owner when a narrower component owns the source.

Keep the boundary rules from the root `AGENTS.md` and the nearest `README.md`.
If an owner does not exist, create a named project inside the appropriate
boundary before choosing its internal layout.

Prefer project-local ownership. Store skills, metadata, docs, and other
supporting content in the project that owns the behavior they support. Avoid
creating a separate generic project that aggregates unrelated content merely
by type, such as a standalone `skills` project; a shared project is justified
only when its members have a common owner, contract, or consumer, not just a
common directory role.

## Choose a source type

Place project-owned source at `<project>/<type>/<name-or-tree>`. Choose
`<type>` for the source's role, not its implementation language:

- `cmd/<command>` for executable entry points.
- `internal/<tree>` for implementation that is not a supported consumer
  surface.
- `pkg/<library>` for deliberately reusable implementation.
- `skills/<skill>` for agent skills.
- `api/<name>`, `docs/<name-or-tree>`, `test/<name-or-tree>`, and
  `assets/<name-or-tree>` when those roles are distinct project content.

Use ecosystem-required role directories when the relevant tool mandates a
different spelling or placement. Keep ordinary unit tests beside their source
when that is the language convention; use `test/` for separately owned test
support or suites.

Name every internal implementation directory `internal`, including one nested
below another source type. Never introduce `private` as a directory name for
this purpose.

## Reject language-first layers

Do not add a `main/<language>` layer or another directory whose only purpose
is grouping source by implementation language. For example, prefer
`<project>/cmd/server` over `<project>/main/go/cmd/server` and
`<project>/internal/parser` over `<project>/main/python/parser`.

Existing `main/<language>`, `private`, and other legacy layouts are
grandfathered. They are not precedent for new paths. Do not broaden a focused
change into an unrelated migration; when the task already moves the affected
source, place it in the current role-based layout and update its references.

## Reference the conventions selectively

The official Go guide documents the `cmd` and `internal` roles used here:
[Organizing a Go module](https://go.dev/doc/modules/layout). The
[community project-layout catalog](https://github.com/golang-standards/project-layout)
is explicitly not an official standard. Use it only as a catalog of possible
roles, and copy only the directories the project actually needs.
