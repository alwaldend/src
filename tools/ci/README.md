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

The [action](action/action.yml) uses `bazel run --script_path` to build
`//tools/ci` with the CI profile and produce a launcher. Node.js runs that
launcher in a separate stage, streams output, and propagates failure. By default
it builds in `GITHUB_WORKSPACE`. With `stable-checkout: "true"`, it verifies the
job checkout is at the event commit, fetches it locally into
`~/.cache/repository-ci/<owner>/<repository>`, checks out that commit using the
job's existing LFS objects, and removes untracked and ignored files there
before launching Bazel. Use this option only on the dedicated runner with one
job slot. The stable path keeps Bazel's on-disk outputs and caches reusable
across jobs.

The command runs `bazel build --config=ci //...` and
`bazel test --config=ci //...` from the repository root. The root wildcard
covers first-party rulesets under their package paths. Bazel applies its
normal exclusions for manual and incompatible targets.

The test phase runs even if the build fails. The command streams diagnostics
and exits nonzero if any phase fails. Forgejo jobs shut down the Bazel
server after both phases to release runner memory.

In Forgejo, build and test logs are collapsible sections. The
job summary above the log shows whether each phase passed and how long it ran.
Expand **Build repository CI command** and then **Bazel launcher diagnostics**
to see the Bazel output for the launcher build.
The stable checkout has its own **Prepare repository checkout** group.
The CI profile omits Bazel configuration traces, progress updates, and
successful target listings. Routine output and warnings are available under
**Bazel diagnostics** inside each phase. Errors and test results are visible
when the phase is expanded. Each failed test's output has its own expandable
section in the test phase.
The same command streams plain-text diagnostics when run locally.

The runner needs the repository-managed Node.js, `bazel`, and Android
SDK/NDK prerequisites. The [runner deployment](../../infra/dev/ansible/README.md)
installs these before CI starts. Agent invocations use
`bazel run --config=ci //tools/ci`.

The pinned Bazel releases its invocation before executing the command, allowing
the build and test to reuse the root workspace's normal output base. The Linux
acceptance fixture runs the compiled CI command against a private Bazel process
fixture. It verifies root-only execution, phase and shutdown failures, missing
Bazel, and reporting. Bazel test outputs include its `ci-execution.json` receipt.

The workflow's checkout action comes from the organization-owned
[Forgejo mirror](https://forgejo.alwaldend.com/alwaldend/com_github_actions_checkout),
at the existing v4.4.0 commit. Its immutable pin selects the action implementation;
`github.sha` separately selects the source revision being built.
