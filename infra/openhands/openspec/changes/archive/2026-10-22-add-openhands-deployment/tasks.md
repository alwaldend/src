## 1. Vault identity

- [x] 1.1 Add the `src_infra_openhands` AppRole as a module under
      `infra/vault/tf/approle_src_infra_openhands/` with its entity, group, SSH
      role, and server PKI role.
- [x] 1.2 Register the AppRole in the shared approle, ansible, and identity
      lookup groups.

## 2. Component and DNS

- [x] 2.1 Declare the site-local component names, `dc1` ingress targets, and
      public canvas CNAME in `infra/openhands/dnsconfig.json`.
- [x] 2.2 Wire the new DNS owner into the `infra` and `infra/dns` record
      lists.
- [x] 2.3 Add `infra/openhands/al.lua`, `BUILD.bazel`, and `README.md`.

## 3. Provisioning

- [x] 3.1 Add `infra/openhands/tf_setup` creating the canvas, secured agent
      server, and automation server VMs with hostname, address, and MAC derived
      from `dnsconfig.json`.
- [x] 3.2 Pin the Xen Orchestra provider and add the package lock and targets.
- [x] 3.3 Assign the OpenHands template, storage, and network in the owning
      `infra/xcp_ng/tf` resource-set inventory.

## 4. Ansible roles

- [x] 4.1 Add `openhands_server`, `openhands_canvas`, and
      `openhands_automation` roles installing each component natively.
- [x] 4.2 Add the deployment inventory, group variables, playbook, and
      per-component Traefik dynamic configuration.
- [x] 4.3 Pin the third-party Agent Canvas compatibility module and stage it
      for the agent server.

## 5. Ingress and the unsecured server

- [x] 5.1 Add the canvas service to `infra/ingress` without changing its
      client-authentication policy.
- [x] 5.2 Move the unsecured agent server into `users/simeonwarren/host_bot`
      and remove it from the OpenHands inventory.

## 6. Validate

- [x] 6.1 `git diff --check` passes.
- [x] 6.2 `//infra/vault/tf:tf_tests.fmt_test` and
      `//infra/openhands/tf_setup:tf_tests.fmt_test` pass.
- [x] 6.3 `//infra/openhands/ansible:ansible_bin` builds.
- [x] 6.4 The three new role targets build and `//infra/dns:config_test` passes.
- [x] 6.5 `//:buildifier_test` and OpenSpec validation pass.

## 7. Live prerequisites

- [x] 7.1 Obtain explicit authorization for the scoped Vault Terraform apply,
      review its plan, and create the OpenHands AppRole and policy.
- [x] 7.2 Generate and write the authorized component and host-bot secrets through
      their owning Vault targets without printing or reading back values.

On 2026-09-12 the user authorized the scoped Vault apply, then allowed
non-destructive dependency policy updates after review. The owning Terraform
targets created 16 OpenHands resources, added the AppRole to all four shared
groups, reconciled eight existing SecretID-cleanup policies, and removed the
unused Yandex grant from the OpenHands role. No resources were destroyed.
The final scoped plan reported no changes. All three authorized KV writes
returned version 1 with create-only CAS; values were generated in memory and
were neither printed nor read back. See the live-operation evidence in
[design.md](design.md).

XO identity synchronization, resource-set reconciliation, VM provisioning, and
Ansible deployment remain outside the authorized live operations in this change.

Validation evidence on 2026-09-12: the eight-target offline batch passed Vault,
OpenHands, and XO Terraform formatting; DNS config; Buildifier; skill discovery;
`repo-infra` eval configuration loading; and the OpenHands OpenSpec check. The
three role targets and both deployment binaries built. OpenHands packaged
Ansible syntax checking passed. Final candidate quality and semantic lint are
bound to the repository delivery receipt in ignored task scratch.
