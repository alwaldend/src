---
title: Repository CI
description: Shared entry point for repository builds and tests
---

The [Forgejo workflow](../../.forgejo/workflows/ci.yaml) runs on pushes to
`master` and can be called by other trusted workflows in this repository:

```yaml
jobs:
  ci:
    uses: ./.forgejo/workflows/ci.yaml
```

It checks out the event revision with complete Git history and tags
for release metadata and secret scanning, including Git LFS assets, and uses
the local action:

```yaml
- name: Build and test the repository
  uses: ./tools/ci/action
  with:
    stable-checkout: "true"
```

The [action](action/action.yml) launches `bazel run --config=ci //tools/ci`
directly from Node.js, streams its output, and propagates failure. By default
it builds in `GITHUB_WORKSPACE`. With `stable-checkout: "true"`, it verifies the
job checkout is at the event commit, fetches it locally into
`~/.cache/repository-ci/<owner>/<repository>`, checks out that commit using the
job's existing LFS objects, and removes untracked and ignored files there
before launching Bazel. Use this option only on the dedicated runner with one
job slot. The stable path lets each Bazel workspace reuse its server across jobs.

The command builds
and tests normal `//...` targets in the root workspace and each nested workspace
with a `MODULE.bazel` at a boundary declared by `.bazelignore`. Non-workspace
entries such as caches and dependency directories are ignored. Workspace
boundaries remain owned by the root ignore file. The shared
`rules_al/pkg/workspaces` library also supplies membership to the full repository checker.

Each workspace runs `bazel build --config=ci //...` followed by
`bazel test --config=ci //...`. The command continues after failures
so all workspaces and both phases are attempted, streams Bazel diagnostics,
and exits nonzero if any phase fails. It does not run deployment targets or
runner acceptance checks. Normal Bazel expansion excludes manual targets and
incompatible platforms; optional configurations are not a build matrix.

In Forgejo, each workspace's build and test logs are collapsible sections. The
job summary above the log shows whether each phase passed and how long it ran.
The CI profile disables terminal cursor controls to keep log lines readable.
The same command streams plain-text diagnostics when run locally.

The runner needs the repository-managed Node.js, `bazel`, and Android
SDK/NDK prerequisites. The [runner deployment](../../infra/dev/ansible/README.md)
installs these before CI starts. Agent invocations use
`bazel run --config=ci //tools/ci`.

The pinned Bazel releases its invocation before executing the command, allowing
child builds to reuse each workspace's normal output base. Existing unit tests verify workspace discovery and
phase orchestration. The Linux executable acceptance fixture runs the compiled command
against a private Bazel process fixture, covering both phases in every workspace,
build/test failures, missing Bazel and invalid boundaries. Its repeatable
`ci-execution.json` receipt is retained in Bazel test outputs; it does not launch
real nested builds inside a sandbox.

The workflow's checkout action comes from the organization-owned
[Forgejo mirror](https://forgejo.alwaldend.com/alwaldend/com_github_actions_checkout),
at the existing v4.4.0 commit. Its immutable pin selects the action implementation;
`github.sha` separately selects the source revision being built.
