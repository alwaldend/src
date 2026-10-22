## Context

See the proposal for scope. Go uses identifier capitalization for package visibility; Python uses naming conventions and explicit package exports. Existing Python guidance already preserves public serialized fields and framework contracts.

## Goals / Non-Goals

Provide equivalent visibility decisions in both skills without equating the languages' export mechanisms. This change does not rename application declarations or introduce visibility lint rules.

## Decisions

- Require an actual consumer or a defined contract rather than hypothetical reuse. Same-package callers and tests alone do not justify an export.
- Use unexported Go names and leading underscores for Python implementation details; retain explicit `__all__` for intentional package exports.
- Preserve existing APIs and required interface, framework, and serialization names. Check consumers before narrowing existing visibility.
- Add behavioral evaluation scenarios before changing guidance. Run the existing offline validation and skill packaging checks; these do not measure agent compliance.

## Risks / Trade-offs

Overly mechanical privatization can break consumers or serialized names. The guidance requires checking consumers and retaining contractual public names. Visibility justification belongs in the API documentation or change rationale, without a separate registry.

## Verification

Both offline evaluation configurations accepted the new cases and staged the canonical skills. The repository formatter completed without unrelated changes. Final exact-candidate checks and publication are recorded in `out/justified-public-apis/prepare.json` and its validation and delivery receipts.
