## Why

This project needs an independently owned Terraform stage for its DNS records
as part of the [repository DNS migration](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/proposal.md).

## What Changes

- Consume the existing project DNS declaration through the shared DNS module.
- Package the project stage with its dedicated AppRole and state backend.
- Keep record creation disabled until an authorized ownership cutover.

## Capabilities

### New Capabilities

- `project-dns`: Project-local Terraform DNS configuration and packaged inputs.

### Modified Capabilities

None.

## Impact

This project's Terraform package, AL configuration, and Bazel packaging. Shared
module behavior and migration sequencing remain with the linked owners.
