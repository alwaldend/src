## Context

This project owns Terraform schemas, planning, resource state, and protocol-6
plugin startup. The merged internal/vastai_sdk owns generated API operations,
bounded transport, and response validation.

internal/provider owns registration and provider configuration;
internal/resources owns resource lifecycles; internal/datasources owns offer
queries. Both consume internal/providerdata for checked dependency extraction,
avoiding dependency cycles through the registration package.

## Goals / Non-Goals

Manage Vast.ai account resources and preserve ownership through failures.
Terraform deployment roots, AppRoles, Vault configuration, token injection,
backend setup, and checked-in Terraform examples remain outside this PR.

## Decisions

- Consume the SDK's single public client with explicit Endpoint and APIKey
  Options. Keep Terraform behavior independent of generated protocol types.
- Share provider-client configuration across resources and the offer data source.
  Nil provider data remains unconfigured; unexpected data produces one diagnostic.
- Pass a providerdata.Data struct containing Client through both Framework data
  channels. Consumers retain that typed payload rather than extracting only its
  client, leaving one shared dependency contract across the package boundaries.
- Keep native Framework schema structs explicit: state-model tags describe value
  mapping, while schema metadata defines required/computed/sensitive fields,
  defaults, and replacement. Expand declarations for readability. Retain the
  documented typed-model ValidateConfig approach without adding a reflection
  layer or an optional validation dependency.
- Keep an explicit guard in the Framework-required API-key Update method because
  every configurable key attribute requires replacement. Verify replacement
  through actual key revocation and recreation. Name the stored identity and
  scoped team client in deletion, and handle already-revoked keys explicitly.
- Normalize compound API launch modes before saving instance state. Retain the
  prior supported mode when the response is equivalent; choose a deterministic
  supported mode on import. Jupyter launch modes take precedence over accompanying
  SSH access modes. Missing or unsupported modes are read errors, preserving state.
- Reject empty image identifiers and non-positive disk allocations before
  mapping any refreshed instance fields into state. Partial successful responses
  must not introduce replacement plans by overwriting creation-only attributes.
- Support SSH public keys, environment variables, and team roles with CRUD and
  import. Treat malformed replies as errors rather than remote deletion.
- Preserve equivalent permission JSON to prevent false drift. Team role names
  remain replacement-only, with reads and deletion using escaped names and
  permission updates using the captured numeric identifier.
- Create teams without transferring credit. Generate an independent management
  key before team creation; never change the personal key's context.
- Record confirmed team creation separately from management-key ownership.
  Failed creation cannot adopt or destroy an existing team with the same name.
  Bind cleanup to the captured account ID.
- Retain management-key ownership after external team deletion, plan explicit
  replacement, and preserve state when key revocation fails.
- Keep generated key values sensitive and support their consumption by other
  Terraform resources. One-time credentials cannot be recovered through import.
- Verify behavior through the actual Terraform CLI against loopback fixtures.
  Fixture credentials are synthetic; no live account or compute is used.

## Delivery Boundary

The provider-only PR targets trunk and consumes merged SDK PR #251. It contains
Go implementation, its approved dependency graph, offline fixtures, provider
packaging, and project documentation. PR #249 retains Terraform configurations
and credential integration and can drop the provider after this prerequisite
merges. The split creates no merge or deployment authority.
