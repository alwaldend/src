---
title: Terraform setup
description: XCP-ng runner provisioning and owner DNS
---

The `runners` map provisions XCP-ng VMs using the shared Fedora cloud template.
The initial allocation is owned by [variables.tf](variables.tf). The single
root disk contains tools and workspaces and is expanded by cloud-init.

The existing Vault HTTP state location is retained. The operator confirmed the
legacy Proxmox VM no longer exists. [Migration configuration](migrations.tf)
forgets its old module without requesting destruction; the setup now manages
only XCP-ng VMs and current DNS records.

Use `runner.plan`, `runner.show`, and `runner.apply` to manage the XCP-ng runner
VMs. Use `dns.plan`, `dns.show`, and `dns.apply` to manage the owner DNS module.
Both apply targets require an inspected saved plan. For example, with an
absolute task-private plan path:

```sh
bazel_agent bazel run //infra/forgejo_runner/tf_setup:runner.plan -- -out=/absolute/worktree/out/task/runner.plan
bazel_agent bazel run //infra/forgejo_runner/tf_setup:runner.show -- /absolute/worktree/out/task/runner.plan
bazel_agent bazel run //infra/forgejo_runner/tf_setup:runner.apply -- /absolute/worktree/out/task/runner.plan
```

Before provisioning, deploy the component's Vault access, bootstrap its XCP-ng
OIDC identity through `//infra/xcp_ng/cmd/xo_login:bootstrap`, and apply its
resource-set assignment in `//infra/xcp_ng/tf`. Review only the resources needed
for this component. Authentication and provider data come through the owning AL
configuration; no provider credentials belong in Terraform source.

The full `tf_setup.plan`, `tf_setup.show` and `tf_setup.apply` commands reconcile
both VM and DNS resources. The scoped commands are useful when one stage needs
independent inspection. Scoped no-change plans establish only their selected
resources' state.

DNS remains enabled after the existing
[adoption](../openspec/changes/archive/2026-09-13-adopt-dns-records/design.md).
Keep existing records enabled and use the shared [DNS procedure](../../dns/README.md)
for any future ownership migration.
