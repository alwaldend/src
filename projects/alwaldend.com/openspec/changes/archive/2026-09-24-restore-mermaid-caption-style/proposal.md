## Why

Removing output rewriting also removed the previous caption styling. Restore
that appearance with Mermaid theme inputs while retaining its existing Dagre
layout. The user explicitly ruled out a layout redesign.

## What Changes

- Restore rounded, bordered edge labels and container titles, and prior spacing.
- Supply the resolved site palette to Mermaid's native input theme CSS.
- Verify both theme modes, standalone images, and unchanged blog diagrams.

## Capabilities

### New Capabilities

### Modified Capabilities

- `project-alwaldend-com`: caption styling for adaptive documentation diagrams.

## Impact

Only documentation rendering changes. No layout extension, custom edge router,
SVG post-processing, dependency change, or live deployment is introduced.
