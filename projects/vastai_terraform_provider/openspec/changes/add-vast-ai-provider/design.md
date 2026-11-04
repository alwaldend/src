## Context

The root module already supplies Go, Terraform CLI, Bazel, and Hugo tooling.
There is no existing provider framework. See proposal.md for authorization.

## Goals / Non-Goals

Provide a working initial on-demand GPU rental provider. Account management,
host administration, volumes, spot bids, and Registry release automation are
outside this initial implementation.

## Decisions

- Use the approved Framework v1.18.0 in the root Go module and the existing
  shared Vast.ai SDK without adding acceptance-test dependencies.
- Separate provider configuration and registration in `internal/provider` from
  resource lifecycles in `internal/resources` and queries in `internal/datasources`.
  Shared client extraction belongs in `internal/providerdata`; HTTP transport
  remains in the shared SDK's `pkg/client`. The command starts the protocol-6 server.
- Use `vastai_offers` with a JSON query and typed results. Raw API responses
  are not persisted because they can contain sensitive fields.
- Use `vastai_instance` with replacement for creation-only attributes.
  Update label and desired state in place. Creation records the contract ID
  before any subsequent API operation so failed refreshes retain ownership.
  Apply acknowledges API acceptance; boot readiness is observed on refresh.
- Never retry mutation requests automatically: rental creation is not known
  to be idempotent. Bound requests and response size, refuse redirects,
  suppress response bodies in errors, and distinguish HTTP 404 or an explicit
  null instance from omitted fields, mismatched IDs, and other failures.
- Run the actual Terraform CLI and provider binary against a loopback HTTP
  fixture. This exercises schema, planning, state consistency, and imports
  without new testing libraries or paid resources.

## Risks / Trade-offs

- Marketplace offers disappear between plan and apply: surface API failure.
- Create responses can be lost: report ambiguity; do not retry rental creation.
- Provisioning is asynchronous: computed connection fields can be null until
  refresh; apply does not guarantee SSH readiness.
- Imported instances lack their original offer and some creation settings:
  document reconciliation and avoid guessing values from unrelated API fields.
- Offline fixtures establish the documented API contract, not live-account
  compatibility. No credentials or infrastructure deployment are required.

## Decision Review

Proceed. The Framework is explicitly requested and its dependencies approved.
The shared SDK owns generated REST protocol operations. The provider owns the
bounded Terraform resource surface and reuses the SDK without duplicating it.
