## Context

The migration replaces repository-specific wrappers with a standalone module.
The initial inventory contained 43 source lockfiles and 14 provider versions.
The user selected one version per provider, with an explicit override policy
possible later. Yandex is standardized on the already used 0.203.0 release;
other provider versions and Terraform 1.14.8 remain unchanged.

## Goals / Non-Goals

Goals: pinned Bazel acquisition; explicit provider runfiles; offline provider
installation; complete target migration; reusable execution independent of AL.
This change does not run infrastructure plans or apply live configuration.

## Decisions

- **Proceed** with Terraform's documented packed filesystem mirror layout and
  a CLI configuration containing no direct installation method. A host plugin
  cache or runtime registry resolution would weaken the declared input contract.
- Reuse rules_binary_toolchain for the Terraform executable. Correct its
  root dependency metadata so dependencies supplied by a nested module are not
  incorrectly reported as root imports.
- Keep one version resolution phase before repository registration. Conflicts
  fail now; a future explicit override can be resolved at that boundary.
- Keep AL plugin/configuration wrappers in projects/al/rules/terraform; the
  standalone module uses Go stdlib and rules_go runfiles only.
- Generated locks belong to runtime workspaces. Bazel source control owns
  provider URL/integrity selections in the shared extension declarations.

## Risks / Trade-offs

Directory runfiles preserve the existing working-directory and relative-plan
semantics. Manifest-only tests materialize declared files in test scratch;
operators without directory runfiles must explicitly choose a working directory.
The mirror controls provider installation, while provider/backend operations
can still require network access when the operator requests them.

## Validation and continuation

Implementation passed 48 Terraform consumer tests and focused extension,
runner, real-provider, and DNS fixture checks. Installed package verification
covers both automatic initialization and direct commands, with stale versions,
modified/missing/extra files, symlink entries, and TF_DATA_DIR variations tested.
A real provider initialized and validated with zero requests to the fixture's
proxy despite poisoned user CLI configuration and host caches. The catalog test
also passed through `bazel run`, exercising AL and directory runfiles.

The full matrix ran `bazel_agent bazel build //...` and
`bazel_agent bazel test //...` in the root and all 13 standalone modules.
All 14 builds and all 13 standalone test commands passed. Root tests reported
295 passed, one secret-scanner failure, and one skipped test (not named by the
terse summary). The scanner found four historical synthetic URI fixtures in
commits preceding the task base d40966740f0e7f708bf4b9ab69e473fb9b56663d:
infra/xcp_ng/cmd/xo_login/main_test.go, tools/agents/evidence/evidence_test.go,
tools/agents/api/v1alpha1/phase3_test.go, and
tools/repo_delivery/main/go/runner_test.go. Its additional finding in the new
GitLab provider SRI was verified against publisher metadata and resolved with a
line-specific annotation; a focused repeat confirmed only the four historical
findings remain. The repository scanner policy was preserved.

Configured formatting and Buildifier passed after explicitly formatting two
changed BUILD files excluded by the aggregate formatter. Representative
main/external runfiles configurations, the canonical GitLab mirror manifest
entry, and both rendered module/launcher documentation pages were verified.
The command matrix and restricted diagnostics are task-local under
out/rules_terraform and out/full-repo-check. No live infrastructure operations
were performed for this migration. Publication identity belongs to Git and the
delivery receipt; final candidate validation follows the archived source.
