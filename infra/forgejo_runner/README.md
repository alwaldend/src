---
title: Forgejo runner
description: XCP-ng Forgejo Actions runners with Vault CI authentication
tags:
  - forgejo
  - ansible
  - xcp_ng
---

This project provisions Forgejo Actions runners on XCP-ng. The initial runner
is `secure`, using the allocation in [Terraform](tf_setup/variables.tf),
address in [DNS](dnsconfig.json), and host/labels in [inventory](ansible/inventory.yaml).
The deployment applies the shared `host`, `dev_vm` and `forgejo_runner` roles,
including Bazel tools for the runner account.

[Terraform setup](tf_setup/README.md) owns VM and DNS provisioning;
[Ansible](ansible/README.md) owns registration and host configuration.
The existing component AppRole owns deployment state and the registration
credential. CI authenticates separately through the
[Vault JWT role](../vault/tf/forgejo_ci.tf); deployment tokens are never installed
on the runner.

## CI trust boundary

The repository-scoped runner executes `secure:host` jobs as its dedicated
unprivileged account. Repository writers must be trusted: host jobs share the
runner account and its persistent workspace. The [workflow](../../.forgejo/workflows/secure.yaml)
runs the shared [repository CI command](../../tools/ci/README.md) to build and
test normal targets. Runner acceptance checks are not part of repository CI.

The declared Vault role accepts only the default branch from the shared
repository catalog, covered by [Forgejo branch protection](../forgejo/tf/alwaldend_repos.tf).
It verifies repository, workflow, branch type, audience and event. Its
five-minute tokens permit only self-lookup and self-revocation. The workflow
schedules default/release branch jobs; release jobs are not granted this Vault
identity. Build and test commands do not request Vault credentials.

Forgejo 15 has no protected-only runner setting. A different workflow on an
unprotected branch can target the same repository runner; workflow conditions
cannot prevent that. Its `ref_protected` OIDC claim is hardcoded false, so the
Vault role binds explicit protected refs instead. These are documented
[upstream limitations](https://forgejo.org/docs/latest/admin/actions/security/),
not a guarantee of runner admission enforcement.

## Add a runner

Add a named allocation to `runners`, a corresponding owner DNS record, and an
Ansible inventory entry with its runner name and label. Review/apply the scoped
runner and DNS plans, then deploy with `--limit <hostname>,localhost` so the
controller trust preparation also runs. Runners with a
different trust boundary need separate repository registration and CI policy.

The operator confirmed the legacy Proxmox VM was already gone. Its stale state
entry is retired without a VM destruction request; [migration configuration](tf_setup/migrations.tf)
records that transition. The active inventory and runner allocation contain
only `secure`.
