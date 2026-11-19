## Why

Project states need one shared implementation of canonical DNS declarations.
The [DNS migration](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/proposal.md)
also needs provider-free normalization for offline compatibility checks.

## What Changes

- Add a provider-free normalization module and a caller-provider DNS module.
- Preserve multi-type entries, repeated destinations, record multiplicity,
  effective TTLs, and unproxied Cloudflare records.
- Package the module and exercise declarations with offline Terraform tests.

## Capabilities

### New Capabilities

- `dns-records`: canonical DNS normalization and per-record provider resources.

### Modified Capabilities

None.

## Impact

Adds `projects/tf_modules/dns_records` for owning infrastructure roots and
provider-free inspection. Provider credentials and state remain caller-owned;
no live imports or provider operations are authorized by this change.
