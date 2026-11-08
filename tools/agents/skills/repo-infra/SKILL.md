---
name: repo-infra
description: >-
  Implement, review, and validate this repository's infrastructure: Terraform
  stages, Bazel-packaged Ansible, DNS declarations, Vault-backed injection, and
  read-only deployment diagnostics. Use for infra/** and host_bot work,
  BUILD.bazel configuration, tf_setup/tf/ansible packages, dnsconfig.json records,
  or repository Vault targets. Do not use it to authorize live deployment.
---

# Work with repository infrastructure

Infrastructure spans several mechanisms that share one packaging and injection
flow. Read [references/flow.md](references/flow.md) first, then the reference
for the mechanism you are changing. Never treat this skill as authority for a
live operation: mutating infrastructure requires the user's explicit request
for that exact operation and scope.

## Start from ownership and authority

1. Read the nearest `AGENTS.md`, the owning `README.md`, `BUILD.bazel`, and the closest sibling component.
2. Identify which stage owns the change: `tf_setup` (provisioning/bootstrap),
   `tf` (service or API configuration), `ansible/` (host deployment), or
   a component-owned `dns/` Terraform root (records). Prefer a separate DNS
   root with its own state and DNS-only credentials, independent of VM
   provisioning. Keep canonical declarations in the component's
   `dnsconfig.json`. Reusable behavior belongs in `projects/tf_modules` or
   `projects/ansible_collection/roles`.
3. Inspect the real generated targets with `bazel query`; label
   maps vary by package and are the source of truth. Follow `al-bazel` and
   `repo-bazel` for every invocation.
4. Load `repo-secrets` for secret-valued work and `al-openspec` when the change
   needs durable requirements or continuation state.

## Keep the shared flow intact

Terraform packages declare `.tf` files and module inputs in `data`, select
provider labels through `providers`, pass one or more `al_config` labels, and
expose commands with `al_terraform_binary_map` from
`//tools/al:al_terraform.bzl`, passing `configs` and `plugin_labels` directly. Follow the
[Terraform rule owner](../../../../tools/rules_terraform/README.md) for provider
pinning and installation; do not add source `.terraform.lock.hcl` files.
Ansible packages package runtime files as
`pkg_files`/`pkg_filegroup`, build with `al_ansible_binary`, and run through
`al_binary_run`/`al_binary_run_map`. Use `al_inject` with named `res_<type>` maps for injector declarations.
Authentication, state backend, and secret values arrive through small authentication, plugin and plugin-call targets
composed with `al_config` and selected by labels (`tf=setup`, `tf=main`,
`ansible=1`). Component auth explicitly extends `//tools/al:auth`; aggregate
reference order does not select overrides. Keep short declarations inline.
Use local `defs.bzl` only for significant reuse or a substantial readability gain
from complex graphs, long templates or construction logic; call common target
macros directly. Use `vault_tf_backend` for each backend call. Plugin executable runfiles
propagate through configuration dependencies; wrappers select labels with
`plugin_labels` dictionaries. Never replace that flow with
committed credentials, a host `vault`/`terraform` binary, or an unprefixed
ad-hoc command.

`al_plugin(plugin = "tf_backend", ...)` is the repository's Vault-backed
HTTP state backend: the plugin stores lock and state as Vault KV entries under
the AppRole's own path and creates them on first use. It is not an object-storage
bucket, and switching to `backend "s3"` or adding a `tf_backend` module bucket
is not required for a new component.

Do not commit `.terraform/`, state, plans, environment files, or credentials.

## Terraform

Prefer the pinned provider's documented environment variables for credentials
and endpoints, supplied through the existing injection flow and shared helpers.
Avoid redundant Terraform variables and provider arguments when those native
variables suffice. Keep explicit provider inputs when the provider lacks
environment support or aliases need different identities in the same process.

Read [references/terraform.md](references/terraform.md) for the stage choice,
module conventions, text-mode interactive review, format/plan targets, and
state-migration rules. Apply the root `AGENTS.md` plan-artifact policy. Treat
`tf.apply`, `destroy`, `import`, and state operations as mutating.

## Ansible

Read [references/ansible.md](references/ansible.md) for the packaged execution
model, idempotency, handler, and injection rules. Never run a playbook target
as a syntax check; it can mutate inventory hosts.

## DNS

Read [references/dns.md](references/dns.md) before adding or changing records.
Record definitions remain in each owner's `dnsconfig.json`; Terraform modules
consume them. The DNS linter discovers those files at runtime, prints their
declarations as a table, and enforces one source file per domain name.

## Vault operations

Read [references/vault.md](references/vault.md) for AppRole bootstrap,
component credential setup, KV capabilities, safe token metadata, refresh,
and failure handling through repository targets.

## Diagnostics

Prefer the owning component's packaged SSH target for remote inspection;
discover its label and authentication before attempting direct host SSH.
If none exists, follow the missing-target procedure in the diagnostics reference.

Read [references/diagnostics.md](references/diagnostics.md) for the read-only
audit procedure and the evidence a finding must carry.

## Validate without mutating

| Change         | Non-mutating check                                                                     |
| -------------- | -------------------------------------------------------------------------------------- |
| Terraform      | `bazel test //path/tf:tf_tests.fmt_test`; an authorized live review can use `:tf.plan` |
| Ansible        | `bazel build //path/ansible:ansible_bin`                                               |
| DNS            | `bazel run //infra/dns:lint`; `bazel test //infra/dns:config_test`                     |
| BUILD/Starlark | `bazel test //:buildifier_test`                                                        |
| Any            | `git diff --check`                                                                     |

Start verification with `git diff --check`, then the narrowest package targets.
A Terraform plan contacts the configured backend and providers; it is a live
review, separate from offline implementation validation. When that access is
unauthorized or unavailable, still run formatting and Bazel build/tests and
report the precise limitation. Summarize plans for additions, changes, destroys, replacements, and
sensitive or security-impacting changes without copying secret values.
