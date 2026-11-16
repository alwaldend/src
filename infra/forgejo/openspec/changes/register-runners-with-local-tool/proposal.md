## Why

Runner credentials need a repeatable IaC bootstrap without a separate Ansible
registration controller or another Terraform provider.

## What Changes

- Package a local Go command that registers a repository-scoped runner through
  Forgejo's API and stores its UUID and token directly in Vault KV v2.
- Invoke the tool from the service Terraform through built-in `terraform_data`
  and `local-exec`, retaining the existing AL credentials and backend.
- Run only when the stable Terraform resource is created, without replacement
  triggers. Cache credentials in Vault and reject same-name registrations
  without cached credentials after interrupted creation.
- Document the command's lifecycle, recovery, and configuration consumers.

## Capabilities

### New Capabilities

- `runner-credentials`: Local Terraform-driven runner credential bootstrap.

### Modified Capabilities

None.

## Impact

`infra/forgejo` command and Terraform packaging, isolated command acceptance
coverage, and documentation. Existing pinned Go dependencies suffice. No live
registration, Vault write, or deployment is authorized by this source change.
