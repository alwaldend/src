---
title: Repository CI
description: Shared entry point for repository builds and tests
---

Every CI workflow checks out its event revision and uses the local action:

```yaml
- name: Build and test the repository
  uses: ./tools/ci/action
```

The [action](action/action.yml) launches `bazel run --config=ci //tools/ci`
directly from Node.js with an argument array and no shell. It uses
`GITHUB_WORKSPACE` as the working directory, streams the command output, and
propagates failure. It needs no npm packages or generated JavaScript bundle.
The action comes from the exact revision already checked out by the workflow.

The command builds
and tests normal `//...` targets in the root workspace and each nested workspace
with a `MODULE.bazel` at a boundary declared by `.bazelignore`. Non-workspace
entries such as caches and dependency directories are ignored. Workspace
boundaries remain owned by the root ignore file.

Each workspace runs `bazel_agent bazel build --config=ci //...` followed by
`bazel_agent bazel test --config=ci //...`. The command continues after failures
so all workspaces and both phases are attempted, streams Bazel diagnostics,
and exits nonzero if any phase fails. It does not run deployment targets or
runner acceptance checks. Normal Bazel expansion excludes manual targets and
incompatible platforms; optional configurations are not a build matrix.

The runner needs the repository-managed Node.js, `bazel`, `bazel_agent`, and Android
SDK/NDK prerequisites. The [runner deployment](../../infra/forgejo_runner/ansible/README.md)
installs these before CI starts. Agent invocations use
`bazel_agent bazel run --config=ci //tools/ci`.

The pinned Bazel releases its invocation before executing the command, allowing
child builds to reuse each workspace's normal output base. No second cache or
source checkout is needed. Unit tests verify workspace discovery, both phases,
and propagation of failures without running nested Bazel inside a test sandbox.

The workflow's checkout action comes from the organization-owned
[Forgejo mirror](https://git.alwaldend.com/alwaldend/com_github_actions_checkout),
managed by the [repository catalog](../../infra/repos/README.md) and Forgejo
Terraform. Its immutable commit pin selects the action implementation;
`github.sha` separately selects the source revision being built.
