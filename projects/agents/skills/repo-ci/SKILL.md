---
name: repo-ci
description: >-
  Implement and diagnose repository CI workflows, reusable job commands, and
  Forgejo CI authentication. Use for workflow changes, job logic, and CI
  verification; use repo-infra for provisioning or deploying runner hosts.
---

# Work with repository CI

## Find the owning sources

Read the job implementation before editing:

- [Action](../../../../tools/ci/action/action.yml): a Node.js adapter that
  launches `bazel run --config=ci //tools/ci` directly without a shell.
- [CI command](../../../../tools/ci/README.md) and
  [package](../../../../tools/ci/BUILD.bazel): repository build/test orchestration.

Inspect the checked-out tree for workflow and deployment configuration. Shared
CI tooling can exist before workflow activation. When present, read the affected
sources:

- `.forgejo/workflows/secure.yaml`: triggers, runner selection, checkout, and
  the local `uses: ./tools/ci/action` invocation.
- `infra/forgejo_runner/README.md` and its `ansible/README.md`: execution
  boundary and prerequisites.
- `infra/vault/tf/forgejo_ci.tf` and `infra/forgejo/tf/alwaldend_repos.tf`:
  signed claim restrictions, token permissions, and protected refs.

If these files are absent, report that this revision has no configuration at
those paths and continue with the available tooling sources. Do not infer
workflow activation or deployed authentication from the tooling alone.

These sources own mutable settings; do not duplicate them in workflow logic or
this skill.

## Keep workflows maximally thin

Workflow YAML owns declarative triggers, job conditions, runner selection,
environment, checkout, and short command invocations. **All executable logic
belongs in `bazel run` targets, preferably, or committed scripts.** Workflow
`run` blocks must not contain inline programs, heredocs, API calls,
authentication flows, parsing, conditionals, or loops. Moving such logic into
composite-action YAML or encoded command arguments does not satisfy this rule.

All repository CI workflows use the local action backed by `//tools/ci`. Extend this entry point instead
of adding a second implementation. Runner acceptance probes are separate
infrastructure verification and do not belong in normal repository CI. Package
its inputs and dependencies through Bazel so local verification exercises the
same code as CI. Use `repo-bazel` and `bazel-agent` for target changes and agent
invocations. Read the current BUILD file before choosing validation commands.

## Preserve authentication and execution boundaries

Use `repo-secrets` for credentials and `repo-infra` for Vault, Forgejo, or runner
configuration. CI receives its signed job identity and restricted Vault token;
deployment identities stay with deployment. Credentials belong only in the
runtime process, never Bazel command arguments, action environments, build
definitions, logs, or artifacts.

Workflow filters, runner admission, and Vault authorization are separate
controls. Preserve the documented host-runner trust boundary. Verify matcher
semantics before treating two branch patterns as equivalent. Diagnose rejected
authentication through bounded, non-secret evidence; repair the owning cause
instead of weakening claims or clock tolerances to make CI pass.

## Verify and deliver

Start with offline package tests. Validate any affected workflow against the
pinned runner when that workflow exists. Exercise failure behavior with dummy credentials when authentication
changes. Distinguish source validation from a live CI result.

For an authorized live check, record the tested commit, source tree, run URL,
and outcome. This skill grants no deployment, protected-branch publication, or
credential-mutation authority. Preserve existing authorization and use
`repo-delivery` for publication and review.
