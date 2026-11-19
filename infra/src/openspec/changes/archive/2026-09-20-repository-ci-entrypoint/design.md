## Context

Review of candidate `a209a2f0b618ec497d940616abd292797baeef9d` requested normal
repository CI instead of repeated infrastructure acceptance. The user confirmed
that `//tools/ci` should build and test everything and that the validation
harness should be removed fully.

## Goals / Non-Goals

Provide one CI command for root and nested normal targets. Remove the completed
acceptance harness and redundant Node certificate override. Keep provisioning,
existing authentication restrictions, and archived acceptance evidence. Do not
expand this change into repairs for unrelated failures uncovered by broad CI.

## Decisions

Decision review verdict: **revise** the existing smoke-only workflow. A wrapper
around smoke would not fulfill the requested build/test behavior. Reusing the
full-repository audit program would couple CI to audit reporting and agent
resume state; the focused command instead reads existing workspace boundaries
and executes two phases per workspace. Root `//...` alone omits nested modules.

The command uses the managed `bazel_agent` already installed by the runner role,
with explicit `--config=ci` for every child invocation. A read-only nested Bazel
probe on pinned 8.7.0 completed immediately; no separate output base is needed.
The old lock warning in help did not describe this observed behavior.

Node v22.23.1 on the runner returned HTTP 200 for Forgejo and Vault with both
`NODE_EXTRA_CA_CERTS` and `NODE_OPTIONS` unset (2026-09-20). Ansible's system CA
installation is sufficient for this Node build. The removed smoke command's
pinned Python had a different trust-store issue; no Python HTTP client remains
in the workflow.

Keep the Android installer archive in Ansible: root module analysis requires
an installed NDK even before Bazel can run `//tools/android:install`. The archive
reuses the existing package pins and avoids that circular bootstrap.

A real nested-workspace CI analysis disproved the assumption that the shared
profile was already usable: `tools/rules_docs` failed while resolving the
root-only `//tools/bazel_configs:ci_flag`. Decision review verdict: **proceed**
with moving that single setting into `root.bazelrc`. Removing the CI profile
from nested commands would discard intended CI behavior; importing the root
package into every nested module would add unnecessary dependencies.

The user then requested an action instead of a workflow shell command before
publication. A local Node.js action starts the existing entry point with a
fixed argument array, `shell: false`, and the checked-out workspace. It adds no
external dependency or bundled artifact; build/test orchestration stays in Go.
A composite action would retain a shell command and does not satisfy this
refinement. The pinned runner validates the local action metadata and the renamed `ci` job.
Action unit tests pass for argument boundaries, working directory, nonzero exit,
failed process startup and signal termination.

## Risks / Trade-offs

Broad CI can expose unrelated failures and requires the existing development
prerequisites. Normal target expansion excludes manual and incompatible targets
and does not expand optional configuration matrices. Desired Vault policy drops
the obsolete validation ref; no live policy deployment is claimed here.

## Verification

The new command's unit tests pass. After the rc correction, a real
`bazel_agent bazel test --config=ci //...` in `tools/rules_docs` built its normal
targets and passed its one test. Packaged success and failure fixtures each
observed 28 phases across 14 workspaces; a root build failure still attempted
all phases and returned exit 1. The pinned runner accepts the workflow schema.
The command and runner Ansible package build with affected semantic lint.
Node trust inspection, invocation-lock probe, and fixture logs are retained in
ignored `out/repo_ci/` for this session. These checks establish orchestration,
not a green full-repository CI run. Historical smoke run 8 remains evidence for
the earlier runner candidate only. Final aggregate delivery validation follows
archival and binds the resulting candidate.
