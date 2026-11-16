# Trunk CI acceptance

The prior runner PR was merged and its remote feature branch deleted. The clean
task branch was rebased to Forgejo master
`47fc30d1b6de21d487e034ca46827629a42297ba` before implementation.

The executable acceptance fixture initially failed because the actual compiled
command attempted the removed `al bazel` interface. After direct Bazel invocation
was restored, command/action tests and all 30 selected checks passed with
`--config=ci`. The fixture covers successful execution, build/test failure,
missing Bazel and invalid boundaries, with six observed commands across three
workspaces. It retains `ci-execution.json` in Bazel test outputs. This is an
isolated process-boundary check, not a full build of the monorepo.

The deployed runner's pinned 13.2.1 validator accepts the workflow. The existing
checkout mirror's live v4.4.0 tag resolves to the selected immutable commit.
Public repository metadata confirms default_branch is master and has_actions is
true; no repository settings were changed.
An HTTPS probe as forgejo-runner using host Node 22.23.1, with Node CA overrides
unset, returned HTTP 200 from Forgejo. The redundant NODE_EXTRA_CA_CERTS workflow
override was removed; trust remains owned by the existing host CA installation.
Checkout includes LFS assets and disables credential persistence. Workflow
inspection confirms only a master push trigger and exactly two steps: checkout
and the local action. A local workspace action must be present before it is
loaded; checkout therefore remains in the workflow.

The existing action launches `bazel run --config=ci //tools/ci`. Workspace
membership remains with the shared discovery implementation and .bazelignore.
Normal manual/incompatible target exclusions apply. No deployment target,
infrastructure probe, workflow dispatch or feature-branch CI run was added.

Strict OpenSpec validation and semantic lint pass. Activation and the first full
trunk CI run require this change to reach master; no merge or protected-ref push
was performed as part of preparing the PR.
