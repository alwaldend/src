## Context

See [proposal](proposal.md). Coordination and live cutover requirements belong
to the [DNS owner change](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/design.md).

## Goals / Non-Goals

Keep transformation provider-free and reusable by offline inspection. Provider
credentials, backend state, ownership scopes, and import receipts belong to
calling roots and the migration coordinator.

## Decisions

Use a child module for the canonical transformation. Separate logical keys
model multiple scalar values at one name/type, avoiding array-index or
address-based identities. The API is documented in the module README.

A real-provider offline plan confirmed RouterOS validates required hosturl
and initializes its API even when its resource set is empty. Global-only
owners therefore use the global entrypoint, containing the sole Cloudflare
implementation and calling normalization. The combined module delegates
there and consumes the resulting dc1 declarations through one RouterOS child.

Provider ownership defaults off for preparation. Once imported, keep it on:
turning it off in an adopted state would plan deletion and is not rollback.

Pinned provider schemas were inspected on 2026-09-12:
[Cloudflare 5.22.0](https://github.com/cloudflare/terraform-provider-cloudflare/blob/v5.22.0/docs/resources/dns_record.md)
provides scalar records and zone/record-ID imports;
[RouterOS 1.99.1](https://github.com/terraform-routeros/terraform-provider-routeros/blob/v1.99.1/docs/resources/ip_dns_record.md)
provides scalar address, cname, mx_exchange/mx_preference, ns, and text fields,
string durations, and static-ID imports. Tests use local normalization and
provider mocks; they assert no live provider behavior.

## Risks / Trade-offs

Provider credentials may still initialize when the resource gate is off;
callers keep provider wiring separate from provider-free inspection.
Exact import IDs and no-change adoption plans require separately authorized
provider inventory and migration operations.

## Validation

On 2026-09-13, all 31 normalization runs and the six combined/global provider
mapping runs passed through the Bazel normalization and provider test targets
(invocation `6e031608-e9c9-4028-9609-4eb3362c4dff`). The provider tests check
all six record types, default-disabled ownership, resource counts, updates,
deletions, and retained provider IDs. Terraform formatting also passed.
Final package formatting checks and strict OpenSpec validation remain pending.

The `routeros-eager-configuration` constraint was reproduced with Terraform
1.14.8 and RouterOS 1.99.1 in ordinary local-backend plans, both with disabled
resource ownership and with an outer module count of zero. Empty host inputs
still triggered the provider's system-resource GET and failed because the
URL had no host. Terraform test setup/teardown also exhibited the constraint.
No remote endpoint was contacted. The implementation does not supply a fake
endpoint, credentials, RouterOS version, or skip-check bypass. Operational
roots must satisfy their provider/Vault prerequisites; the module's disabled
flag controls resource ownership only. Empty-credential plan fixtures were
removed after preserving this evidence because they asserted an unsupported
operational guarantee. Offline provider lifecycle checks use provider mocks.

The user removed the optional BIND export from migration scope; canonical
normalization remains the Terraform API and its provider-free offline surface.

## Source validation

The coordinator passed the 53-target Terraform/skill package test batch and
built all 45 owning Terraform wrappers plus Vault and the infrastructure skill.
The runtime linter discovered 46 canonical files and printed 97 declaration
rows without ownership conflicts. Provider mocks cover both views; no live
provider, backend or Vault operation was run. This change completes source
preparation; operational adoption remains in the DNS coordination change.
