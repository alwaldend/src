---
title: Forgejo runner CI
description: Bazel entry point for the runner and Vault smoke check
---

The [workflow](../../../.forgejo/workflows/secure.yaml) checks out its exact
event revision and runs `bazel run --config=ci //infra/forgejo_runner/ci:smoke`.
The committed [implementation](smoke.py) owns resource checks, reference
selection, OIDC issuance, Vault login, policy/TTL validation and revocation.
It reads the checked-out [CI configuration](../ci.json) through Bazel runfiles.
The same script can run with `python3 infra/forgejo_runner/ci/smoke.py` when
Bazel is unavailable; the workflow uses the pinned Bazel/Python toolchain.

Authentication uses the job's runtime environment. Never pass tokens through
Bazel flags, `--action_env`, `--repo_env`, BUILD attributes, logs or artifacts.
The command contacts Forgejo and Vault at runtime; build actions and unit tests
do not contact those services. Unit tests use mocked HTTP exchanges to cover
authentication success/failure, revocation and redirect rejection.

The [checkout action](https://github.com/actions/checkout/releases/tag/v4.4.0)
is pinned to an immutable commit in the workflow, uses the event SHA and does
not persist its credential in the checkout. The runner role supplies Node.js
for JavaScript actions. The workflow selects the host's existing Fedora CA
bundle for Node.js and the pinned Python runtime to trust the internal services.

For agent validation, use `bazel_agent bazel test
//infra/forgejo_runner/ci:smoke_test`. Live validation uses the owning
[Ansible CI entry point](../ansible/README.md); the
[Vault policy and runner trust boundary](../README.md#ci-trust-boundary) still
apply. Repository CI authoring follows the
[repo-ci skill](../../../projects/agents/skills/repo-ci/SKILL.md).
