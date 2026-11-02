## Context

The implementation and command-level protocol fixtures currently belong to
infra/local_inference. See proposal.md for the extraction rationale.

## Goals / Non-Goals

Preserve the existing command flags, loopback endpoints, bounded bodies,
namespace mapping, and JSON/SSE behavior. Host services remain deployment-owned.

## Decisions

Use the plural `responses` to match the Responses API. Keep the existing
`cmd/responses_bridge`, `internal/responses`, and `test/protocol` roles under
the new project. Export a root command alias for repository consumers; keep
implementation visibility within the project. A separate root-based PR avoids
including unrelated host infrastructure or a stacked base. Existing protocol
fixtures move with the implementation and retain the optional Codex round trip.

## Risks / Trade-offs

- Stale imports or runfile labels: verify the actual command-level acceptance.
- Accidental infra dependencies: build the new project from master and inspect
  its imports and target dependencies.
- Premature consumer migration: keep the extraction PR a prerequisite for the
  independently based infrastructure PR; publish the consumer only after it
  can pass its checks.

## Migration Plan

Publish the extracted project first. Update the existing infrastructure PR to
consume it and remove duplicate implementation once the prerequisite is merged.
No deployment is performed.
