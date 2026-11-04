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
The provider implementation is delivered separately from Terraform roots,
AppRoles, token injection, and backend setup. Their integration is described
below; live deployment and Registry publication remain outside this change.

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
packaging, and project documentation. That implementation is merged in PR #255.
PR #249 now retains only Terraform configurations and credential integration.
The split creates no merge or deployment authority.

## Test-root integration

The test/team root consumes its generated team token through the project's own
AppRole and injector, with a separate tf=test Vault backend. It manages a test_
API key, TEST_ environment marker, and unassigned test_ role requesting current-account
reads. The bootstrap catalog and token delivery belong to infra/vastai;
the AppRole declaration belongs to infra/vault. The test root contains no
compute, storage rental, credit transfer, or membership invitation, and adds no
account-ID guard. Reuse the merged provider's package layout and reviewed
lifecycle behavior without restoring the older combined provider snapshot.

## Live verification

The user subsequently authorized deployment on 2026-11-04. Vast.ai rejects SSH
keys in team context, so the live root uses an API key instead. Personal SSH
behavior remains covered by the offline lifecycle fixture. Masked environment
reads preserve stored values; imports require a configured value on the next
apply. Role state separates requested permissions from effective server grants,
retaining a prior effective document to distinguish normalization from drift.

The AppRole and bootstrap were applied, followed by the three-resource team
root. Both Vast.ai roots produced no-change plans; the personal bootstrap
context stayed unchanged. No compute was rented and no credit transferred.
