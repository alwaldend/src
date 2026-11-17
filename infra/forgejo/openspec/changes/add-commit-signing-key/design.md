## Context

See proposal.md. The existing setup stage precedes service deployment and owns
the Forgejo AppRole's Vault-backed backend. Its KV subtree already permits the
new signing-key record. Service Terraform requires a running Forgejo instance,
so generation belongs in setup rather than the service API stage.

## Goals / Non-Goals

Preserve a dedicated signing identity across rebuilds. Source work only; no
live generation, Vault writes, deployment, rotation, or account-key registration.

## Decisions

Proceed with Ed25519 SSH signing. The pinned Forgejo 17 documentation recommends
SSH when supported, and the host already uses modern Git and OpenSSH. GPG would
add keyring management without improving the requested lifecycle.

A creation-only Terraform provisioner runs a packaged Go tool with the existing
Vault environment. The tool generates in memory and writes directly to KV v2
with CAS=0, never into Terraform attributes. Cached records must contain a valid
matching private/public pair and fingerprint. Access errors, soft deletion,
destroyed versions, and existing metadata without usable data fail closed.
Normal plans and applies do not rerun generation; explicit replacement reuses
the cache. The setup backend lock owns serialization.

The reusable command is owned by tools/vault/cmd/ssh_key and accepts neutral
VAULT_SSH_KEY_MOUNT and VAULT_SSH_KEY_PATH inputs. Forgejo retains its secret
path, provisioning resource, deployment wiring, and combined CLI/Terraform
acceptance fixture under infra/forgejo/test/signing_key. The tool has no
infrastructure source dependency.

Ansible injects the cached pair through AL environment resources. The shared
Forgejo role installs an optional signing pair before activation. Its private
key path remains under the existing backed-up Forgejo root. Existing private
content must match the declared key; rotation requires separate reconciliation.
The public key can be repaired. Configuration points SIGNING_KEY at the public
file, as required by Forgejo, and retains merge/CRUD signing and trust policies.

## Risks / Trade-offs

- Missing Vault data must not create a new signing identity after deletion:
  check metadata as well as data before creation; preserve CAS=0.
- A lost Vault write response is recoverable by rerunning and validating the
  stored pair; no remote service object is created before persistence.
- Terraform does not detect Vault drift after bootstrap. Deployment still
  requires the record, and explicit replacement revalidates it.
- Existing unconfigured GPG state is preserved; enabling this dedicated SSH
  identity is a future operator-reviewed deployment, not an automatic rotation.

## Migration Plan

After separate live authorization, review setup Terraform's normal interactive
apply, then run the service-only Ansible deployment. Confirm the instance SSH
signing endpoint and a Forgejo-created signed commit. Preserve Vault history and
the private recovery archive. Do not delete or rotate this long-lived identity
as a rollback mechanism.

Reference: https://forgejo.org/docs/v17.0/admin/advanced/signing/
