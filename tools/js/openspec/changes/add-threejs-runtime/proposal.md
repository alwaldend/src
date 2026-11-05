## Why

The 3D Mermaid renderer needs the already approved Three.js runtime. Extract
its dependency pin from PR #211 so it can be reviewed and merged independently
of the renderer and snapshot implementations.

## What Changes

- Add `three@0.186.1` to the existing tools pnpm workspace using its lock workflow.
- Exercise the runtime from the existing strict TypeScript Node E2E fixture.
- Preserve the merged TypeScript setup, existing dependency versions, and disabled
  lifecycle scripts.

## Capabilities

No product requirements change. This dependency and build-support record uses
`skip_specs: true`.

## Impact

The owning files are `tools/package.json`, `tools/pnpm-lock.yaml`, and the
`tools/js` documentation and acceptance fixture. All diagram implementation
remains in `projects/mermaid_3d` and is delivered separately. The runtime pin
was explicitly approved by the user; its type declarations are already present.
