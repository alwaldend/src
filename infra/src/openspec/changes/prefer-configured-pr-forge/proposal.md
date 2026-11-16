## Why

Repository agent policy names GitHub as the only delivery destination even
though repositories may use Forgejo as their configured origin or upstream.

## What Changes

- Authorize normal delivery to the verified selected repository remote, keeping
  PR destination on its forge and preserving unrelated-work and approval guards.
- Route detailed selection and CLI instructions to the canonical AL delivery skill.
- Expand the fj tool landing page with invocation and login guidance.

## Capabilities

Documentation-only integration correction; `skip_specs: true`.

## Impact

Root agent policy and `tools/fj/README.md`. Procedure and evals remain owned by
`projects/al/openspec/changes/document-forge-remote-selection`. This change does
not alter remotes, credentials or deployments.
