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
runner account and its persistent workspace. The smoke workflow checks resource
capacity through the [Bazel CI command](ci/README.md), exchanges Forgejo's signed
job identity for a five-minute Vault token,
checks its limited policies, and revokes it. That identity currently grants
only token self-lookup and self-revocation.

Vault accepts the default branch from the shared repository catalog and the
exact validation branch in [CI configuration](ci.json), both covered by the
[Forgejo branch protections](../forgejo/tf/alwaldend_repos.tf). It also verifies
repository, workflow, branch type, audience and event. The workflow schedules
default/release branch jobs and reads the candidate's CI configuration before
running authentication checks; other release branches skip those checks.

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
