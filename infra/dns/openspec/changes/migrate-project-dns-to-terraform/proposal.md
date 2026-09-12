## Why

The central DNSControl apply reconciles records for every project, coupling
otherwise independent infrastructure work. Project-local Terraform state can
own individual records in the existing zones without changing DNS providers
or purchasing separate child-zone hosting.

## What Changes

- Introduce a reusable Terraform module that creates DNS records from each
  project's existing `dnsconfig.json`, preserving global and dc1 views.
  The request's `dnscontrol.json` refers to this existing input format;
  renaming files is not required.
- Invoke the module from each project's `tf_setup` root when present, and
  otherwise its `tf` root, with exactly one state owning each record.
- Migrate gradually: exclude transferred records from the central desired
  configuration and ignore their exact names and descendant subdomains in
  the corresponding operational DNSControl views before Terraform writes.
- Import existing provider records without recreating them. Preserve values,
  effective TTLs, record multiplicity, and provider-specific behavior.
- Migrate shared apex/mail records to an explicit infrastructure owner before
  retiring the central DNSControl deployment entry points.
- Investigate retaining BIND documentation through offline generation from
  the same canonical declarations and shared transformation. Retain it only
  if it needs no live-provider access or duplicated record definitions or
  translation logic; otherwise record the limitation explicitly.
- **BREAKING**: after migration, DNS changes are applied through the owning
  project's Terraform root, not the central DNSControl push command.

## Capabilities

### New Capabilities

None; this changes the existing infrastructure DNS ownership contract.

### Modified Capabilities

- `infra-dns`: project-state record ownership, staged coexistence with
  DNSControl, preserved view semantics, and optional offline BIND output.

## Impact

The migration affects `infra/dns`, the reusable module collection in
`projects/tf_modules`, DNS-owning projects under `infra`, `projects`, and
`users`, and nested project modules contributing DNS declarations. Terraform
provider packaging, project state, AL/Vault injection, Bazel runfiles, and DNS
documentation require corresponding integration work. Cloudflare remains
the global provider and MikroTik remains the dc1 provider.

This change currently contains planning artifacts only. It authorizes no
live imports, applies, credential changes, or DNSControl pushes.
