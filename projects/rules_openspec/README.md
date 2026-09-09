---
title: OpenSpec
description: Pinned OpenSpec CLI and repository specification validation
statuses:
  - active
---

This package runs the upstream [OpenSpec](https://openspec.dev/) CLI with
Bazel's pinned Node toolchain. `tools/package.json` pins
`@fission-ai/openspec` to 1.11.0; `tools/pnpm-lock.yaml` records the package
integrities and transitive dependency versions. Package lifecycle scripts are
disabled. No global npm installation is needed.

Each component owns an OpenSpec workspace at `<owner>/openspec/`. The
repository's evolution belongs to [`infra/src/openspec/`](../../infra/src/openspec/README.md),
which is the default for the runnable target. Run commands from the enclosing
repository root; `OPENSPEC_PROJECT` selects a different owner directory:

```sh
bazel_agent bazel run //tools/openspec -- --version
bazel_agent bazel run //tools/openspec -- list --specs --json
bazel_agent bazel run //tools/openspec -- status --change <change-name> --json
bazel_agent bazel run //tools/openspec -- validate --all --strict --no-interactive
OPENSPEC_PROJECT=projects/agents bazel_agent bazel run //tools/openspec -- list --specs --json
OPENSPEC_PROJECT=infra/vault bazel_agent bazel run //tools/openspec -- list --specs --json
```

`OPENSPEC_PROJECT` is a repository launcher setting, not an upstream CLI flag.
It is relative to the enclosing repository root, and OpenSpec reads and writes
that owner's source artifacts directly. Standalone nested projects use the
same root invocation, for example `OPENSPEC_PROJECT=projects/rules_skills`.

Run upstream validation for all 51 owner workspaces in sandboxed tests against
declared source inputs:

```sh
bazel_agent bazel test //infra/src/openspec/validation:validate_test //infra/src/openspec/validation:archive_test
```

`validate_test` aggregates one strict test per owner for current specs and
active change deltas. `archive_test` covers the repository, agent-system, and
MCP workspaces and checks that archived changes have completed task
checkboxes; upstream does not reapply historical deltas during that check.
Both are structural checks. Requirement scenarios still need evidence from
their owning implementation and validation workflow.

The source-to-test mapping lives in
`infra/src/openspec/validation/BUILD.bazel`. Each owner's
`openspec/BUILD.bazel` exposes config, baseline specs, and complete native
change artifacts. Standalone projects expose this source through their
existing module repository, without taking an OpenSpec dependency. The shared
`openspec_validation` macro retains declared source runfiles through a
`js_library` with `no_copy_to_bin`: files keep their owning package paths,
including across module boundaries. Tests select that owner inside their
sandbox and receive only declared inputs.

When adding an owner workspace, add its source label to `_WORKSPACE_SOURCES`.
When it first archives a change, also add it to `_ARCHIVE_WORKSPACES`. Keep
these aggregate memberships aligned with the checked-in workspaces.

The launcher disables telemetry and update checks. OpenSpec's XDG config and
data paths are isolated to ignored `out/openspec/` in the source workspace for
`bazel run`, and to Bazel's test temporary directory for tests. These targets
do not read or change user-global OpenSpec configuration. Commands that create,
edit, or archive changes still modify source artifacts as documented by the
[upstream CLI](https://github.com/Fission-AI/OpenSpec/blob/v1.11.0/docs/cli.md).

Update the version in `tools/package.json`, then use the owning pnpm workflow
described in [`tools/js`](../js/README.md) to regenerate its lock with lifecycle
scripts disabled. Re-run the version command and both validation tests after
upgrading.
