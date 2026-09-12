## Why

Project DNS records need the same ownership and independent Terraform state as
the projects they serve. The previous central DNSControl writer coupled all projects.

## What Changes

- Reuse canonical per-owner `dnsconfig.json` from each `tf_setup` root when
  present, otherwise `tf`, with a shared provider-aware Terraform module.
- Add missing dedicated AppRoles as modules under
  `infra/vault/tf/approles/<name>/*.tf` and grant only required DNS reads.
- Discover declarations at runtime, reject multiple source files managing one
  DNS name, and print a table. Do not maintain a central ownership registry.
- Remove central DNSControl entrypoints. No exporter or BIND generation is
  retained; historical snapshots are labeled as such.
- Complete owner-local cutovers under the user's separate deployment authorization,
  preserving existing records and retaining recovery evidence.

## Capabilities

### Modified Capabilities

- `infra-dns`: project-state management and runtime declaration ownership lint.

## Impact

Affects DNS-owning infrastructure, projects (including nested modules), user
hosts, reusable Terraform modules, Vault AppRoles, and infrastructure guidance.
Cloudflare and MikroTik remain the providers. [Operational completion](operations.md#verified-completion)
records the separately authorized deployment; source implementation alone does
not authorize future live writes.
