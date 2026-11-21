---
title: OpenSpec
description: Pinned OpenSpec CLI and repository specification validation
statuses:
  - active
---

This package runs the upstream [OpenSpec](https://openspec.dev/) CLI with
Bazel's pinned Node toolchain. `tools/rules_openspec/package.json` pins
`@fission-ai/openspec` to 1.11.0; `tools/rules_openspec/pnpm-lock.yaml` records the package
integrities and transitive dependency versions. Package lifecycle scripts are
disabled. No global npm installation is needed.

Each component owns an OpenSpec workspace at `<owner>/openspec/`. The
repository's evolution belongs to [`infra/src/openspec/`](../../infra/src/openspec/README.md),
which is the default for the runnable target. Run commands from the enclosing
repository root; `OPENSPEC_PROJECT` selects a different owner directory:

```sh
bazel run //tools/openspec -- --version
bazel run //tools/openspec -- list --specs --json
bazel run //tools/openspec -- status --change <change-name> --json
bazel run //tools/openspec -- validate --all --strict --no-interactive
OPENSPEC_PROJECT=tools/agents bazel run //tools/openspec -- list --specs --json
OPENSPEC_PROJECT=infra/vault bazel run //tools/openspec -- list --specs --json
```

`OPENSPEC_PROJECT` is a repository launcher setting, not an upstream CLI flag.
It is relative to the enclosing repository root, and OpenSpec reads and writes
that owner's source artifacts directly. Standalone nested projects use the
same root invocation, for example `OPENSPEC_PROJECT=tools/rules_skills`.

Run the workspace's sandboxed validation suite from the repository root.

```sh
bazel test //infra/src/openspec:validate_test
bazel test //tools/agents/openspec:validate_test
```

The suite runs the CLI in strict mode against declared workspace runfiles.
The current check validates specifications and active changes. The archive
check verifies that archived task lists are complete. Requirement scenarios
still need evidence from their owning implementation.

Gazelle generates the source target and validation suite from the workspace's
`openspec/config.yaml`. See the [OpenSpec Gazelle extension](../rules_openspec_gazelle/README.md).

The launcher disables telemetry and update checks. Tests isolate the CLI's
XDG config and data paths in Bazel's test temporary directory. Commands that
create, edit, or archive changes modify source artifacts as documented by the
[upstream CLI](https://github.com/Fission-AI/OpenSpec/blob/v1.11.0/docs/cli.md).

Update the version in `tools/rules_openspec/package.json`, then use the owning pnpm workflow
described in [`tools/js`](../js/README.md) to regenerate its lock with lifecycle
scripts disabled. Re-run the version command and the workspace validation suite after
upgrading.
