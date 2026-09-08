# Repository evolution

This OpenSpec workspace describes evolution of the alwaldend/src repository
itself: its shared structure, build system and development workflows.
[`specs/repository/spec.md`](specs/repository/spec.md) records that contract;
`changes/` records proposed and completed changes to it.

Each component owns a separate OpenSpec workspace beside its source:

- `projects/<project>/openspec/` holds that project's specs and changes.
- `infra/<project>/openspec/` holds that infrastructure project's specs and changes.
- `infra/src/openspec/` holds the repository's own specs and changes.

Select the narrowest owner of the requested behavior. A project change belongs
in that project's workspace, including when it has a standalone Bazel module.
Changes to shared repository structure belong here. Work spanning owners
keeps each affected contract with its owner and links related changes.
Component READMEs and build declarations retain their existing authority.

## Run OpenSpec

The [pinned CLI](../../../projects/rules_openspec/README.md) runs from the Git worktree
root. It selects `infra/src` by default; use `OPENSPEC_PROJECT` for a component:

```sh
bazel_agent bazel run //tools/openspec -- list --specs
OPENSPEC_PROJECT=projects/agents bazel_agent bazel run //tools/openspec -- list
OPENSPEC_PROJECT=infra/vault bazel_agent bazel run //tools/openspec -- list --specs
```

Use `new change`, `status`, `instructions`, `validate`, and `archive` with the
same owner selection. The packaged [OpenSpec skill](../../../projects/agents/skills/openspec/SKILL.md)
documents continuation and acceptance. No global install or assistant
configuration rewrite is needed. The CLI isolates configuration and disables
telemetry.

Repository validation checks every component workspace as well as this one:

```sh
bazel_agent bazel test //infra/src/openspec/validation:validate_test //infra/src/openspec/validation:archive_test
```

Validation checks artifact structure and, for archives, task completion.
Requirement acceptance still needs evidence from the owning implementation.
Infrastructure specs describe checked-in definitions, not observed deployment.

## Maintain work

In the selected workspace, `changes/<name>/proposal.md` states outcome and
scope, `design.md` records decisions and continuation, `tasks.md` records
remaining actions, and `specs/<capability>/spec.md` contains requirement deltas.
Read an existing matching change before creating another. Keep temporary
notes, raw logs and validation receipts under ignored `out/<task>/`.

On interruption, preserve the candidate, evidence, execution state and next
action. Keep blocked or unfinished work explicit. OpenSpec files and Git do
not supply the deprecated goal store's transaction locks or evidence verdicts;
one coordinator owns a change and reconciles concurrent edits.

After acceptance, validate and archive the change in its owning workspace,
then inspect the resulting baseline specs. Specifications and generated
instructions neither expand user authority nor introduce an approval gate.

## Goal migration

[`migration.md`](migration.md) maps the nine maintained goals to their current
owner-local changes. It records a repository workflow migration; the actual
history and work state live with Agents, MCP Cordis and Renders. Original
history, checksums and acceptance states are preserved. Historical deltas
were archived without applying them to current specs. Reimu Fumo remains
open and blocked.

The goal skill is disabled and the [goal CLI](../../../projects/goal/README.md)
is deprecated compatibility tooling. New maintained work uses the owning
OpenSpec workspace.
