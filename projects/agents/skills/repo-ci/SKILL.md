---
name: repo-ci
description: >-
  Implement and diagnose repository CI workflows, reusable job commands, and
  Forgejo CI authentication. Use for workflow changes, job logic, and CI
  verification; use repo-infra for provisioning or deploying runner hosts.
---

# Work with repository CI

## Find the owning sources

Read the affected workflow and job implementation before editing. For the
current Forgejo smoke workflow, use:

- [Workflow](../../../../.forgejo/workflows/secure.yaml): triggers, runner
  selection, checkout, environment, and command invocation.
- [CI package](../../../../infra/forgejo_runner/ci/BUILD.bazel): executable
  targets and their runtime inputs; the smoke entry point is
  `//infra/forgejo_runner/ci:smoke`.
- [Runner owner](../../../../infra/forgejo_runner/README.md) and
  [CI configuration](../../../../infra/forgejo_runner/ci.json): execution
  boundary and validation branch.
- [Vault role](../../../../infra/vault/tf/forgejo_ci.tf) and
  [Forgejo protection](../../../../infra/forgejo/tf/alwaldend_repos.tf): signed
  claim restrictions, token permissions, and protected refs.
- [Validation deployment](../../../../infra/forgejo_runner/ansible/README.md):
  publishing and checking an exact candidate on the protected validation branch.

These sources own mutable settings; do not duplicate them in workflow logic or
this skill.

## Keep workflows maximally thin

Workflow YAML owns declarative triggers, job conditions, runner selection,
environment, checkout, and short command invocations. **All executable logic
belongs in `bazel run` targets, preferably, or committed scripts.** Workflow
`run` blocks must not contain inline programs, heredocs, API calls,
authentication flows, parsing, conditionals, or loops. Moving such logic into
composite-action YAML or encoded command arguments does not satisfy this rule.

Extend the owning command instead of adding a second implementation. Package
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

Start with offline package tests and workflow validation against the pinned
runner. Exercise failure behavior with dummy credentials, including redirect
rejection and safe diagnostics when authentication changes. Distinguish source
validation from a live CI result.

For an authorized live check, use the owning validation deployment and record
the tested commit, source tree, run URL, and outcome. This skill grants no
deployment or credential-mutation authority. Preserve existing authorization
and use `repo-delivery` for publication and review.
