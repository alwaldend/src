## Why

Deferred; implementation has not started. Review finding 1 identified several
component AppRoles with the same SSH signing policy and `ansible` principal.
Separate Vault identities therefore do not establish host administration
isolation; this is a checked-in configuration finding, not observed live abuse.

## What Changes

- Scope deployment certificates and host authorization to each component's
  intended hosts, retaining explicit shared administration where necessary.
- Review deployment certificate lifetimes and extensions against actual jobs.
- Define an ordered migration that verifies authorized access and rejection of
  cross-component access before retiring broad grants.
- Accept only when an isolated SSH exercise demonstrates both outcomes and
  records non-secret evidence of the effective policy.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-vault`: constrain deployment signing and host authorization together.

## Impact

Planning only is authorized. Future implementation affects Vault Terraform,
component injection configuration, and the shared Ansible host role with
owner-specific principal inputs. Any host_bot changes must remain declared in
its checked-in Ansible. No deployment, credential changes, or new dependency
is authorized by this record.

Baseline: `e0d24cb3a`, inspected 2026-10-22. Evidence:
[shared policy membership](../../../tf/group_ansible.tf),
[SSH signing roles](../../../tf/ssh_clients.tf), and
[shared host role](../../../../../projects/ansible_collection/roles/host/).
