## Why

The [shared DNS migration](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/proposal.md)
requires `rules_dnscontrol` to own its DNS records in a dedicated Terraform state.
Its existing `dnsconfig.json` remains the canonical declaration source.

## What Changes

- Add a project-owned `tf` root consuming the existing `dnsconfig.json`.
- Export its source for a root-workspace operational wrapper under
  `infra/dns/projects`, preserving the standalone module dependency boundary.
- Inject provider credentials through the normal `tf=main` AL flow, with the
  project AppRole and required Vault credentials deployed before root use.
- Default DNS record creation to disabled until authorized adoption.

## Capabilities

### New Capabilities

- `project-dns`: project Terraform source and independent DNS state integration.

### Modified Capabilities

None.

## Impact

Project Terraform source and Bazel exports, the root-workspace DNS wrapper,
and the project's Vault/AppRole wiring are affected. The runtime
[DNS linter](../../../../../../infra/dns/README.md) discovers canonical declarations
and displays its table directly. Live adoption remains tracked by the shared
migration and requires separate operational authority.
