## Context

See [proposal.md](proposal.md) and the linked DNS migration. Existing component
AppRoles already own their state subtrees. The generic AppRole module separates
operator SecretID creation from the role token's secret access.

## Goals / Non-Goals

Group each missing owner's identity and DNS policy into a module under
`tf/approles/<name>/`. Preserve existing resource addresses and Ansible
authentication. No live Vault reads, writes, plans, or applies are part of
implementation validation.

## Decisions

- Create `src_infra_mikrotik`, `src_infra_nas`,
  `src_users_simeonwarren_host_bot`, and the missing project identities. The
  existing router1 identity owns router-specific SSH and PKI privileges, while
  the Mikrotik DNS component also covers switch and bare-metal records; a
  dedicated identity gives the new DNS stage only its required permissions.
- Each new owner module composes `projects/tf_modules/vault_approle` with a
  shared local DNS access module. The latter owns the existing provider-secret
  references and read-only policy generation. This avoids duplicating secret
  path facts across owner modules.
- Existing roles receive the DNS access policy through their established
  policy inputs. No existing AppRole, identity, or policy address moves.
- New DNS identities join a dedicated group without extra policies. Reusing
  the existing general AppRole group would grant child-token creation and PKI
  permissions that DNS-only stages do not need. Disable Yandex policy
  attachment through the generic module's existing input.

## Risks / Trade-offs

- Shared provider credentials retain their existing provider-side scope;
  Vault grants are limited to reads of the exact required KV entries. Record
  ownership enforcement belongs to the linked DNS migration.
- The generic AppRole module still declares its unused Yandex policy resource
  when attachment is disabled. Changing that shared lifecycle is outside this
  additive identity migration; no new role token receives that policy.
- Declarations do not prove bootstrap access or deployed state. Operator
  verification follows an explicitly authorized Vault apply.

## Migration Plan

Validate and publish source first. An authorized operator then applies the
additive Vault configuration before preparing owner DNS plans. Provider
credentials remain at their existing KV references; provision the additional
Cloudflare zone identifier field using the owning secret workflow described
by the DNS migration. Removing the newly added DNS policy attachments and
unused roles reverses this source change without moving existing resources.

## Evidence and Next Action

Role and view inventory derives from the DNS migration's source inventory at
revision `cb4b5bd3a37d29dc778ed140ab72d05efef7816a`, observed 2026-09-12.
Live identity and provider state were not inspected. Source checks and an
independent review confirm all 33 new identities and 11 existing policy
attachments match owner views, all relative module paths resolve, local source
globs include the modules, the DNS group adds no policy or downstream
membership, and existing Terraform declarations retain their addresses.
`git diff --check` passes. Next action: run the owning Terraform formatting
and packaging checks, then strict OpenSpec validation and source delivery.

## Source validation

The coordinator passed the 53-target Terraform/skill package test batch and
built all 45 owning Terraform wrappers plus Vault and the infrastructure skill.
The runtime linter discovered 46 canonical files and printed 97 declaration
rows without ownership conflicts. Provider mocks cover both views; no live
provider, backend or Vault operation was run. This change completes source
preparation; operational adoption remains in the DNS coordination change.
