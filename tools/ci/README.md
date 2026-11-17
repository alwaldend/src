---
title: Repository CI
description: Shared entry point for repository builds and tests
---

The [Forgejo workflow](../../.forgejo/workflows/ci.yaml) runs only for pushes to
`master`. It checks out the event revision with complete Git history and tags
for release metadata and secret scanning, including Git LFS assets, and uses
the local action:

```yaml
- name: Build and test the repository
  uses: ./tools/ci/action
```

The [action](action/action.yml) launches `bazel run --config=ci //tools/ci`
directly from Node.js with an argument array and no shell. It uses
`GITHUB_WORKSPACE` as the working directory, streams the command output, and
propagates failure. It needs no npm packages or generated JavaScript bundle.
The action comes from the exact revision already checked out by the workflow.
Checkout stays in the workflow because a local action must already exist in the
workspace before the runner can load it. Moving checkout inside this action
would require loading the action remotely with a separate source reference.

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

The runner needs the repository-managed Node.js, `bazel`, and Android
SDK/NDK prerequisites. The [runner deployment](../../infra/dev/ansible/README.md)
installs these before CI starts. Agent invocations use
`bazel run --config=ci //tools/ci`.

The pinned Bazel releases its invocation before executing the command, allowing
child builds to reuse each workspace's normal output base. No second cache or
source checkout is needed. Existing unit tests verify workspace discovery and
phase orchestration. The Linux executable acceptance fixture runs the compiled command
against a private Bazel process fixture, covering both phases in every workspace,
build/test failures, missing Bazel and invalid boundaries. Its repeatable
`ci-execution.json` receipt is retained in Bazel test outputs; it does not launch
real nested builds inside a sandbox.

The workflow's checkout action comes from the organization-owned
[Forgejo mirror](https://forgejo.alwaldend.com/alwaldend/com_github_actions_checkout),
at the existing v4.4.0 commit. Its immutable pin selects the action implementation;
`github.sha` separately selects the source revision being built.
